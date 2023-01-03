package routing_table

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketContain(t *testing.T) {
	for i, c := range []struct {
		bytesSlice [][]byte
		expect     []byte
	}{
		{[][]byte{{0x00}}, []byte{0x00}},
		{[][]byte{{0xff}}, []byte{0xff}},
		{[][]byte{{0xab}, {0xcd}, {0xef}}, []byte{0xcd}},
	} {
		ids := make([]ID, len(c.bytesSlice))
		for i, bytes := range c.bytesSlice {
			ids[i] = FromBytes(bytes)
		}
		assert.True(t, bucket(ids).Contain(FromBytes(c.expect)), "index=%v, bytesSlice=%v, expect=%v", i, c.bytesSlice, c.expect)
	}
}

func TestBucketIndexOf(t *testing.T) {
	for i, c := range []struct {
		bytesSlice [][]byte
		x          []byte
		expect     int
		err        error
	}{
		{[][]byte{{0xff}}, []byte{0xff}, 0, nil},
		{[][]byte{{0xff}, {0x00}, {0x11}, {0x22}}, []byte{0x11}, 2, nil},
		{[][]byte{{0xff}, {0x00}, {0x11}, {0x22}}, []byte{0x33}, 0, ErrIDNotInBucket},
	} {
		ids := make([]ID, len(c.bytesSlice))
		for i, bytes := range c.bytesSlice {
			ids[i] = FromBytes(bytes)
		}
		idx, err := bucket(ids).IndexOf(FromBytes(c.x))
		assert.Equal(t, c.err, err, "index=%v, bytesSlice=%v, x=%v", i, c.bytesSlice, c.x)
		assert.Equal(t, c.expect, idx, "index=%v, bytesSlice=%v, x=%v", i, c.bytesSlice, c.x)

	}
}

func TestBucketRemove(t *testing.T) {
	for i, c := range []struct {
		bytesSlice [][]byte
		x          []byte
		expect     [][]byte
	}{
		{[][]byte{{0xff}}, []byte{0xff}, [][]byte{}},
		{[][]byte{{0xff}, {0x11}, {0x22}}, []byte{0x11}, [][]byte{{0xff}, {0x22}}},
	} {
		ids := make([]ID, len(c.bytesSlice))
		for i, bytes := range c.bytesSlice {
			ids[i] = FromBytes(bytes)
		}
		out := bucket(ids).Remove(FromBytes(c.x))
		expectIDs := make([]ID, len(c.expect))
		for i, bytes := range c.expect {
			expectIDs[i] = FromBytes(bytes)
		}
		assert.Equal(t, bucket(expectIDs), out, "index=%v, bytesSlice=%v, x=%v", i, c.bytesSlice, c.x)
	}
}

func TestRoutingTableGenRandIDWithCpl(t *testing.T) {
	for i, c := range []struct {
		local  []byte
		rand   []byte
		cpl    int
		expect []byte
	}{
		{[]byte{0xff}, []byte{0x01}, 4, []byte{0xf1}},
		{[]byte{0xaa}, []byte{0x55}, 1, []byte{0xd5}},
		{[]byte{0xff}, []byte{0x00}, 0, []byte{0x00}},
		{[]byte{0xff}, []byte{0x00}, 8, []byte{0xff}},
		{[]byte{0xab, 0xcd}, []byte{0x12, 0x34}, 0, []byte{0x12, 0x34}},
		{[]byte{0xab, 0xcd}, []byte{0x12, 0x34}, 2, []byte{0x92, 0x34}},
		{[]byte{0xab, 0xcd}, []byte{0x12, 0x34}, 4, []byte{0xa2, 0x34}},
		{[]byte{0xab, 0xcd}, []byte{0x12, 0x34}, 8, []byte{0xab, 0x34}},
		{[]byte{0xab, 0xcd}, []byte{0x12, 0x34}, 12, []byte{0xab, 0xc4}},
		{[]byte{0xab, 0xcd}, []byte{0x12, 0x34}, 15, []byte{0xab, 0xcc}},
		{[]byte{0xab, 0xcd}, []byte{0x12, 0x34}, 16, []byte{0xab, 0xcd}},
	} {
		rt := &routingTable{
			localID:         FromBytes(c.local),
			randBytesReader: bytes.NewReader(c.rand),
			bucketCount:     BYTE_WIDTH * len(c.local),
		}
		out := rt.GenRandIDWithCpl(c.cpl)
		assert.Equal(t, c.expect, out.Bytes(), "index=%v, local=%v, rand=%v, cpl=%v", i, c.local, c.rand, c.cpl)
	}
}

