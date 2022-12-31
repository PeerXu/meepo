package routing_table

import "github.com/PeerXu/meepo/pkg/internal/cpl"

type comparableBucket struct {
	x ID
	b bucket
}

func NewComparableBucket(x ID, b bucket) *comparableBucket {
	var copy bucket = make(bucket, 0)
	copy = append(copy, b...)
	return &comparableBucket{x, copy}
}

func (t *comparableBucket) Len() int {
	return t.b.Len()
}

func (t *comparableBucket) Less(i, j int) bool {
	di := cpl.CommonPrefixLen(t.x, t.b[i])
	dj := cpl.CommonPrefixLen(t.x, t.b[j])
	return di < dj
}

func (t *comparableBucket) Swap(i, j int) {
	t.b[i], t.b[j] = t.b[j], t.b[i]
}
