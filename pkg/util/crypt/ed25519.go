package crypt

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
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

	v, err := ssh.ParseRawPrivateKey(buf)
	if err != nil {
		return
	}

	switch v.(type) {
	case *ed25519.PrivateKey:
		prik = *v.(*ed25519.PrivateKey)
		pubk = prik.Public().(ed25519.PublicKey)
		return
	default:
		err = fmt.Errorf("expect openssh ed25519 private key")
		return
	}
}
