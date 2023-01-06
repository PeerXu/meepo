package cmd

import (
	"bytes"
	"io/ioutil"

	"github.com/imdario/mergo"
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
