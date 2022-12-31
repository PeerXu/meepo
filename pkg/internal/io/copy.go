package io

import (
	"io"
	"sync"
)

const copyBufSize = 32 * 1024

var copyBufPool *sync.Pool

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := copyBufPool.Get().([]byte)
	defer copyBufPool.Put(buf)
	return io.CopyBuffer(dst, src, buf)
}

func init() {
	copyBufPool = &sync.Pool{
		New: func() any { return make([]byte, copyBufSize) },
	}
}
