package random

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func NewRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = rune(charset[seededRand.Intn(len(charset))])
	}
	return string(b)
}
