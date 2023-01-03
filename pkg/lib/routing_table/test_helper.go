package routing_table

func bytesSlice2Bucket(xs [][]byte) bucket {
	ys := make([]ID, 0)
	for _, x := range xs {
		ys = append(ys, FromBytes(x))
	}
	return ys
}

func bytesSliceMap2Buckets(x map[int][][]byte, bucketCount int) map[int]bucket {
	buckets := newBuckets(bucketCount)
	for i, bytesSlice := range x {
		bucket := bytesSlice2Bucket(bytesSlice)
		buckets[i] = bucket
	}
	return buckets
}
