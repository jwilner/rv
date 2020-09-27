package platform

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/jwilner/rv/internal/models"
	"github.com/jwilner/rv/pkg/pb/rvapi"
)

type tokenManager struct {
	signingKey []byte
}

func (t *tokenManager) parse(token string) (*Claims, error) {
	var c Claims
	_, err := jwt.ParseWithClaims(token, &c, func(*jwt.Token) (interface{}, error) {
		return t.signingKey, nil
	})
	if err != nil {
		if vErr := new(jwt.ValidationError); errors.As(err, &vErr) {
			if vErr.Errors&jwt.ValidationErrorExpired > 0 {
				return nil, status.New(codes.Unauthenticated, "token invalid").Err()
			}
		}
		return nil, err
	}
	return &c, nil
}

func (t *tokenManager) write(c *Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(t.signingKey)
}

func loadClaims(ctx context.Context) *Claims {
	c, _ := ctx.Value(claimsKey).(*Claims)
	return c
}

func userID(ctx context.Context) int64 {
	c := loadClaims(ctx)
	if c == nil {
		return 0
	}
	i, _ := strconv.ParseInt(c.Subject, 16, 64)
	return i
}

func tokenInterceptor(tokM *tokenManager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var token string
		if tok, ok := ctx.Value(parsedCookieKey).(string); ok {
			token = tok
		}
		if token == "" {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return handler(ctx, req)
			}
			if tok, ok := md["rv-token"]; ok && len(tok) != 0 {
				token = tok[0]
			}
		}
		if token == "" {
			return handler(ctx, req)
		}
		claims, err := tokM.parse(token)
		if err != nil {
			return handler(ctx, req)
		}
		return handler(context.WithValue(ctx, claimsKey, claims), req)
	}
}

func newClaims(userID int64) *Claims {
	return &Claims{jwt.StandardClaims{Subject: strconv.FormatInt(userID, 16)}}
}

type Claims struct {
	jwt.StandardClaims
}

func (c *Claims) Valid() error {
	if _, err := strconv.ParseInt(c.Subject, 16, 64); err != nil {
		return jwt.NewValidationError("invalid subject", 0)
	}
	return c.StandardClaims.Valid()
}

func (h *handler) TrustedCheckIn(
	ctx context.Context,
	req *rvapi.TrustedCheckInRequest,
) (*rvapi.TrustedCheckInResponse, error) {
	var claims *Claims
	err := h.txM.inTx(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		client, err := models.Clients(models.ClientWhere.Name.EQ(req.ClientId)).One(ctx, tx)
		if errors.Is(err, sql.ErrNoRows) {
			return status.Error(codes.InvalidArgument, "Unknown client id")
		} else if err != nil {
			return err
		} else if !rvapi.Valid(req, time.Now(), client.Secret) {
			return status.Error(codes.Unauthenticated, "Invalid message")
		}

		alias, err := client.Aliases(models.AliasWhere.Alias.EQ(req.UserName)).One(ctx, tx)
		if notFound := errors.Is(err, sql.ErrNoRows); err != nil && !notFound {
			return err
		} else if notFound {
			var user models.User
			if err = user.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}

			alias = &models.Alias{ClientID: client.ID, UserID: user.ID, Alias: req.UserName}
			if err = alias.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}
		}

		claims = newClaims(alias.UserID)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := h.sendUpdatedClaims(ctx, claims); err != nil {
		return nil, err
	}
	return &rvapi.TrustedCheckInResponse{}, nil
}

func (h *handler) CheckIn(ctx context.Context, _ *rvapi.CheckInRequest) (*rvapi.CheckInResponse, error) {
	claims := loadClaims(ctx)
	if claims == nil {
		var u models.User
		if err := h.txM.inTx(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
			return u.Insert(ctx, tx, boil.Infer())
		}); err != nil {
			return nil, err
		}
		claims = newClaims(u.ID)
	}
	if err := h.sendUpdatedClaims(ctx, claims); err != nil {
		return nil, err
	}
	return &rvapi.CheckInResponse{}, nil
}

func (h *handler) sendUpdatedClaims(ctx context.Context, claims *Claims) error {
	now := time.Now().Truncate(time.Second)

	exp := now.Add(h.tokLife)

	claims.IssuedAt = now.Unix()
	claims.NotBefore = now.Unix()
	claims.ExpiresAt = exp.Unix()

	tok, err := h.tokM.write(claims)
	if err != nil {
		return err
	}

	// returns immediately if not a grpc web proxy req
	if err := setIdentityTokenCookie(ctx, &http.Cookie{
		Name:     "rv-token",
		Value:    tok,
		Path:     "/",
		Expires:  exp,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}); err != nil {
		return err
	}

	return grpc.SendHeader(ctx, metadata.Pairs("rv-token", tok))
}

type cookieRequest struct {
	cookie *http.Cookie
	done   chan<- struct{}
}

func tokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("rv-token"); err == nil {
			r = r.WithContext(context.WithValue(r.Context(), parsedCookieKey, cookie.Value))
		}

		ch := make(chan *cookieRequest)
		r = r.WithContext(context.WithValue(r.Context(), cookieRequestKey, ch))

		tokenCtx, cncl := context.WithCancel(r.Context())
		defer cncl()

		go func() {
			var cookieR *cookieRequest
			select {
			case <-tokenCtx.Done():
				return
			case cookieR = <-ch:
			}
			http.SetCookie(w, cookieR.cookie)
			select {
			case <-tokenCtx.Done():
			case cookieR.done <- struct{}{}:
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func setIdentityTokenCookie(ctx context.Context, cookie *http.Cookie) error {
	ch, _ := ctx.Value(cookieRequestKey).(chan *cookieRequest)
	if ch == nil {
		return nil
	}
	done := make(chan struct{})
	select {
	case ch <- &cookieRequest{cookie, done}:
	case <-ctx.Done():
		return ctx.Err()
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
