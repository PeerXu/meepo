package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/Pallinder/go-randomdata"
	"github.com/VividCortex/godaemon"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
	"github.com/PeerXu/meepo/pkg/api"
	http_api "github.com/PeerXu/meepo/pkg/api/http"
	"github.com/PeerXu/meepo/pkg/meepo"
	"github.com/PeerXu/meepo/pkg/signaling"
	redis_signaling "github.com/PeerXu/meepo/pkg/signaling/redis"
)

var (
	serveCmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"summon"},
		Short:   "Summon a Meepo",
		RunE:    meepoSummon,
	}
)

func meepoSummon(cmd *cobra.Command, args []string) error {
	fs := cmd.Flags()
	configStr, _ := fs.GetString("config")

	cfg, loaded, err := config.Load(configStr)
	if err != nil {
		return err
	}

	if fs.Lookup("daemon").Changed {
		cfg.Meepo.Daemon, _ = fs.GetBool("daemon")
	}

	if fs.Lookup("log-level").Changed {
		cfg.Meepo.Log.Level, _ = fs.GetString("log-level")
	}

	logger := logrus.New()
	logLevel, err := logrus.ParseLevel(cfg.Meepo.Log.Level)
	if err != nil {
		return err
	}
	logger.SetLevel(logLevel)

	if logLevel < logrus.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}

	if cfg.Meepo.Daemon {
		godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	}

	if !loaded {
		logger.Warningf("Config file not found, load default config")
	}

	id := cfg.Meepo.ID
	if id == "" {
		id = randomdata.SillyName()
	}

	engineOptions := []signaling.NewEngineOption{
		signaling.WithID(id),
		signaling.WithLogger(logger),
	}

	switch cfg.Meepo.Signaling.Name {
	case "redis":
		rsCfg := cfg.Meepo.SignalingI.(*config.RedisSignalingConfig)
		engineOptions = append(
			engineOptions,
			redis_signaling.WithURL(rsCfg.URL),
		)
	}

	signalingEngine, err := signaling.NewEngine(cfg.Meepo.Signaling.Name, engineOptions...)
	if err != nil {
		return err
	}

	newMeepoOptions := []meepo.NewMeepoOption{
		meepo.WithSignalingEngine(signalingEngine),
		meepo.WithLogger(logger),
		meepo.WithID(id),
		meepo.WithICEServers(cfg.Meepo.TransportI.(*config.WebrtcTransportConfig).ICEServers),
	}
	meepo, err := meepo.NewMeepo(newMeepoOptions...)
	if err != nil {
		return err
	}

	apiCfg := cfg.Meepo.ApiI.(*config.HttpApiConfig)
	api, err := api.NewServer(
		"http",
		http_api.WithHost(apiCfg.Host),
		http_api.WithPort(apiCfg.Port),
		api.WithMeepo(meepo),
	)
	if err != nil {
		return err
	}

	if err = api.Start(context.TODO()); err != nil {
		return err
	}
	logger.Infof("api server startd")

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c
		logger.Debugf("api server terminating")

		if err = api.Stop(context.TODO()); err != nil {
			logger.WithError(err).Errorf("failed to stop api server")
		}
	}()

	api.Wait()
	logger.Infof("api server terminated")

	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringP("config", "c", "~/.meepo/config.yaml", "Location of meepo config file")
	serveCmd.PersistentFlags().BoolP("daemon", "d", true, "Run as daemon")
}
