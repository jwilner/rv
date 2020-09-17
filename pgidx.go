package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/jackc/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jwilner/rv/pkg/pb/rvapi"

	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jwilner/rv/models"
)

type indexPage struct {
	Form indexForm

	Now      time.Time
	Overview []*election
}

type indexForm struct {
	form

	Question string
	Choices  string
}

func (i *indexForm) unmarshal(vals url.Values) {
	i.Question = vals.Get("question")
	i.Choices = vals.Get("choices")
}

func (i *indexForm) validate() (normalizedSlice, bool) {
	if len(i.Question) == 0 {
		i.setErrorf("Question", "must provide a question")
	}
	choices := parseChoices(i.Choices)
	if len(choices) == 0 {
		i.setErrorf("Choices", "must have at least one choice")
	} else {
		validateNonDupe(choices, &i.form)
	}
	return choices, i.checkErrors()
}

func (h *handler) getIndex(w http.ResponseWriter, r *http.Request) {
	var overview []*election
	err := h.txM.inTx(r.Context(), &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		overview, err = loadElectionOverview(ctx, tx)
		return
	})
	if handleError(w, r, err) {
		return
	}
	h.tmpls.render(r.Context(), w, "index.html", &indexPage{Overview: overview, Now: time.Now().UTC()})
}

func (h *handler) Create(ctx context.Context, request *rvapi.CreateRequest) (*rvapi.CreateResponse, error) {
	norm, err := grpcValidate(request)
	if err != nil {
		return nil, err
	}
	el, err := h.create(ctx, request.Question, norm)
	if err != nil {
		return nil, err
	}
	return &rvapi.CreateResponse{Election: el.proto()}, nil
}

func (h *handler) Overview(ctx context.Context, _ *rvapi.OverviewRequest) (*rvapi.OverviewResponse, error) {
	var els []*election
	err := h.txM.inTx(ctx, &sql.TxOptions{ReadOnly: true}, func(ctx context.Context, tx *sql.Tx) (err error) {
		els, err = loadElectionOverview(ctx, tx)
		return
	})
	if err != nil {
		return nil, err
	}

	resp := rvapi.OverviewResponse{Elections: make([]*rvapi.Election, 0, len(els))}
	for _, el := range els {
		resp.Elections = append(resp.Elections, el.proto())
	}

	return &resp, nil
}

func (e *election) proto() *rvapi.Election {
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

func (h *handler) postIndex(w http.ResponseWriter, r *http.Request) {
	if handleError(w, r, r.ParseForm()) {
		return
	}

	var i indexForm

	i.unmarshal(r.PostForm)

	choices, ok := i.validate()
	if !ok {
		log.Printf("validation failed: %v", i.Errors)
		var overview []*election
		err := h.txM.inTx(
			r.Context(),
			&sql.TxOptions{ReadOnly: true},
			func(ctx context.Context, tx *sql.Tx) (err error) {
				overview, err = loadElectionOverview(ctx, tx)
				return
			},
		)
		if handleError(w, r, err) {
			return
		}
		h.tmpls.render(r.Context(), w, "index.html", &indexPage{Form: i, Overview: overview, Now: time.Now().UTC()})
		return
	}

	e, err := h.create(r.Context(), i.Question, choices)
	if handleError(w, r, err) {
		return
	}

	http.Redirect(w, r, "/e/"+e.Key, http.StatusFound)
}

func (h *handler) create(ctx context.Context, question string, choices normalizedSlice) (*election, error) {
	e := newElection(
		&models.Election{
			Key:       h.kGen.newKey(keyCharSet, 8),
			BallotKey: h.kGen.newKey(keyCharSet, 8),
			Question:  question,
			Choices:   choices.raw(),
		},
	)
	err := h.txM.inTx(ctx, nil, func(ctx context.Context, tx *sql.Tx) error {
		return e.Insert(ctx, tx, boil.Infer())
	})
	return e, err
}

func parseChoices(s string) normalizedSlice {
	var raw []string

	for _, c := range strings.Split(s, "\n") {
		if c := strings.TrimSpace(c); c != "" {
			raw = append(raw, c)
		}
	}

	return normalize(raw)
}

func loadElectionOverview(ctx context.Context, exec boil.ContextExecutor) ([]*election, error) {
	var ms []*models.Election
	err := queries.Raw(
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
	if err != nil {
		return nil, err
	}
	els := make([]*election, 0, len(ms))
	for _, m := range ms {
		els = append(els, newElection(m))
	}
	return els, nil
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

func validateNonDupe(choices normalizedSlice, f *form) {
	counts := make(map[string]int)
	for _, c := range choices {
		counts[c.normalized]++
	}

	for _, c := range choices {
		if counts[c.normalized] > 1 {
			f.setErrorf("Choices", "%q occurs more %d times -- can only occur once", c.raw, counts[c.normalized])
		}
		counts[c.normalized] = 0
	}
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
