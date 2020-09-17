package platform

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jwilner/rv/pkg/pb/rvapi"

	"github.com/jackc/pgtype"

	"github.com/jackc/pgconn"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jwilner/rv/internal/models"
)

func (h *handler) Vote(ctx context.Context, req *rvapi.VoteRequest) (*rvapi.VoteResponse, error) {
	norm, details := validateChoices(req)
	if req.Name == "" {
		details = append(
			details,
			&errdetails.BadRequest_FieldViolation{Field: "Name", Description: "must not be empty"},
		)
	}
	if len(details) > 0 {
		s, err := status.New(codes.InvalidArgument, "Invalid vote").WithDetails(details...)
		if err != nil {
			panic(fmt.Sprintf("unexpected error: %v", err))
		}
		return nil, s.Err()
	}

	err := h.txM.inTx(ctx, nil, func(ctx context.Context, tx *sql.Tx) (err error) {
		var el *models.Election
		if el, err = models.Elections(models.ElectionWhere.BallotKey.EQ(req.BallotKey)).One(ctx, tx); err != nil {
			return fmt.Errorf("model.Elections: %w", err)
		}
		if el.Close.Status == pgtype.Present && !el.Close.Time.After(time.Now()) {
			return fmt.Errorf("election has already closed: %v %v", el.Close.Time, time.Now())
		}
		if undefined := difference(norm, normalize(el.Choices)); len(undefined) > 0 {
			return fmt.Errorf("unknown choices: %v", strings.Join(undefined.raw(), ", "))
		}
		v := models.Vote{
			ElectionID: el.ID,
			Name:       req.Name,
			Choices:    norm.raw(),
		}
		_ = v.CreatedAt.Set(time.Now().UTC())

		return v.Insert(ctx, tx, boil.Infer())
	})
	if isDupe(err, "vote_election_id_name_key") {
		return nil, status.Error(codes.AlreadyExists, "name already used")
	}
	if err != nil {
		return nil, err
	}
	return &rvapi.VoteResponse{}, nil
}

func isDupe(err error, constraintName string) bool {
	pgErr := new(pgconn.PgError)
	return errors.As(err, &pgErr) &&
		pgErr.Code == "23505" &&
		pgErr.SchemaName == "rv" &&
		pgErr.ConstraintName == constraintName
}
