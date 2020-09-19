package platform

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jackc/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jwilner/rv/pkg/pb/rvapi"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jwilner/rv/internal/models"
)

func (h *handler) Create(ctx context.Context, request *rvapi.CreateRequest) (*rvapi.CreateResponse, error) {
	norm, details := validateChoices(request, false)
	if len(request.Question) == 0 {
		details = append(details, &errdetails.BadRequest_FieldViolation{
			Field:       "Question",
			Description: "Cannot be empty",
		})
	}
	if len(details) > 0 {
		return nil, detailedErr(
			codes.InvalidArgument,
			"invalid create request",
			&errdetails.BadRequest{FieldViolations: details},
		)
	}

	e := &models.Election{
		Key:       h.kGen.newKey(keyCharSet, 8),
		BallotKey: h.kGen.newKey(keyCharSet, 8),
		Question:  request.Question,
		Choices:   norm.raw(),
		UserID:    userID(ctx),
	}
	_ = e.Close.Set(nil) // set null
	_ = e.CreatedAt.Set(time.Now())

	err := h.txM.inTx(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		return e.Insert(ctx, tx, boil.Infer())
	})
	if err != nil {
		return nil, err
	}
	return &rvapi.CreateResponse{Election: protoElection(e)}, nil
}

func (h *handler) Get(ctx context.Context, req *rvapi.GetRequest) (*rvapi.GetResponse, error) {
	var el *models.Election
	err := h.txM.inTx(ctx, &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		if el, err = models.Elections(models.ElectionWhere.Key.EQ(req.Key)).One(ctx, tx); err != nil {
			return fmt.Errorf("models.Elections key=%v: %w", req.Key, err)
		}
		return
	})
	if errors.Is(err, sql.ErrNoRows) {
		s, err := status.New(codes.NotFound, "election not found").WithDetails(
			&errdetails.ResourceInfo{ResourceType: "Election", ResourceName: req.Key},
		)
		if err != nil {
			panic(fmt.Sprintf("impossible outcome: %v", err))
		}
		err = s.Err()
	}
	if err != nil {
		return nil, err
	}
	return &rvapi.GetResponse{Election: protoElection(el)}, nil
}

func (h *handler) GetView(ctx context.Context, req *rvapi.GetViewRequest) (*rvapi.GetViewResponse, error) {
	var el *models.Election
	err := h.txM.inTx(ctx, &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		if el, err = models.Elections(models.ElectionWhere.BallotKey.EQ(req.BallotKey)).One(ctx, tx); err != nil {
			return fmt.Errorf("models.Elections ballotKey=%v: %w", req.BallotKey, err)
		}
		return
	})
	if errors.Is(err, sql.ErrNoRows) {
		s, err := status.New(codes.NotFound, "election not found").WithDetails(
			&errdetails.ResourceInfo{ResourceType: "Election", ResourceName: req.BallotKey},
		)
		if err != nil {
			panic(fmt.Sprintf("impossible outcome: %v", err))
		}
		err = s.Err()
	}
	if err != nil {
		return nil, err
	}
	return &rvapi.GetViewResponse{Election: protoElectionView(el)}, nil
}

