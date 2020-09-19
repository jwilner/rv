package platform

import (
	"math/rand"
)

// TODO don't use built-in rand's shitty concurrency
func newStringGener() *stringGener {
	return new(stringGener)
}

type stringGener struct{}

func (k *stringGener) newKey(charSet string, length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b) // can't actually fail

	l := uint8(len(charSet))
	for i := range b {
		b[i] = charSet[b[i]%l]
	}
	return string(b)
}
