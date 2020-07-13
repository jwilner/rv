package main

import "context"

type ctxKey int

const (
	_ ctxKey = iota
	debugKey
)

func setDebug(ctx context.Context, debug bool) context.Context {
	return context.WithValue(ctx, debugKey, debug)
}

func isDebug(ctx context.Context) bool {
	b, _ := ctx.Value(debugKey).(bool)
	return b
}
