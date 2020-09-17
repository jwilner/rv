package platform

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jwilner/rv/pkg/pb/rvapi"
)

type handler struct {
	rvapi.UnimplementedRVerServer

	txM  *txMgr
	kGen *stringGener
	tzes []string
}

var _ rvapi.RVerServer = new(handler)

func (h *handler) Close() error {
	return h.txM.db.Close()
}

const keyCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

type ctxKey int

const (
	_ ctxKey = iota
	debugKey
	requestIDKey
)

func setDebug(ctx context.Context, debug bool) context.Context {
	return context.WithValue(ctx, debugKey, debug)
}

func isDebug(ctx context.Context) bool {
	b, _ := ctx.Value(debugKey).(bool)
	return b
}

func requestID(ctx context.Context) string {
	s, _ := ctx.Value(requestIDKey).(string)
	if s == "" {
		return "x-no-request-id-set"
	}
	return s
}
