package http_api

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func ParseError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}

func ExtractError(buf []byte) error {
	var e struct {
		Error string `json:"error"`
	}

	json.NewDecoder(bytes.NewReader(buf)).Decode(&e)

	return fmt.Errorf(e.Error)
}
