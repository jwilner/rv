package platform

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/jwilner/rv/internal/models"
	"github.com/jwilner/rv/pkg/pb/rvapi"
)

func (h *handler) ListViews(ctx context.Context, req *rvapi.ListViewsRequest) (*rvapi.ListViewsResponse, error) {
	var op func(context.Context, boil.ContextExecutor) ([]*models.Election, error)
	switch req.Filter {
	case rvapi.ListViewsRequest_PUBLIC:
		op = listPublic
	case rvapi.ListViewsRequest_VOTED_IN:
		op = func(ctx context.Context, exec boil.ContextExecutor) ([]*models.Election, error) {
			return listVotedIn(ctx, exec, userID(ctx))
		}
	default:
		return nil, invalidArgument("filter")
	}

	var els []*models.Election
	err := h.txM.inTx(ctx, &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		els, err = op(ctx, tx)
		return
	})
	if err != nil {
		return nil, err
	}

	resp := rvapi.ListViewsResponse{Elections: make([]*rvapi.ElectionView, 0, len(els))}
	for _, e := range els {
		resp.Elections = append(resp.Elections, protoElectionView(e))
	}
	return &resp, nil
}

func (h *handler) List(ctx context.Context, _ *rvapi.ListRequest) (*rvapi.ListResponse, error) {
	var els []*models.Election
	err := h.txM.inTx(ctx, &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		els, err = loadUserElections(ctx, tx, userID(ctx))
		return
	})
	if err != nil {
		return nil, err
	}
	resp := rvapi.ListResponse{Elections: make([]*rvapi.Election, 0, len(els))}
	for _, e := range els {
		resp.Elections = append(resp.Elections, protoElection(e))
	}
	return &resp, nil

}

func validateChoices(
	req interface{ GetChoices() []string },
	permitZero bool,
) (norm normalizedSlice, details []*errdetails.BadRequest_FieldViolation) {
	if norm = normalize(req.GetChoices()); len(norm) == 0 && !permitZero {
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

func listVotedIn(ctx context.Context, exec boil.ContextExecutor, userID int64) (els []*models.Election, err error) {
	err = queries.Raw(
		`
SELECT 
	e.*
FROM
	rv.election e
JOIN
	rv.vote v
ON
	(v.election_id = e.id)
WHERE
	v.user_id = $1
`,
		userID,
	).Bind(ctx, exec, &els)
	return
}

func loadUserElections(ctx context.Context, exec boil.ContextExecutor, userID int64) ([]*models.Election, error) {
	return models.Elections(
		models.ElectionWhere.UserID.EQ(userID),
		qm.OrderBy(models.ElectionColumns.CreatedAt+" DESC"),
		qm.Limit(10),
	).All(ctx, exec)
}

func listPublic(
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
