package rvapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"sync"
	"time"
)

// Digester encapsulates concerns around message digests.
type Digester struct {
	pool *sync.Pool
}

// NewDigester returns a new digester for the provided key. It's thread safe.
func NewDigester(key string) *Digester {
	return &Digester{
		pool: &sync.Pool{
			New: func() interface{} {
				return hmac.New(sha256.New, []byte(key))
			},
		},
	}
}

func (d *Digester) calculate(r *TrustedCheckInRequest, b []byte) []byte {
	h := d.pool.Get().(hash.Hash)
	defer func() {
		h.Reset()
		d.pool.Put(h)
	}()
	return calculate(h, r, b)
}

func calculate(h hash.Hash, r *TrustedCheckInRequest, b []byte) []byte {
	// can't fail
	_, _ = fmt.Fprintf(h, "v0:%v:%v", r.Timestamp.AsTime().Unix(), r.UserName)
	return h.Sum(b)
}

// SetDigest calculates the digest for this request and sets it; the timestamp and username must be set.
func (d *Digester) SetDigest(r *TrustedCheckInRequest) error {
	if r.Timestamp == nil {
		return errors.New("timestamp cannot be nil")
	}
	if r.UserName == "" {
		return errors.New("username cannot be empty")
	}
	if cap(r.Digest) < sha256.Size {
		r.Digest = make([]byte, 0, sha256.Size)
	}
	r.Digest = d.calculate(r, r.Digest[:0])
	return nil
}

var bufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, sha256.Size)
		return &b
	},
}

// Valid returns whether or not this request is valid for the provided key
func Valid(r *TrustedCheckInRequest, now time.Time, key []byte) bool {
	{
		t := now.Sub(r.Timestamp.AsTime())
		if t < 0 {
			t *= -1
		}
		if t >= time.Minute {
			return false
		}
	}

	b := bufPool.Get().(*[]byte)
	defer func() { bufPool.Put(b) }()

	return hmac.Equal(r.Digest, calculate(hmac.New(sha256.New, key), r, (*b)[:0]))
}
