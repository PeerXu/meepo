package crypt

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
)

func Ed25519GenerateKey() (pubk ed25519.PublicKey, prik ed25519.PrivateKey) {
	pubk, prik, _ = ed25519.GenerateKey(rand.Reader)
	return
}

func LoadEd25519Key(filename string) (pubk ed25519.PublicKey, prik ed25519.PrivateKey, err error) {
	filename, err = homedir.Expand(filename)
	if err != nil {
		return
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	blk, _ := pem.Decode(buf)
	if err != nil {
		return
	}

	v, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
	if err != nil {
		return
	}
	switch v.(type) {
	case ed25519.PrivateKey:
		prik = v.(ed25519.PrivateKey)
		pubk = prik.Public().(ed25519.PublicKey)
	default:
		err = fmt.Errorf("Expect ed25519 private key")
		return
	}

	return
}
