package msgpack

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type Marshaler = msgpack.Marshaler
type Unmarshaler = msgpack.Unmarshaler

func Marshal(v interface{}) ([]byte, error) {
	enc := msgpack.GetEncoder()

	var buf bytes.Buffer
	enc.Reset(&buf)
	enc.SetSortMapKeys(true)
	enc.UseCompactInts(true)
	enc.UseCompactFloats(true)

	err := enc.Encode(v)
	b := buf.Bytes()

	msgpack.PutEncoder(enc)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func Unmarshal(b []byte, v interface{}) error {
	return msgpack.Unmarshal(b, v)
}
