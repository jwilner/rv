package platform

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	_ "github.com/jackc/pgx/v4/stdlib" // register driver
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/jwilner/rv/internal/slack"
	"github.com/jwilner/rv/pkg/pb/rvapi"
)

type Config struct {
	Debug                                        bool
	DBURL, Addr, GRPCAddr, StaticDir, SigningKey string
	TokenLength                                  time.Duration
	Slack                                        slack.Config
}

// Run runs the application, connecting to the database at dbURL and listening for HTTP at the provided address.
func Run(config *Config) error {
	// construct a ctx that we can cancel
	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()

	ctx = setDebug(ctx, config.Debug)

	// cancel it if interrupted
	cancelOnSignal(ctx, cncl, os.Interrupt)

	tokM := &tokenManager{signingKey: []byte(config.SigningKey)}

	h, err := newHandler(ctx, tokM, config.DBURL, config.TokenLength)
	if err != nil {
		return err
	}
	defer func() {
		if err := h.Close(); err != nil {
			log.Printf("error while closing handler: %v", err)
		}
	}()

	var grpcL net.Listener
	if config.GRPCAddr != "" {
		if grpcL, err = net.Listen("tcp", config.GRPCAddr); err != nil {
			return err
		}
		defer func() {
			_ = grpcL.Close()
		}()
	}

	var staticMux http.Handler
	if config.StaticDir != "" {
		staticMux = http.FileServer(http.Dir(config.StaticDir))
	}

	var (
		slackMux  http.Handler
		grpcUnixL net.Listener
	)
	if config.Slack.Token != "" {
		dir, err := ioutil.TempDir("", "rv*")
		if err != nil {
			return err
		}
		defer func() {
			_ = os.RemoveAll(dir)
		}()

		sockPath := path.Join(dir, "rv.sock")

		if grpcUnixL, err = net.Listen("unix", sockPath); err != nil {
			return err
		}
		defer func() {
			_ = grpcUnixL.Close()
		}()

		conn, err := grpc.Dial("passthrough:///unix://"+sockPath, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer func() {
			_ = conn.Close()
		}()

		slackMux = slack.New(&config.Slack, conn)
	}

	server := newGRPCServer(tokM, h)

	mux := http.Handler(buildMux(server, staticMux, slackMux))
	if config.Debug {
		mux = logMiddleware(mux)
	}
	mux = tokenMiddleware(mux)
	mux = requestIDer(ctx, mux)

	errC := make(chan error, 3)

	var wg sync.WaitGroup

	if grpcL != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case errC <- listenAndServeGRPC(ctx, server, grpcL):
			case <-ctx.Done():
				errC <- ctx.Err()
			}
		}()
	}

	if grpcUnixL != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case errC <- listenAndServeGRPC(ctx, server, grpcUnixL):
			case <-ctx.Done():
				errC <- ctx.Err()
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case errC <- listenAndServe(ctx, config.Addr, mux):
		case <-ctx.Done():
			select {
			case errC <- ctx.Err():
			default:
			}
		}
	}()

	// return first error and wait for all to complete
	defer wg.Wait()
	return <-errC
}

func buildMux(server *grpc.Server, staticMux, slackMux http.Handler) *http.ServeMux {
	grpcWebServer := grpcweb.WrapServer(server, grpcweb.WithWebsockets(true))

	mux := http.NewServeMux()

	if slackMux != nil {
		mux.Handle("/api/slack/", http.StripPrefix("/api/slack", slackMux))
	}

	mux.Handle("/api/", http.StripPrefix("/api/", grpcWebServer))

	if staticMux != nil {
		mux.Handle("/static/", staticMux)

		// any requests that come in at the following paths should be rewritten to be served
		// the index and the frontend router.
		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			original := r.URL.Path

			r.URL.Path = "/"
			staticMux.ServeHTTP(w, r)

			r.URL.Path = original
		}))
	}
	return mux
}

func newHandler(
	ctx context.Context,
	tokM *tokenManager,
	dbURL string,
	tokenLife time.Duration,
) (*handler, error) {
	h := handler{kGen: newStringGener(), tokM: tokM, tokLife: tokenLife}
	db, err := connectDB(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %w", err)
	}

	h.txM = &txMgr{db}
	return &h, nil
}

type codeInterceptor struct {
	http.ResponseWriter

	codeSeen int
}

// necessary for grpc web proxy logic
func (c *codeInterceptor) Flush() {
	if f, ok := c.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (c *codeInterceptor) WriteHeader(statusCode int) {
	c.codeSeen = statusCode
	c.ResponseWriter.WriteHeader(statusCode)
}

var _ http.Flusher = new(codeInterceptor)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := codeInterceptor{w, http.StatusOK}
		t := time.Now()
		next.ServeHTTP(&c, r)
		elapsed := time.Since(t)
		log.Printf(
			"request_id=%v method=%v url=%q response_code=%v elapsed_ms=%v",
			requestID(r.Context()),
			r.Method,
			r.URL,
			c.codeSeen,
			int64(elapsed/time.Millisecond),
		)
	})
}

