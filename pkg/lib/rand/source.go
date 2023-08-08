package rand

import "math/rand"

var globalSource rand.Source

func GetSource() rand.Source { return globalSource }
