package crypto_core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/teserakt-io/golang-ed25519/extra25519"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ssh"
)

func Curve25519FromEd25519PublicKey(x ed25519.PublicKey) [32]byte {
	var curve [32]byte
	var pubk32 [32]byte
	copy(pubk32[:], x[:32])
	extra25519.PublicKeyToCurve25519(&curve, &pubk32)
	return curve
}

func Curve25519FromEd25519PrivateKey(x ed25519.PrivateKey) [32]byte {
	var curve [32]byte
	var prik64 [64]byte
	copy(prik64[:], x[:64])
	extra25519.PrivateKeyToCurve25519(&curve, &prik64)
	return curve
}

func SecretFromEd25519(pubk ed25519.PublicKey, prik ed25519.PrivateKey) []byte {
	var secret []byte
	curvePubk := Curve25519FromEd25519PublicKey(pubk)
	curvePrik := Curve25519FromEd25519PrivateKey(prik)
	secret, _ = curve25519.X25519(curvePrik[:], curvePubk[:])
	return secret[:]
}

func NewGCM(secret []byte) (cipher.AEAD, error) {
	blk, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(blk)
}

func WriteBufferWithSize(wr io.Writer, buf []byte) (err error) {
	if err = binary.Write(wr, binary.LittleEndian, uint32(len(buf))); err != nil {
		return
	}
	if _, err = wr.Write(buf); err != nil {
		return
	}
	return
}

func ReadBufferWithSize(rd io.Reader, ptr *[]byte) (err error) {
	var sz uint32
	if err = binary.Read(rd, binary.LittleEndian, &sz); err != nil {
		return
	}

	*ptr = make([]byte, sz)
	if _, err = rd.Read(*ptr); err != nil {
		return
	}

	return
}

func LoadEd25519Key(filename string) (pubk ed25519.PublicKey, prik ed25519.PrivateKey, err error) {
	filename, err = homedir.Expand(filename)
	if err != nil {
		return
	}

	buf, err := os.ReadFile(filename)
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
