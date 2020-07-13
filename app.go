package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	paramBallotKey   = "ballotKey"
	paramElectionKey = "electionKey"
)

func ballotKey(ctx context.Context) string {
	return chi.URLParamFromCtx(ctx, paramBallotKey)
}

func electionKey(ctx context.Context) string {
	return chi.URLParamFromCtx(ctx, paramElectionKey)
}

type handler struct {
	tmpls *tmplMgr
	txM   *txMgr
}

// form contains common form behavior for template handling
type form struct {
	// Errors is a map of field to error messages
	Errors map[string][]string
}

func (f *form) setErrorf(field, msg string, args ...interface{}) {
	if f.Errors == nil {
		f.Errors = make(map[string][]string)
	}
	f.Errors[field] = append(f.Errors[field], fmt.Sprintf(msg, args...))
}

func (f *form) checkErrors() bool {
	return len(f.Errors) == 0
}

func route(h *handler) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.getIndex)
		r.Post("/", h.postIndex)
	})
	r.Route("/b/{"+paramBallotKey+"}", func(r chi.Router) {
		r.Get("/", h.getBallot)
		r.Post("/", h.postBallot)
	})
	r.Route("/e/{"+paramElectionKey+"}", func(r chi.Router) {
		r.Get("/", h.getElection)
		r.Post("/", h.postElection)
	})

	r.NotFound(serve404)
	r.MethodNotAllowed(serve405)

	return r
}

func serve404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func serve405(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, sql.ErrNoRows):
		serve404(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	return true
}

type txMgr struct {
	db *sql.DB
}

func (t *txMgr) inTx(
	ctx context.Context,
	opts *sql.TxOptions,
	op func(ctx context.Context, tx *sql.Tx) error,
) (err error) {
	var tx *sql.Tx
	if tx, err = t.db.BeginTx(ctx, opts); err != nil {
		return fmt.Errorf("db.BeginTx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil { // is there a panic?
			log.Printf("rolling back tx after panic: %v\n", p)
			if err := tx.Rollback(); err != nil {
				log.Printf("failed rolling back after panic: %v\n", p)
			}
			panic(p) // continue panic
		}

		if err != nil { // did an error occur during the transaction?
			log.Printf("rolling back tx after error: %v", err)
			if rErr := tx.Rollback(); rErr != nil {
				log.Printf("failed rolling back after error: %v\n", rErr)
			}
			return
		}

		// no error yet -- commit transaction, report any error
		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("tx.Commit: %w", err)
		}
	}()

	// execute the transaction and let the deferred handle any error
	err = op(ctx, tx)

	return
}

