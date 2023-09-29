package rand

import (
	"math/rand"
	"os"
	"strconv"
	"sync"
)

type lockedSource struct {
	lk sync.Mutex
	s  rand.Source
}

func (r *lockedSource) Int63() int64 {
	r.lk.Lock()
	defer r.lk.Unlock()
	return r.s.Int63()
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	defer r.lk.Unlock()
	r.s.Seed(seed)
}

func init() {
	var seed int64

	randomSeedStr := os.Getenv("MPO_EXPERIMENTAL_RANDOM_SEED")
	if randomSeedStr != "" {
		var err error
		seed, err = strconv.ParseInt(randomSeedStr, 10, 64)
		if err != nil {
			panic(err)
		}
	} else {
		seed = rand.Int63()
	}

	globalSource = &lockedSource{
		s: rand.NewSource(seed),
	}
}
