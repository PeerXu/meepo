package meepo_debug_sdk_http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ParseError(res *http.Response) error {
	var errInst struct{ Error string }

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &errInst); err != nil {
		return fmt.Errorf("unexpected error: %s", body)
	}

	return fmt.Errorf(errInst.Error)
}
