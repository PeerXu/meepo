package auth

import "fmt"

var (
	PermissionDenied              = fmt.Errorf("Permission denied")
	UnsupportedAuthEngineError    = fmt.Errorf("Unsupported auth engine")
	UnsupportedHashAlgorithmError = fmt.Errorf("Unsupported hash algorithm")
)
