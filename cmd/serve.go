package cmd

import (
	"context"
	"crypto/ed25519"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/PeerXu/meepo/cmd/config"
	"github.com/PeerXu/meepo/pkg/api"
	http_api "github.com/PeerXu/meepo/pkg/api/http"
	"github.com/PeerXu/meepo/pkg/meepo"
	"github.com/PeerXu/meepo/pkg/signaling"
	redis_signaling "github.com/PeerXu/meepo/pkg/signaling/redis"
	mcrypt "github.com/PeerXu/meepo/pkg/util/crypt"
	mdaemon "github.com/PeerXu/meepo/pkg/util/daemon"
	meg "github.com/PeerXu/meepo/pkg/util/group"
)

var (
	serveCmd = &cobra.Command{
		Use:     "serve [-c config] [-d daemon]",
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

	switch logLevel {
	case logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel:
		gin.SetMode(gin.ReleaseMode)
	case logrus.DebugLevel:
	case logrus.TraceLevel:
	}

	if cfg.Meepo.Daemon {
		mdaemon.Daemon()
	}

	if !loaded {
		logger.Warningf("Config file not found, load default config\n")
	}

	if fs.Lookup("identity-file").Changed {
		cfg.Meepo.IdentityFile, _ = fs.GetString("identity-file")
	}

	var pubk ed25519.PublicKey
	var prik ed25519.PrivateKey
	if cfg.Meepo.IdentityFile != "" {
		pubk, prik, err = mcrypt.LoadEd25519Key(cfg.Meepo.IdentityFile)
		if err != nil {
			return err
		}
	} else {
		pubk, prik = mcrypt.Ed25519GenerateKey()
	}

	signalingEngineOptions := []signaling.NewEngineOption{
		signaling.WithID(meepo.Ed25519PublicKeyToMeepoID(pubk)),
		signaling.WithLogger(logger),
	}

	switch cfg.Meepo.Signaling.Name {
	case "redis":
		rsCfg := cfg.Meepo.SignalingI.(*config.RedisSignalingConfig)
		signalingEngineOptions = append(
			signalingEngineOptions,
			redis_signaling.WithURL(rsCfg.URL),
		)
	}

	signalingEngine, err := signaling.NewEngine(cfg.Meepo.Signaling.Name, signalingEngineOptions...)
	if err != nil {
		return err
	}

	newMeepoOptions := []meepo.NewMeepoOption{
		meepo.WithSignalingEngine(signalingEngine),
		meepo.WithLogger(logger),
		meepo.WithED25519KeyPair(pubk, prik),
		meepo.WithICEServers(cfg.Meepo.TransportI.(*config.WebrtcTransportConfig).ICEServers),
	}

	if cfg.Meepo.AsSignaling {
		newMeepoOptions = append(newMeepoOptions, meepo.WithAsSignaling(true))
	}

	switch cfg.Meepo.Auth.Name {
	case "secret":
		sa := cfg.Meepo.AuthI.(*config.SecretAuthConfig)
		newMeepoOptions = append(
			newMeepoOptions,
			meepo.WithAuthorizationName("secret"),
			meepo.WithAuthorizationSecret(sa.Secret),
		)
	case "dummy":
		fallthrough
	default:
		newMeepoOptions = append(
			newMeepoOptions,
			meepo.WithAuthorizationName("dummy"),
		)
	}

	var allows, blocks []meepo.AclPolicy
	if cfg.Meepo.Acl != nil {
		if allows, err = meepo.ParseAclPolicies(cfg.Meepo.Acl.Allows); err != nil {
			return err
		}
		if blocks, err = meepo.ParseAclPolicies(cfg.Meepo.Acl.Blocks); err != nil {
			return err
		}
	}
	newMeepoOptions = append(
		newMeepoOptions,
		meepo.WithAcl(meepo.NewAcl(meepo.WithAllowPolicies(allows), meepo.WithBlockPolicies(blocks))),
	)

	mp, err := meepo.NewMeepo(newMeepoOptions...)
	if err != nil {
		return err
	}

	apiCfg := cfg.Meepo.ApiI.(*config.HttpApiConfig)
	api, err := api.NewServer(
		"http",
		http_api.WithHost(apiCfg.Host),
		http_api.WithPort(apiCfg.Port),
		api.WithMeepo(mp),
	)
	if err != nil {
		return err
	}

	if err = api.Start(context.TODO()); err != nil {
		return err
	}
	logger.Infof("api server startd")

	var socks5 meepo.Socks5Server
	if cfg.Meepo.Proxy != nil && cfg.Meepo.Proxy.Socks5 != nil {
		socks5Cfg := cfg.Meepo.Proxy.Socks5
		socks5, err = meepo.NewSocks5Server(
			meepo.WithMeepo(mp),
			meepo.WithHost(socks5Cfg.Host),
			meepo.WithPort(socks5Cfg.Port),
		)
		if err != nil {
			return err
		}

		if err = socks5.Start(context.TODO()); err != nil {
			return err
		}
		logger.Infof("socks5 server started")
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c

		if err = api.Stop(context.TODO()); err != nil {
			logger.WithError(err).Errorf("failed to stop api server")
		} else {
			logger.Debugf("api server terminating")
		}

		if socks5 != nil {
			if err = socks5.Stop(context.TODO()); err != nil {
				logger.WithError(err).Errorf("failed to stop socks5 server")
			} else {
				logger.Debugf("socks5 server terminating")
			}
		}
	}()

	eg := meg.NewAllGroupFunc()

	eg.Go(func() (interface{}, error) {
		er := api.Wait()
		logger.WithError(er).Debugf("api server terminated")
		return nil, er
	})

	if socks5 != nil {
		eg.Go(func() (interface{}, error) {
			er := socks5.Wait()
			logger.WithError(er).Debugf("socks5 server terminated")
			return nil, er
		})
	}

	_, err = eg.Wait()
	logger.WithError(err).Infof("meepo terminated")

	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringP("config", "c", config.GetDefaultConfigPath(), "Location of meepo config file")
	serveCmd.PersistentFlags().BoolP("daemon", "d", true, "Run as daemon")
	serveCmd.PersistentFlags().StringP("identity-file", "i", "", "Select a file from which the identity for Meepo to read")
}
