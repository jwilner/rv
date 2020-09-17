package platform

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jwilner/rv/pkg/pb/rvapi"

	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jwilner/rv/internal/models"
)

func (h *handler) Overview(ctx context.Context, _ *rvapi.OverviewRequest) (*rvapi.OverviewResponse, error) {
	var els []*models.Election
	err := h.txM.inTx(ctx, &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		els, err = loadElectionOverview(ctx, tx)
		return
	})
	if err != nil {
		return nil, err
	}

	resp := rvapi.OverviewResponse{Elections: make([]*rvapi.Election, 0, len(els))}
	for _, el := range els {
		resp.Elections = append(resp.Elections, protoElection(el))
	}

	return &resp, nil
}

func validateChoices(
	req interface{ GetChoices() []string },
) (norm normalizedSlice, details []proto.Message) {
	if norm = normalize(req.GetChoices()); len(norm) == 0 {
		details = append(details, &errdetails.BadRequest_FieldViolation{
			Field:       "Choices",
			Description: "Cannot be empty",
		})
	}

	counts := make(map[string]int)
	for _, c := range norm {
		counts[c.normalized]++
	}

	for _, c := range norm {
		if count := counts[c.normalized]; count > 1 {
			details = append(details, &errdetails.BadRequest_FieldViolation{
				Field:       "Choices",
				Description: fmt.Sprintf("%q occurs %d times", c.raw, count),
			})
		}
		delete(counts, c.normalized)
	}
	return
}

func grpcValidate(req *rvapi.CreateRequest) (normalizedSlice, error) {
	norm, details := validateChoices(req)

	if len(req.Question) == 0 {
		details = append(details, &errdetails.BadRequest_FieldViolation{
			Field:       "Question",
			Description: "Cannot be empty",
		})
	}

	if len(details) > 0 {
		s, err := status.New(codes.InvalidArgument, "Invalid create request").WithDetails(details...)
		if err != nil {
			panic(fmt.Sprintf("failed to construct proper err: %v", err))
		}
		return nil, s.Err()
	}

	return norm, nil
}

func loadElectionOverview(
	ctx context.Context,
	exec boil.ContextExecutor,
) (ms []*models.Election, err error) {
	err = queries.Raw(
		`
SELECT 
	*,
	coalesce(close, 'infinity') > NOW() AS active,
	CASE
		WHEN close - NOW() > INTERVAL '0' THEN close - NOW()
		WHEN close - NOW() <= INTERVAL '0' THEN -(close - NOW())
		ELSE NULL
	END AS distance
FROM rv.election e
WHERE 
	'public' = ANY(e.flags)
ORDER BY
	active DESC,
	distance ASC NULLS LAST
LIMIT
	10
`).Bind(
		ctx,
		exec,
		&ms,
	)
	return
}

func normalize(raw []string) normalizedSlice {
	n := make(normalizedSlice, 0, len(raw))
	for _, r := range raw {
		n = append(n, &struct {
			raw, normalized string
		}{r, strings.ToLower(r)})
	}
	return n
}

type normalizedSlice []*struct {
	raw, normalized string
}

func (n normalizedSlice) raw() []string {
	o := make([]string, 0, len(n))
	for _, c := range n {
		o = append(o, c.raw)
	}
	return o
}

// return, in original order, every element in left that is not in right -- case insensitive
func difference(left, right normalizedSlice) (diff normalizedSlice) {
	present := make(map[string]bool, len(right))
	for _, r := range right {
		present[r.normalized] = true
	}
	for _, l := range left {
		if !present[l.normalized] {
			diff = append(diff, l)
		}
	}
	return
}
