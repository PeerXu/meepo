package cmd

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/imdario/mergo"
	"github.com/pion/webrtc/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/PeerXu/meepo/pkg/lib/config"
)

func initConfig() {
	var buf []byte
	var err error

	cfgFile = config.FindConfigPath(cfgFile)
	if cfgFile != "" {
		var dstBuf []byte

		buf, err = ioutil.ReadFile(cfgFile)
		cobra.CheckErr(err)

		src := make(map[string]any)
		err = yaml.Unmarshal(buf, src)
		cobra.CheckErr(err)

		dst := make(map[string]any)
		dstBuf, err = yaml.Marshal(config.Default())
		cobra.CheckErr(err)

		err = yaml.Unmarshal(dstBuf, &dst)
		cobra.CheckErr(err)

		err = mergo.MergeWithOverwrite(&dst, src)
		cobra.CheckErr(err)

		buf, err = yaml.Marshal(dst)
	} else {
		buf, err = yaml.Marshal(config.Default())
	}
	cobra.CheckErr(err)

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadConfig(bytes.NewReader(buf))
	cobra.CheckErr(err)

	err = viper.Unmarshal(config.Get())
	cobra.CheckErr(err)

	config.Get().Init()
}

func renderBool(x bool) string {
	if x {
		return "O"
	} else {
		return "X"
	}
}

type BindFlagsStruct struct {
	ConfigKey  string
	CommandKey string
}

func bindFlags(fs *pflag.FlagSet, cs []BindFlagsStruct) {
	for _, c := range cs {
		viper.BindPFlag(c.ConfigKey, fs.Lookup(c.CommandKey)) // nolint:errcheck
	}
}

func encodeSessionDescription(sd webrtc.SessionDescription) (string, error) {
	buf, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		return "", err
	}

	buf, err = zip(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}

func decodeSessionDescription(s string, sd *webrtc.SessionDescription) error {
	buf, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	buf, err = unzip(buf)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, sd)
}

func readFromStdin() (string, error) {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				return "", err
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}
	return in, nil
}

func zip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		return nil, err
	}
	err = gz.Flush()
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func unzip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
