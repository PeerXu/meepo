package routing_table

import (
	"bytes"
	"crypto/rand"
	"io"
	"sort"
	"sync"

	"github.com/PeerXu/meepo/pkg/lib/cpl"
	"github.com/PeerXu/meepo/pkg/lib/option"
)

const (
	BUCKET_COUNT    = 256
	MAX_BUCKET_SIZE = 8
)

type bucket []ID

func (b bucket) Contain(x ID) bool {
	_, err := b.IndexOf(x)
	return err == nil
}

func (b bucket) IndexOf(x ID) (int, error) {
	for i, y := range b {
		if bytes.Equal(x.Bytes(), y.Bytes()) {
			return i, nil
		}
	}
	return 0, ErrIDNotInBucket
}

func (b bucket) Len() int {
	return len(b)
}

func (b bucket) Remove(x ID) bucket {
	i, err := b.IndexOf(x)
	if err != nil {
		return b
	}
	return append(b[0:i], b[i+1:b.Len()]...)
}

func defaultNewRoutingTableOptions() option.Option {
	return option.NewOption(map[string]any{
		OPTION_BUCKET_COUNT:      BUCKET_COUNT,
		OPTION_MAX_BUCKET_SIZE:   MAX_BUCKET_SIZE,
		OPTION_RAND_BYTES_READER: rand.Reader,
	})
}

type routingTable struct {
	localID ID
	buckets map[int]bucket

	mtx             sync.Mutex
	bucketCount     int
	maxBucketSize   int
	randBytesReader io.Reader
}

func NewRoutingTable(opts ...option.ApplyOption) (RoutingTable, error) {
	o := option.ApplyWithDefault(defaultNewRoutingTableOptions(), opts...)

	localID, err := GetID(o)
	if err != nil {
		return nil, err
	}

	bucketCount, err := GetBucketCount(o)
	if err != nil {
		return nil, err
	}

	maxBucketSize, err := GetMaxBucketSize(o)
	if err != nil {
		return nil, err
	}

	randBytesReader, err := GetRandBytesReader(o)
	if err != nil {
		return nil, err
	}

	buckets := newBuckets(bucketCount)

	return &routingTable{
		localID:         localID,
		buckets:         buckets,
		bucketCount:     bucketCount,
		maxBucketSize:   maxBucketSize,
		randBytesReader: randBytesReader,
	}, nil
}

func (rt *routingTable) LocalID() ID {
	return rt.localID
}

func (rt *routingTable) TableSize() int {
	return rt.bucketCount
}

func (rt *routingTable) BucketSize(cpl int) int {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	return rt.getBucket(cpl).Len()
}

func (rt *routingTable) GenRandIDWithCpl(cpl int) ID {
	sz := (rt.bucketCount + 1) / BYTE_WIDTH
	x := rt.LocalID().Bytes()
	y := make([]byte, sz)
	_, err := rt.randBytesReader.Read(y)
	if err != nil {
		panic(err)
	}
	for i := 0; i < sz; i++ {
		n := cpl - i*BYTE_WIDTH
		if n > BYTE_WIDTH {
			n = 8
		} else if n < 0 {
			n = 0
		}
		m := BYTE_WIDTH - n
		y[i] = x[i]&(_0XFF<<m) + y[i]&(_0XFF>>n)
	}
	return FromBytes(y)
}

func (rt *routingTable) CommonPrefixLen(x ID) int {
	return cpl.CommonPrefixLen(rt.LocalID(), x)
}

func (rt *routingTable) AddID(x ID) error {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	i := rt.bucketOf(x)
	b := rt.getBucket(i)

	if b.Contain(x) {
		return nil
	}

	if b.Len() >= rt.maxBucketSize {
		return ErrOutOfBucketSize
	}

	rt.setBucket(i, append(b, x))

	return nil
}

func (rt *routingTable) RemoveID(x ID) error {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	i := rt.bucketOf(x)
	b := rt.getBucket(i)

	rt.setBucket(i, b.Remove(x))

	return nil
}

func (rt *routingTable) NearestIDs(x ID, count int, excludes []ID) (ys []ID, found bool) {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	return rt.nearestIDs(x, count, false, excludes)
}

func (rt *routingTable) ClosestIDs(x ID, count int) (ys []ID, found bool) {
	rt.mtx.Lock()
	defer rt.mtx.Unlock()

	return rt.nearestIDs(x, count, true, nil)
}

func (rt *routingTable) nearestIDs(x ID, count int, strict bool, excludes []ID) (ys []ID, found bool) {
	defer func() {
		cb := NewComparableBucket(x, ys)
		sort.Sort(sort.Reverse(cb))
		if cb.b.Len() > count {
			ys = cb.b[:count]
		} else {
			ys = cb.b
		}

		if len(ys) > 0 {
			found = bytes.Equal(ys[0].Bytes(), x.Bytes())
		}
	}()

	i := rt.bucketOf(x)
	for t := i; t < rt.bucketCount && len(ys) < count; t++ {
		b := rt.getBucket(t)
		ys = append(ys, excludeIDs(b, excludes)...)
		if strict {
			return
		}

	}
	for t := i - 1; t >= 0 && len(ys) < count; t-- {
		b := rt.getBucket(t)
		ys = append(ys, excludeIDs(b, excludes)...)
	}

	return
}

func (rt *routingTable) bucketOf(x ID) int {
	return cpl.CommonPrefixLen(rt.LocalID(), x)
}

func (rt *routingTable) getBucket(cpl int) bucket {
	return rt.buckets[cpl]
}

func (rt *routingTable) setBucket(cpl int, b bucket) {
	rt.buckets[cpl] = b
}

func newBuckets(bucketCount int) (y map[int]bucket) {
	y = make(map[int]bucket, bucketCount)
	for i := 0; i < bucketCount; i++ {
		y[i] = make([]ID, 0)
	}
	return
}

func excludeIDs(xs []ID, excludes []ID) (ys []ID) {
	for _, x := range xs {
		inExcludes := false
		for _, exclude := range excludes {
			if x.Equal(exclude) {
				inExcludes = true
				break
			}
		}
		if !inExcludes {
			ys = append(ys, x)
		}
	}
	return
}
