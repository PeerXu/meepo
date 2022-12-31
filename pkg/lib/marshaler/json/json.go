package marshaler_json

import (
	"encoding/json"

	marshaler_interface "github.com/PeerXu/meepo/pkg/lib/marshaler/interface"
)

var (
	Marshaler   marshaler_interface.MarshalFunc   = json.Marshal
	Unmarshaler marshaler_interface.UnmarshalFunc = json.Unmarshal
)
