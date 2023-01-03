package routing_table

import (
	"io"

	"github.com/PeerXu/meepo/pkg/lib/option"
)

const (
	OPTION_ID                = "id"
	OPTION_BUCKET_COUNT      = "bucketCount"
	OPTION_MAX_BUCKET_SIZE   = "maxBucketSize"
	OPTION_RAND_BYTES_READER = "randBytesReader"
)

var (
	WithID = func(v ID) option.ApplyOption {
		return func(o option.Option) {
			o[OPTION_ID] = v
		}
	}

	GetID = func(o option.Option) (ID, error) {
		var x ID
		i := o.Get(OPTION_ID).Inter()
		if i == nil {
			return x, option.ErrOptionRequiredFn(OPTION_ID)
		}

		v, ok := i.(ID)
		if !ok {
			return x, option.ErrUnexpectedTypeFn(x, i)
		}

		return v, nil
	}

	WithBucketCount, GetBucketCount     = option.New[int](OPTION_BUCKET_COUNT)
	WithMaxBucketSize, GetMaxBucketSize = option.New[int](OPTION_MAX_BUCKET_SIZE)

	WithRandBytesReader = func(rd io.Reader) option.ApplyOption {
		return func(o option.Option) {
			o[OPTION_RAND_BYTES_READER] = rd
		}
	}

	GetRandBytesReader = func(o option.Option) (io.Reader, error) {
		var x io.Reader
		i := o.Get(OPTION_RAND_BYTES_READER).Inter()
		if i == nil {
			return x, option.ErrOptionRequiredFn(OPTION_RAND_BYTES_READER)
		}

		v, ok := i.(io.Reader)
		if !ok {
			return x, option.ErrUnexpectedTypeFn(x, i)
		}

		return v, nil
	}
)
