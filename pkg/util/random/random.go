package random

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
)

var (
	Random *rand.Rand
)

func SillyName() string {
	return fmt.Sprintf("%s%s%d", randomdata.Noun(), randomdata.Adjective(), Random.Int31n(1000000))
}

func init() {
	Random = rand.New(rand.NewSource(time.Now().UnixNano()))
}
