package rand

import (
	"math/rand"
	"strings"
)

const (
	FULL_LETTERS_AND_DIGITS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type StringGenerator struct {
	alphabet string
	rand     *rand.Rand
}

func NewStringGenerator(alphabet string, rand *rand.Rand) *StringGenerator {
	return &StringGenerator{
		alphabet: alphabet,
		rand:     rand,
	}
}

func (sg *StringGenerator) Generate(sz int) string {
	var sb strings.Builder
	alphabetSize := len(sg.alphabet)
	for i := 0; i < sz; i++ {
		sb.WriteByte(sg.alphabet[sg.rand.Intn(alphabetSize)])
	}
	return sb.String()
}

var DefaultStringGenerator *StringGenerator

func init() {
	DefaultStringGenerator = NewStringGenerator(FULL_LETTERS_AND_DIGITS, rand.New(GetSource()))
}
