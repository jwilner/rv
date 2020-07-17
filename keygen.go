package main

import (
	"math/rand"
	"time"
)

// TODO don't use built-in rand's shitty concurrency
func newStringGener() *stringGener {
	rand.Seed(time.Now().UnixNano())
	return new(stringGener)
}

type stringGener struct{}

func (k *stringGener) newKey(charSet string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(b)
}