func newGRPCServer(tokM *tokenManager, h *handler) *grpc.Server {
	srvr := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			tokenInterceptor(tokM),
			interceptor,
		),
	)
	rvapi.RegisterRVerServer(srvr, h)
	reflection.Register(srvr)
	return srvr
}

func interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	t := time.Now()
	resp, err = handler(ctx, req)
	elapsed := int64(time.Since(t) / time.Millisecond)

	if err == nil {
		if isDebug(ctx) {
			log.Printf(
				`request_id=%v user_id=%v rpc=%v elapsed_ms=%v msg="rpc success"`,
				requestID(ctx),
				userID(ctx),
				info.FullMethod,
				elapsed,
			)
		}
		return
	}

	if s, ok := status.FromError(err); ok {
		log.Printf(
			`request_id=%v user_id=%v rpc=%v code=%d err=%q elapsed_ms=%v msg="rpc failure"`,
			requestID(ctx),
			userID(ctx),
			info.FullMethod,
			s.Code(),
			s.Message(),
			elapsed,
		)
		return
	}

	// make a new error
	var s *status.Status
	switch {
	case errors.Is(err, sql.ErrNoRows):
		s = status.New(codes.NotFound, "resource not found")
	default:
		s = status.New(codes.Internal, "internal server error")
	}
	if isDebug(ctx) {
		s, _ = s.WithDetails(&errdetails.DebugInfo{Detail: err.Error()})
	}
	log.Printf(
		`request_id=%v user_id=%v rpc=%v code=%d err=%q original_err=%q elapsed_ms=%v msg="rpc failure"`,
		requestID(ctx),
		userID(ctx),
		info.FullMethod,
		s.Code(),
		s.Message(),
		err,
		elapsed,
	)
	err = s.Err()
	return
}

// cancelOnSignal cancels the provided context when any of the signals is received
func cancelOnSignal(ctx context.Context, cncl context.CancelFunc, sigs ...os.Signal) {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, sigs...)

	// cancel the context if either a os.Interrupt arrives
	go func() {
		defer cncl()
		select {
		case <-sigC:
		case <-ctx.Done():
		}
	}()
}

func connectDB(ctx context.Context, url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := func() error {
		ctx, cncl := context.WithTimeout(ctx, 10*time.Second)
		defer cncl()
		return db.PingContext(ctx)
	}(); err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	return db, nil
}

func listenAndServeGRPC(ctx context.Context, server *grpc.Server, l net.Listener) error {
	errs := make(chan error, 1)
	go func() {
		errs <- server.Serve(l)
	}()
	select {
	case <-ctx.Done():
		done := make(chan struct{}, 1)
		go func() {
			server.GracefulStop()
			done <- struct{}{}
		}()
		select {
		case <-time.After(5 * time.Second):
			server.Stop()
		case <-done:
		}
		return ctx.Err()
	case err := <-errs:
		return err
	}
}

func listenAndServe(ctx context.Context, addr string, handler http.Handler) error {
	srvr := http.Server{
		Addr:        addr,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Handler:     handler,
	}

	errs := make(chan error, 1)
	go func() {
		errs <- srvr.ListenAndServe()
	}()

	select {
	case err := <-errs: // server shutdown on its own
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("srvr.ListenAndServe: %w", err)

	case <-ctx.Done():
		ctx, cncl := context.WithTimeout(context.Background(), time.Second*5)
		defer cncl()

		log.Println("Beginning shutdown.")

		cancelOnSignal(ctx, cncl, os.Interrupt) // if you interrupt again, we'll shutdown immediately

		if err := srvr.Shutdown(ctx); err != nil {
			return fmt.Errorf("srvr.Shutdown: %w", err)
		}
		return nil
	}
}

func requestIDer(ctx context.Context, next http.Handler) http.Handler {
	// we'll have some buffered request ids ready to go if a proxy isn't setting request ids for us.
	// wrap in a routine because rnd isn't safe.
	requestIDs := make(chan string, 10)
	go func() {
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		buf := make([]byte, 20)
		chrSet := []byte("abcdefghijklmnopqrstuvwxyz012345679")
		l := uint8(len(chrSet))
		for {
			_, _ = rnd.Read(buf) // cannot fail
			for i := range buf {
				buf[i] = chrSet[buf[i]%l]
			}
			select {
			case requestIDs <- string(buf):
			case <-ctx.Done():
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			select {
			case <-r.Context().Done():
			case reqID = <-requestIDs:
			}
		}
		r = r.WithContext(context.WithValue(r.Context(), requestIDKey, reqID))
		next.ServeHTTP(w, r)
	})
}
