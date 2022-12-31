package crypto_core

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteBufferWithSize(t *testing.T) {
	for _, c := range []struct {
		buf    []byte
		expect []byte
	}{
		{[]byte{0x00}, []byte{0x01, 0x00, 0x00, 0x00, 0x00}},
		{[]byte{0x01}, []byte{0x01, 0x00, 0x00, 0x00, 0x01}},
		{[]byte{0x01, 0x02, 0x03}, []byte{0x03, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}},
	} {
		wr := new(bytes.Buffer)
		assert.Nil(t, WriteBufferWithSize(wr, c.buf))
		assert.Equal(t, c.expect, wr.Bytes())
	}
}

func TestReadBufferWithSize(t *testing.T) {
	for _, c := range []struct {
		buf    []byte
		expect []byte
	}{
		{[]byte{0x01, 0x00, 0x00, 0x00, 0x01}, []byte{0x01}},
		{[]byte{0x01, 0x00, 0x00, 0x00, 0x02}, []byte{0x02}},
		{[]byte{0x03, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03}, []byte{0x01, 0x02, 0x03}},
	} {
		var ptr []byte
		rd := bytes.NewReader(c.buf)
		assert.Nil(t, ReadBufferWithSize(rd, &ptr))
		assert.Equal(t, c.expect, ptr)
	}
}
