package io

import (
	"io"
	"os"
	"strconv"
	"sync"
)

const defualtCopyBufSize = 32 * 1024

var copyBufPool *sync.Pool

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := copyBufPool.Get().([]byte)
	defer copyBufPool.Put(buf) // nolint:staticcheck
	return io.CopyBuffer(dst, src, buf)
}

func init() {
	copyBufSize := defualtCopyBufSize
	var err error
	copyBufSizeStr := os.Getenv("MPO_EXPERIMENTAL_COPY_BUF_SIZE")
	if copyBufSizeStr != "" {
		copyBufSize, err = strconv.Atoi(copyBufSizeStr)
		if err != nil {
			panic(err)
		}
	}
	copyBufPool = &sync.Pool{
		New: func() any { return make([]byte, copyBufSize) },
	}
}
