package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/PeerXu/meepo/pkg/lib/config"
	"github.com/PeerXu/meepo/pkg/lib/constant"
)

var cfgFile string

var (
	rootCmd = &cobra.Command{
		Use:          "meepo",
		SilenceUsage: true,
	}
)

func Execute() {
	rootCmd.Execute() // nolint:errcheck
}

// nolint:errcheck
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Configuration file")

	rootCmd.PersistentFlags().StringVar(&config.Get().Meepo.Log.Level, "log-level", "error", "Log level")
	viper.BindPFlag("meepo.log.level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().StringVarP(&config.Get().Meepo.API.Host, "host", "H", constant.API_HOST, "Meepo API server address")
	viper.BindPFlag("meepo.api.host", rootCmd.PersistentFlags().Lookup("host"))
}
