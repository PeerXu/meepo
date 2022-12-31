package crypto_interface

type Packet struct {
	Source      []byte `json:"source"`
	Destination []byte `json:"destination"`
	Nonce       []byte `json:"nonce"`
	CipherText  []byte `json:"cipherText"`
	Signature   []byte `json:"signature"`
}
