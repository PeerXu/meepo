package routing_table

import "fmt"

var (
	ErrIDNotInBucket   = fmt.Errorf("id not in bucket")
	ErrOutOfBucketSize = fmt.Errorf("out of bucket size")
)
