package rtoken

import (
	"crypto/rand"
	"encoding/base64"
	mrand "math/rand"
	"time"
)

// CategoryteRandomBytes returns securely generated random bytes.
func CategoryteRandomBytes(n int) ([]byte, error) {
	mrand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// CategoryteRandomString returns a URL-safe, base64 encoded securely generated random string.
func CategoryteRandomString(s int) (string, error) {
	b, err := CategoryteRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// CategoryteRandomID generates random id for a session
func CategoryteRandomID(s int) string {

	mrand.Seed(time.Now().UnixNano())

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, s)
	for i := range b {
		b[i] = letterBytes[mrand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
