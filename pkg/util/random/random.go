package random

import (
	"math/rand"
	"time"
)

var (
	Random *rand.Rand
)

func init() {
	Random = rand.New(rand.NewSource(time.Now().UnixNano()))
}
