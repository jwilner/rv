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

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/jwilner/rv/models"
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
		log.Printf("validation failed: %v", i.Errors)
		h.tmpls.render(r.Context(), w, "index.html", &indexPage{Form: i})
		return
	}

	e := models.Election{
		Key:       h.kGen.newKey(keyCharSet, 8),
		BallotKey: h.kGen.newKey(keyCharSet, 8),
		Name:      i.Name,
		CreatedAt: time.Now().UTC(),
		Choices:   parseChoices(i.Choices),
	}
	if err := h.txM.inTx(r.Context(), nil, func(ctx context.Context, tx *sql.Tx) error {
		return e.Insert(ctx, tx, boil.Infer())
	}); err != nil {
		handleError(w, r, err)
		return
	}

	http.Redirect(w, r, "/e/"+e.Key, http.StatusFound)
}

func parseChoices(s string) []string {
	row, err := csv.NewReader(strings.NewReader(s)).Read()
	if err != nil {
		log.Printf("csv.Reader.read; %v\n", err)
	}
	return row
}
