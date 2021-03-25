package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/objx"

	mrand "github.com/PeerXu/meepo/pkg/util/random"
)

const (
	SECRET_CONTEXT_SIGNATURE_HASH_ALGORITHM = "hashAlgorithm"
	SECRET_CONTEXT_SIGNATURE_TIMESTAMP      = "timestamp"
	SECRET_CONTEXT_SIGNATURE_SESSION        = "session"
	SECRET_CONTEXT_SIGNATURE_PAYLOAD_HASH   = "payloadHash"
)

var (
	SECRET_ENGINE_DEFAULT_TEMPLATE = fmt.Sprintf(`{{.%s}}${{.%s}}${{.%s}}${{.%s}}`,
		SECRET_CONTEXT_SIGNATURE_HASH_ALGORITHM,
		SECRET_CONTEXT_SIGNATURE_TIMESTAMP,
		SECRET_CONTEXT_SIGNATURE_SESSION,
		SECRET_CONTEXT_SIGNATURE_PAYLOAD_HASH,
	)
)

var (
	newHashFuncs = map[string]func() hash.Hash{
		"sha256": sha256.New,
	}
)

func getNewHashFunc(algo string) (func() hash.Hash, error) {
	fn, ok := newHashFuncs[algo]
	if !ok {
		return nil, UnsupportedHashAlgorithmError
	}

	return fn, nil
}

type SecretEngine struct {
	opt objx.Map

	secret        []byte
	hashAlgorithm string
	template      *template.Template
}

func (e *SecretEngine) hashPayload(payload Context) []byte {
	var keys []string
	var sb strings.Builder

	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		buf, _ := json.Marshal(payload[k])
		sort.Slice(buf, func(i, j int) bool { return buf[i] < buf[j] })
		sb.Write(buf)
		sb.Write([]byte{'$'})
	}

	sha := sha256.New()
	sha.Write([]byte(sb.String()))
	return sha.Sum(nil)

}

func (e *SecretEngine) Sign(payload Context) (Context, error) {
	fn, err := getNewHashFunc(e.hashAlgorithm)
	if err != nil {
		return nil, err
	}

	ts := time.Now().String()
	sess := mrand.Random.Int31()

	payload[SECRET_CONTEXT_SIGNATURE_HASH_ALGORITHM] = e.hashAlgorithm
	payload[SECRET_CONTEXT_SIGNATURE_SESSION] = sess
	payload[SECRET_CONTEXT_SIGNATURE_TIMESTAMP] = ts
	payload[SECRET_CONTEXT_SIGNATURE_PAYLOAD_HASH] = e.hashPayload(payload)

	var sb strings.Builder
	if err = e.template.Execute(&sb, payload); err != nil {
		return nil, err
	}

	mac := hmac.New(fn, e.secret)
	if _, err = mac.Write([]byte(sb.String())); err != nil {
		return nil, err
	}
	signStr := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	signature := map[string]interface{}{
		CONTEXT_NAME:                            "secret",
		SECRET_CONTEXT_SIGNATURE_HASH_ALGORITHM: e.hashAlgorithm,
		SECRET_CONTEXT_SIGNATURE_TIMESTAMP:      ts,
		SECRET_CONTEXT_SIGNATURE_SESSION:        sess,
		CONTEXT_SIGNATURE:                       signStr,
	}

	return signature, nil
}

func (e *SecretEngine) Verify(payload, signature Context) error {
	sx := objx.New(signature)

	if cast.ToString(sx.Get(CONTEXT_NAME).Inter()) != "secret" {
		return PermissionDenied
	}

	hashAlgo := cast.ToString(sx.Get(SECRET_CONTEXT_SIGNATURE_HASH_ALGORITHM).Inter())
	fn, err := getNewHashFunc(hashAlgo)
	if err != nil {
		return err
	}

	sess := cast.ToInt32(sx.Get(SECRET_CONTEXT_SIGNATURE_SESSION).Inter())
	ts := cast.ToString(sx.Get(SECRET_CONTEXT_SIGNATURE_TIMESTAMP).Inter())
	xptMac, err := base64.StdEncoding.DecodeString(cast.ToString(sx.Get(CONTEXT_SIGNATURE).Inter()))
	if err != nil {
		return err
	}

	payload[SECRET_CONTEXT_SIGNATURE_HASH_ALGORITHM] = hashAlgo
	payload[SECRET_CONTEXT_SIGNATURE_SESSION] = sess
	payload[SECRET_CONTEXT_SIGNATURE_TIMESTAMP] = ts
	payload[SECRET_CONTEXT_SIGNATURE_PAYLOAD_HASH] = e.hashPayload(payload)

	// TODO(Peer): signature reply attack: check timestamp and session
	var sb strings.Builder
	if err = e.template.Execute(&sb, payload); err != nil {
		return err
	}

	mac := hmac.New(fn, e.secret)
	if _, err = mac.Write([]byte(sb.String())); err != nil {
		return err
	}
	plMac := mac.Sum(nil)

	if !hmac.Equal(xptMac, plMac) {
		return PermissionDenied
	}

	return nil
}

func newNewSecretEngineOption() objx.Map {
	return objx.New(map[string]interface{}{})
}

func NewSecretEngine(opts ...NewEngineOption) (Engine, error) {
	o := newNewSecretEngineOption()

	for _, opt := range opts {
		opt(o)
	}

	srt := []byte(cast.ToString(o.Get("secret").Inter()))
	if len(srt) == 0 {
		return nil, fmt.Errorf("Require secret")
	}

	algo := cast.ToString(o.Get("hashAlgorithm").Inter())
	if len(algo) == 0 {
		algo = "sha256"
	}

	tmplStr := cast.ToString(o.Get("template").Inter())
	if len(tmplStr) == 0 {
		tmplStr = SECRET_ENGINE_DEFAULT_TEMPLATE
	}

	tmpl := template.Must(template.New("authSecretTemplate").Parse(tmplStr))

	return &SecretEngine{
		opt:           o,
		secret:        srt,
		hashAlgorithm: algo,
		template:      tmpl,
	}, nil
}

func init() {
	RegisterNewEngineFunc("secret", NewSecretEngine)
}
