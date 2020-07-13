package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"net/http"
	"net/url"
	"strings"
)

type indexPage struct {
	Form indexForm
}

type indexForm struct {
	form

	Name    string
	Choices string
}

func (i *indexForm) unmarshal(vals url.Values) {
	i.Name = vals.Get("name")
	i.Choices = vals.Get("choices")
}

func (i *indexForm) validate() bool {
	if len(i.Name) == 0 {
		i.setErrorf("Name", "must provide a name")
	}
	if choices := parseChoices(i.Choices); len(choices) == 0 {
		i.setErrorf("Choices", "must have at least one choice")
	} else {
		counts := make(map[string]int)
		for _, c := range choices {
			counts[c]++
		}

		for _, c := range choices {
			if counts[c] > 1 {
				i.setErrorf("Choices", "%q occurs more %d times -- can only occur once", c, counts[c])
			}
			counts[c] = 0
		}
	}
	return i.checkErrors()
}

func (h *handler) getIndex(w http.ResponseWriter, r *http.Request) {
	h.tmpls.render(r.Context(), w, "index.html", &indexPage{})
}

func (h *handler) postIndex(w http.ResponseWriter, r *http.Request) {
	if handleError(w, r, r.ParseForm()) {
		return
	}

	var i indexForm

	i.unmarshal(r.PostForm)

	if !i.validate() {
		h.tmpls.render(r.Context(), w, "index.html", &indexPage{Form: i})
		return
	}

	e := election{
		Name:    i.Name,
		Choices: parseChoices(i.Choices),
	}
	if err := h.txM.inTx(r.Context(), nil, func(ctx context.Context, tx *sql.Tx) error {
		return insertElection(ctx, tx, &e)
	}); err != nil {
		handleError(w, r, err)
		return
	}
}

func parseChoices(s string) []string {
	row, _ := csv.NewReader(strings.NewReader(s)).Read()
	return row
}

func insertElection(ctx context.Context, tx *sql.Tx, e *election) error {
	return tx.QueryRowContext(
		ctx,
		`
INSERT INTO 
    rv.election	(name, choices) VALUES ($1, $2) 
RETURNING
	key,
    created_at
	`,
		e.Name,
		e.Choices,
	).Scan(
		&e.Key,
		&e.CreatedAt,
	)
}
