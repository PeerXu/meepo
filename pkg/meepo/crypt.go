package meepo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/PeerXu/meepo/pkg/signaling"
	"github.com/PeerXu/meepo/pkg/util/base36"
	mcrypt "github.com/PeerXu/meepo/pkg/util/crypt"
	"github.com/pion/webrtc/v3"
)

func b64DecodeStringFromMap(k string, m map[string]interface{}) (r []byte, err error) {
	v, ok := m[k]
	if !ok {
		err = NotFoundError
		return
	}

	s, ok := v.(string)
	if !ok {
		err = UnexpectedTypeError
		return
	}

	r, err = base64.StdEncoding.DecodeString(s)

	return
}

func generateRandomEd25519KeyPair() (pubk ed25519.PublicKey, prik ed25519.PrivateKey) {
	pubk, prik, _ = ed25519.GenerateKey(rand.Reader)
	return
}

func generateGcmNonce() []byte {
	nonce := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce)
	return nonce
}

const (
	MEEPO_ID_MAGIC_CODE = byte(0x22)
)

func Ed25519PublicKeyToMeepoID(pubk ed25519.PublicKey) string {
	return base36.Encode(append([]byte{MEEPO_ID_MAGIC_CODE}, pubk...))
}

func MeepoIDToEd25519PublicKey(peerID string) (pubk ed25519.PublicKey, err error) {
	buf := base36.Decode(peerID)
	if len(buf) == 0 {
		return nil, InvalidPeerIDError
	}

	if buf[0] != MEEPO_ID_MAGIC_CODE {
		return nil, InvalidPeerIDError
	}

	return ed25519.PublicKey(buf[1:]), nil
}

func (mp *Meepo) newGCM(secret []byte) (cipher.AEAD, error) {
	blk, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(blk)
}

func (mp *Meepo) signDescriptor(d *signaling.Descriptor) error {
	buf, err := json.Marshal(d)
	if err != nil {
		return err
	}

	sig := ed25519.Sign(mp.prik, buf)
	d.UserData["sig"] = base64.StdEncoding.EncodeToString(sig)

	return nil
}

func (mp *Meepo) verifyDescriptor(peerID string, d *signaling.Descriptor) (err error) {
	peerPubk, err := MeepoIDToEd25519PublicKey(peerID)
	if err != nil {
		return
	}

	sig, err := b64DecodeStringFromMap("sig", d.UserData)
	if err != nil {
		return
	}

	delete(d.UserData, "sig")

	buf, err := json.Marshal(d)
	if err != nil {
		return
	}

	if !ed25519.Verify(peerPubk, buf, sig) {
		return IncorrectSignatureError
	}

	return nil
}

func (mp *Meepo) marshalRequestDescriptor(peerID string, offer *webrtc.SessionDescription) (req *signaling.Descriptor, gcm cipher.AEAD, nonce []byte, err error) {
	randPubk, randPrik := generateRandomEd25519KeyPair()

	peerPubk, err := MeepoIDToEd25519PublicKey(peerID)
	if err != nil {
		return
	}

	nonce = generateGcmNonce()
	secret := mcrypt.CalcSharedSecret(peerPubk, randPrik)
	gcm, err = mp.newGCM(secret[:])
	if err != nil {
		return
	}

	plaintext, err := json.Marshal(offer)
	if err != nil {
		return
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	req = &signaling.Descriptor{
		ID: mp.GetID(),
		UserData: map[string]interface{}{
			"randPubk": base64.StdEncoding.EncodeToString(randPubk),
			"ct":       base64.StdEncoding.EncodeToString(ciphertext),
			"nonce":    base64.StdEncoding.EncodeToString(nonce),
		},
	}

	if err = mp.signDescriptor(req); err != nil {
		return
	}

	return
}

func (mp *Meepo) unmarshalResponseDescriptor(res *signaling.Descriptor, gcm cipher.AEAD, nonce []byte) (answer *webrtc.SessionDescription, err error) {
	err = mp.verifyDescriptor(res.ID, res)
	if err != nil {
		return
	}

	ciphertext, err := b64DecodeStringFromMap("ct", res.UserData)
	if err != nil {
		return
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(plaintext, &answer); err != nil {
		return
	}

	return
}

func (mp *Meepo) unmarshalRequestDescriptor(req *signaling.Descriptor) (offer *webrtc.SessionDescription, gcm cipher.AEAD, nonce []byte, err error) {
	peerID := req.ID
	if err = mp.verifyDescriptor(peerID, req); err != nil {
		return
	}

	bRandPubk, err := b64DecodeStringFromMap("randPubk", req.UserData)
	if err != nil {
		return
	}
	randPubk := ed25519.PublicKey(bRandPubk)

	if nonce, err = b64DecodeStringFromMap("nonce", req.UserData); err != nil {
		return
	}

	ciphertext, err := b64DecodeStringFromMap("ct", req.UserData)
	if err != nil {
		return
	}

	secret := mcrypt.CalcSharedSecret(randPubk, mp.prik)

	if gcm, err = mp.newGCM(secret[:]); err != nil {
		return
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(plaintext, &offer); err != nil {
		return
	}

	return
}

func (mp *Meepo) marshalResponseDescriptor(answer *webrtc.SessionDescription, gcm cipher.AEAD, nonce []byte) (res *signaling.Descriptor, err error) {
	plaintext, err := json.Marshal(answer)
	if err != nil {
		return
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	res = &signaling.Descriptor{
		ID: mp.GetID(),
		UserData: map[string]interface{}{
			"ct": base64.StdEncoding.EncodeToString(ciphertext),
		},
	}

	if err = mp.signDescriptor(res); err != nil {
		return
	}

	return
}