func TestRoutingTableAddID(t *testing.T) {
	for i, c := range []struct {
		localID       []byte
		buckets       map[int][][]byte
		bucketCount   int
		maxBucketSize int
		xs            [][]byte
		expectBuckets map[int][][]byte
		expectError   error
	}{
		{[]byte{0x00}, map[int][][]byte{}, 8, 4, [][]byte{{0xff}}, map[int][][]byte{0: {{0xff}}}, nil},
		{[]byte{0x00}, map[int][][]byte{}, 8, 0, [][]byte{{0xff}}, map[int][][]byte{}, ErrOutOfBucketSize},
		{[]byte{0x00}, map[int][][]byte{}, 8, 1, [][]byte{{0xff}, {0xf0}}, map[int][][]byte{0: {{0xff}}}, ErrOutOfBucketSize},
	} {
		var err error

		rt := &routingTable{
			localID:       FromBytes(c.localID),
			buckets:       newBuckets(c.bucketCount),
			bucketCount:   c.bucketCount,
			maxBucketSize: c.maxBucketSize,
		}
		rt.buckets = bytesSliceMap2Buckets(c.buckets, c.bucketCount)
		expectBuckets := bytesSliceMap2Buckets(c.expectBuckets, c.bucketCount)

		for _, x := range c.xs {
			if err = rt.AddID(FromBytes(x)); err != nil {
				break
			}
		}

		assert.Equal(t, c.expectError, err, "index=%v, data=%v", i, c)
		assert.Equal(t, expectBuckets, rt.buckets, "index=%v, data=%v", i, c)
	}
}

func TestRoutingTableRemoveID(t *testing.T) {
	for i, c := range []struct {
		localID       []byte
		buckets       map[int][][]byte
		bucketCount   int
		maxBucketSize int
		xs            [][]byte
		expectBuckets map[int][][]byte
	}{
		{[]byte{0x00}, map[int][][]byte{0: {{0xff}}}, 8, 1, [][]byte{{0xff}}, map[int][][]byte{}},
		{[]byte{0x00}, map[int][][]byte{0: {{0xff}}}, 8, 1, [][]byte{{0x0f}}, map[int][][]byte{0: {{0xff}}}},
	} {
		rt := &routingTable{
			localID:       FromBytes(c.localID),
			buckets:       newBuckets(c.bucketCount),
			bucketCount:   c.bucketCount,
			maxBucketSize: c.maxBucketSize,
		}
		rt.buckets = bytesSliceMap2Buckets(c.buckets, c.bucketCount)
		expectBuckets := bytesSliceMap2Buckets(c.expectBuckets, c.bucketCount)

		for _, x := range c.xs {
			assert.Nil(t, rt.RemoveID(FromBytes(x)), "index=%v, data=%v", i, c)
		}

		assert.Equal(t, expectBuckets, rt.buckets, "index=%v, data=%v", i, c)
	}
}

