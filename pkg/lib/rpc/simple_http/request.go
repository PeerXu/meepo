package rpc_simple_http

type SimpleDoRequest struct {
	Method      string `json:"method"`
	CallRequest []byte `json:"callRequest,omitempty"`
}