func (h *handler) Update(ctx context.Context, req *rvapi.UpdateRequest) (*rvapi.UpdateResponse, error) {
	var el *models.Election
	err := h.txM.inTx(ctx, nil, func(ctx context.Context, tx *sql.Tx) (err error) {
		if el, err = models.Elections(models.ElectionWhere.Key.EQ(req.Key)).One(ctx, tx); err != nil {
			return fmt.Errorf("models.Elections key=%v: %w", req.Key, err)
		}
		var modifiedColumns []string
		addToSet := func(c string) {
			modifiedColumns = add(modifiedColumns, c)
		}
		for _, op := range req.Operations {
			switch op := op.Operation.(type) {
			case *rvapi.UpdateRequest_Operation_SetClose:
				addToSet(models.ElectionColumns.Close)
				if op.SetClose.Close == nil {
					_ = el.Close.Set(nil) // unset
					break
				}
				_ = el.Close.Set(op.SetClose.Close.AsTime())
			case *rvapi.UpdateRequest_Operation_ModifyFlags:
				addToSet(models.ElectionColumns.Flags)
				for _, addFlag := range op.ModifyFlags.Add {
					if flag := mapFlag(addFlag); flag != "" {
						el.Flags = add(el.Flags, flag)
					}
				}
				for _, remFlag := range op.ModifyFlags.Remove {
					if flag := mapFlag(remFlag); flag != "" {
						last := remove(el.Flags, flag)
						el.Flags = el.Flags[:last]
					}
				}
			}
		}
		if len(modifiedColumns) > 0 {
			_, err = el.Update(ctx, tx, boil.Whitelist(modifiedColumns...))
		}
		return
	})
	if err != nil {
		return nil, err
	}
	return &rvapi.UpdateResponse{Election: protoElection(el)}, nil
}

func mapFlag(flag rvapi.Election_Flag) string {
	switch flag {
	case rvapi.Election_PUBLIC:
		return electionFlagPublic
	case rvapi.Election_RESULTS_HIDDEN:
		return electionFlagResultsHidden
	default:
		return ""
	}
}

// Elections are either:
// - open (close NULL)
// - scheduled (close a time in the future)
// - closed (close a time not in the future)
// Updates:
// - schedule close: open or scheduled -> scheduled (set a close time in future)
// - close now: open -> closed (set close to now)
// - unset: closed -> open (set close to null)
// - unset: scheduled -> open (set close to null)

const (
	electionFlagPublic        = "public"
	electionFlagResultsHidden = "results_hidden"
)

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func remove(haystack []string, needle string) int {
	cur := 0
	for i, v := range haystack {
		if haystack[i] == needle {
			continue
		}
		haystack[cur] = v
		cur++
	}
	return cur
}

func add(haystack []string, needle string) []string {
	if contains(haystack, needle) {
		return haystack
	}
	return append(haystack, needle)
}

func protoElectionView(e *models.Election) *rvapi.ElectionView {
	el := rvapi.ElectionView{
		Question:  e.Question,
		Choices:   e.Choices,
		BallotKey: e.BallotKey,
	}

	if e.Close.Status == pgtype.Present {
		var t time.Time
		_ = e.Close.AssignTo(&t)
		if closeTS, err := ptypes.TimestampProto(t); err == nil {
			el.Close = closeTS
		}
	}

	el.Flags = make([]rvapi.Election_Flag, 0, len(e.Flags))
	for _, f := range e.Flags {
		var fl rvapi.Election_Flag
		switch f {
		case electionFlagPublic:
			fl = rvapi.Election_PUBLIC
		case electionFlagResultsHidden:
			fl = rvapi.Election_RESULTS_HIDDEN
		}
		el.Flags = append(el.Flags, fl)
	}

	return &el
}

func protoElection(e *models.Election) *rvapi.Election {
	el := rvapi.Election{
		Question:  e.Question,
		Choices:   e.Choices,
		Key:       e.Key,
		BallotKey: e.BallotKey,
	}

	if e.Close.Status == pgtype.Present {
		var t time.Time
		_ = e.Close.AssignTo(&t)
		if closeTS, err := ptypes.TimestampProto(t); err == nil {
			el.Close = closeTS
		}
	}

	el.Flags = make([]rvapi.Election_Flag, 0, len(e.Flags))
	for _, f := range e.Flags {
		var fl rvapi.Election_Flag
		switch f {
		case electionFlagPublic:
			fl = rvapi.Election_PUBLIC
		case electionFlagResultsHidden:
			fl = rvapi.Election_RESULTS_HIDDEN
		}
		el.Flags = append(el.Flags, fl)
	}

	return &el
}
