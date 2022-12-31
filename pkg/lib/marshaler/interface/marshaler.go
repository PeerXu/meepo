package marshaler_interface

type Marshaler interface {
	Marshal(any) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte, any) error
}

type MarshalFunc func(any) ([]byte, error)

func (fn MarshalFunc) Marshal(x any) ([]byte, error) {
	return fn(x)
}

type UnmarshalFunc func([]byte, any) error

func (fn UnmarshalFunc) Unmarshal(x []byte, y any) error {
	return fn(x, y)
}
