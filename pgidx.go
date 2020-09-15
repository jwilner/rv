package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

	e := newElection(
		&models.Election{
			Key:       h.kGen.newKey(keyCharSet, 8),
			BallotKey: h.kGen.newKey(keyCharSet, 8),
			Question:  i.Question,
			Choices:   choices.raw(),
		},
	)

	if err := h.txM.inTx(r.Context(), nil, func(ctx context.Context, tx *sql.Tx) error {
		return e.Insert(ctx, tx, boil.Infer())
	}); err != nil {
		handleError(w, r, err)
		return
	}

	http.Redirect(w, r, "/e/"+e.Key, http.StatusFound)
}

func parseChoices(s string) normalizedSlice {
	reader := csv.NewReader(strings.NewReader(s))
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		log.Printf("csv.Reader.read; %v\n", err)
	}

	var combined []string
	for _, row := range rows {
		for _, c := range row {
			if c := strings.TrimSpace(c); c != "" {
				combined = append(combined, c)
			}
		}
	}
	return normalize(combined)
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
