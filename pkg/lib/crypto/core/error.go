package crypto_core

import "fmt"

var (
	ErrInvalidBuffer    = fmt.Errorf("invalid buffer")
	ErrInvalidSignature = fmt.Errorf("invalid signature")
)
