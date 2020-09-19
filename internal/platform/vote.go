package platform

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

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
	norm, details := validateChoices(req, true)
	if req.Name == "" {
		details = append(
			details,
			&errdetails.BadRequest_FieldViolation{Field: "Name", Description: "must not be empty"},
		)
	}
	if len(details) > 0 {
		return nil, detailedErr(
			codes.InvalidArgument,
			"invalid fields",
			&errdetails.BadRequest{
				FieldViolations: details,
			},
		)
	}

	err := h.txM.inTx(ctx, nil, func(ctx context.Context, tx *sql.Tx) (err error) {
		var el *models.Election
		if el, err = models.Elections(models.ElectionWhere.BallotKey.EQ(req.BallotKey)).One(ctx, tx); err != nil {
			return detailedErr(codes.NotFound, "election not found", &errdetails.ResourceInfo{
				ResourceType: "election",
				ResourceName: req.BallotKey,
				Description:  fmt.Sprintf("no election with ballot key %v", req.BallotKey),
			})
		}
		if el.Close.Status == pgtype.Present && !el.Close.Time.After(time.Now()) {
			return detailedErr(codes.FailedPrecondition, "election has already closed")
		}
		if undefined := difference(norm, normalize(el.Choices)); len(undefined) > 0 {
			return invalidArgument("unknown choices", "Choices", "unknown choices")
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
		return nil, invalidArgument("value already used", "Name", "already used")
	}
	if err != nil {
		return nil, err
	}
	return &rvapi.VoteResponse{}, nil
}

func invalidArgument(message string, fieldDescriptions ...string) error {
	descs := make([]*errdetails.BadRequest_FieldViolation, 0, len(fieldDescriptions)/2)
	for i := 0; i < len(fieldDescriptions); i += 2 {
		descs = append(descs, &errdetails.BadRequest_FieldViolation{
			Field:       fieldDescriptions[0],
			Description: fieldDescriptions[1],
		})
	}
	return detailedErr(codes.InvalidArgument, message, &errdetails.BadRequest{FieldViolations: descs})
}

func detailedErr(code codes.Code, msg string, deets ...proto.Message) error {
	s, _ := status.New(code, msg).WithDetails(deets...)
	return s.Err()
}

func isDupe(err error, constraintName string) bool {
	pgErr := new(pgconn.PgError)
	return errors.As(err, &pgErr) &&
		pgErr.Code == "23505" &&
		pgErr.SchemaName == "rv" &&
		pgErr.ConstraintName == constraintName
}
