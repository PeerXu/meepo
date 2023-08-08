package meepo_eventloop_core

import (
	"math/rand"
	"strings"
)

const (
	LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func randomID(sz int) string {
	var sb strings.Builder
	for i := 0; i < sz; i++ {
		sb.WriteByte(LETTERS[rand.Intn(62)])
	}
	return sb.String()
}