func TestRoutingTableNearestIDs(t *testing.T) {
	for i, c := range []struct {
		localID     []byte
		buckets     map[int][][]byte
		bucketCount int
		x           []byte
		count       int
		excludeIDs  [][]byte
		expectIDs   [][]byte
		expectFound bool
	}{
		{[]byte{0x00}, map[int][][]byte{7: {{0x01}}}, 8, []byte{0x01}, 1, nil, [][]byte{{0x01}}, true},
		{[]byte{0x00}, map[int][][]byte{7: {{0x01}}}, 8, []byte{0x01}, 2, nil, [][]byte{{0x01}}, true},
		{[]byte{0x00}, map[int][][]byte{7: {{0x01}}}, 8, []byte{0x01}, 2, nil, [][]byte{{0x01}}, true},
		{[]byte{0x00}, map[int][][]byte{7: {{0x01}}}, 8, []byte{0x02}, 1, nil, [][]byte{{0x01}}, false},
		{[]byte{0x00}, map[int][][]byte{7: {{0x01}}}, 8, []byte{0x01}, 0, nil, [][]byte{}, false},
		{[]byte{0x00}, map[int][][]byte{7: {{0x01}}}, 8, []byte{0x01}, 1, nil, [][]byte{{0x01}}, true},
		{[]byte{0x00}, map[int][][]byte{
			0: {{0x80}, {0x90}, {0xa0}},
			1: {{0x40}, {0x50}, {0x60}},
		}, 8, []byte{0xb0}, 1, nil, [][]byte{{0xa0}}, false},
		{[]byte{0x00}, map[int][][]byte{
			0: {{0x80}, {0xa0}},
			1: {{0x40}},
		}, 8, []byte{0xb0}, 3, nil, [][]byte{{0xa0}, {0x80}, {0x40}}, false},
		{[]byte{0x00}, map[int][][]byte{
			0: {{0x80}, {0xa0}},
			1: {{0x40}},
			2: {{0x20}},
		}, 8, []byte{0x60}, 2, nil, [][]byte{{0x40}, {0x20}}, false},
		{[]byte{0x00}, map[int][][]byte{
			0: {{0x80}},
			1: {{0x40}},
			2: {{0x20}},
		}, 8, []byte{0x40}, 3, nil, [][]byte{{0x40}, {0x20}, {0x80}}, true},
		{[]byte{0x00}, map[int][][]byte{
			7: {{0x01}},
			6: {{0x03}},
		}, 8, []byte{0x01}, 1, [][]byte{{0x01}}, [][]byte{{0x03}}, false},
		{[]byte{0x00}, map[int][][]byte{
			7: {{0x01}},
			6: {{0x03}},
		}, 8, []byte{0x01}, 3, [][]byte{{0x03}}, [][]byte{{0x01}}, true},
		{[]byte{0x00}, map[int][][]byte{
			7: {{0x01}},
			6: {{0x03}, {0x02}},
		}, 8, []byte{0x01}, 3, [][]byte{{0x03}}, [][]byte{{0x01}, {0x02}}, true},
		{[]byte{0x00}, map[int][][]byte{
			7: {{0x01}},
			6: {{0x03}, {0x02}},
			5: {{0x04}},
		}, 8, []byte{0x01}, 3, [][]byte{{0x03}}, [][]byte{{0x01}, {0x02}, {0x04}}, true},
		{[]byte{0x00}, map[int][][]byte{
			7: {{0x01}},
			6: {{0x03}, {0x02}},
			5: {{0x04}},
		}, 8, []byte{0x03}, 3, [][]byte{{0x03}}, [][]byte{{0x02}, {0x01}, {0x04}}, false},
		{[]byte{0x00}, map[int][][]byte{
			7: {{0x01}},
			6: {{0x03}, {0x02}},
			5: {{0x04}},
		}, 8, []byte{0x03}, 3, [][]byte{{0x02}}, [][]byte{{0x03}, {0x01}, {0x04}}, true},
	} {
		rt := &routingTable{
			localID:     FromBytes(c.localID),
			buckets:     newBuckets(c.bucketCount),
			bucketCount: c.bucketCount,
		}
		rt.buckets = bytesSliceMap2Buckets(c.buckets, c.bucketCount)
		var excludeIDs []ID
		for _, x := range c.excludeIDs {
			excludeIDs = append(excludeIDs, FromBytes(x))
		}

		actualIDs, actualFound := rt.NearestIDs(FromBytes(c.x), c.count, excludeIDs)
		assert.Equal(t, []ID(bytesSlice2Bucket(c.expectIDs)), actualIDs, "index=%v, data=%v", i, c)
		assert.Equal(t, c.expectFound, actualFound, "index=%v, data=%v", i, c)
	}
}

func TestRoutingTableClosestIDs(t *testing.T) {
	for i, c := range []struct {
		localID     []byte
		buckets     map[int][][]byte
		bucketCount int
		x           []byte
		count       int
		expectIDs   [][]byte
		expectFound bool
	}{
		{[]byte{0x00}, map[int][][]byte{7: {{0x01}}}, 8, []byte{0x01}, 1, [][]byte{{0x01}}, true},
		{[]byte{0x00}, map[int][][]byte{7: {{0x01}}}, 8, []byte{0x02}, 1, [][]byte{}, false},
		{[]byte{0x00}, map[int][][]byte{
			7: {{0x01}},
			6: {{0x02}},
		}, 8, []byte{0x02}, 1, [][]byte{{0x02}}, true},
		{[]byte{0x00}, map[int][][]byte{
			7: {{0x01}},
			6: {{0x02}, {0x03}},
		}, 8, []byte{0x02}, 2, [][]byte{{0x02}, {0x03}}, true},
		{[]byte{0x00}, map[int][][]byte{
			6: {{0x02}},
			5: {{0x04}},
			4: {{0x08}},
		}, 8, []byte{0x01}, 3, [][]byte{}, false},
	} {
		rt := &routingTable{
			localID:     FromBytes(c.localID),
			buckets:     newBuckets(c.bucketCount),
			bucketCount: c.bucketCount,
		}
		rt.buckets = bytesSliceMap2Buckets(c.buckets, c.bucketCount)

		actualIDs, actualFound := rt.ClosestIDs(FromBytes(c.x), c.count)
		assert.Equal(t, []ID(bytesSlice2Bucket(c.expectIDs)), actualIDs, "index=%v, data=%v", i, c)
		assert.Equal(t, c.expectFound, actualFound, "index=%v, data=%v", i, c)
	}
}
