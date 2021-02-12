package http_api

func ParseError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}
