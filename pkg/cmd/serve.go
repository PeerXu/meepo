package cmd

import (
	"context"
	"crypto/ed25519"
	"net"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"

	"github.com/PeerXu/meepo/pkg/lib/acl"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	simple_logger "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/logger"
	"github.com/PeerXu/meepo/pkg/lib/config"
	C "github.com/PeerXu/meepo/pkg/lib/constant"
	crypto_core "github.com/PeerXu/meepo/pkg/lib/crypto/core"
	"github.com/PeerXu/meepo/pkg/lib/daemon"
	"github.com/PeerXu/meepo/pkg/lib/logging"
	"github.com/PeerXu/meepo/pkg/lib/marshaler"
	marshaler_json "github.com/PeerXu/meepo/pkg/lib/marshaler/json"
	"github.com/PeerXu/meepo/pkg/lib/pprof"
	"github.com/PeerXu/meepo/pkg/lib/rpc"
	rpc_core "github.com/PeerXu/meepo/pkg/lib/rpc/core"
	rpc_http "github.com/PeerXu/meepo/pkg/lib/rpc/http"
	mpo_webrtc "github.com/PeerXu/meepo/pkg/lib/webrtc"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	meepo_core "github.com/PeerXu/meepo/pkg/meepo/core"
	meepo_interface "github.com/PeerXu/meepo/pkg/meepo/interface"
	socks5 "github.com/PeerXu/meepo/pkg/meepo/socks5"
	"github.com/PeerXu/meepo/pkg/meepo/tracker"
	tracker_core "github.com/PeerXu/meepo/pkg/meepo/tracker/core"
)

var (
	serveCmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"summon"},
		Short:   "Summon a Meepo",
		RunE:    meepoSummon,
		Args:    cobra.NoArgs,
	}
)

func meepoSummon(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cfg := config.Get()

	if cfg.Meepo.Daemon {
		daemon.Daemon()
	}

	if !slices.Contains([]string{"trace", "debug"}, cfg.Meepo.Log.Level) {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	logger, err := simple_logger.GetLogger()
	if err != nil {
		return err
	}

	fs := cmd.Flags()
	summonLogger := logger.WithField("#method", "meepoSummon")

	aclFlag := fs.Lookup("acl")
	switch cfg.Meepo.Mode {
	case C.MODE_MAIN:
	case C.MODE_MINOR:
		if aclFlag.Changed || cfg.Meepo.Acl != C.ACL_BLOCK_ALL {
			summonLogger.Warningf("failed to apply acl in minor mode, please use main mode to apply acl")
			return nil
		}
	case C.MODE_DEV:
		if aclFlag.Changed || cfg.Meepo.Acl != C.ACL_ALLOW_ALL {
			summonLogger.Warningf("failed to apply acl in dev mode, force set to allow all request")
			cfg.Meepo.Acl = C.ACL_ALLOW_ALL
		}
	}

	if cfg.Meepo.Pprof != "" {
		pprof.Setup(cfg.Meepo.Pprof)
		summonLogger.Debugf("pprof listen on %s", cfg.Meepo.Pprof)
	}

	var pubk ed25519.PublicKey
	var prik ed25519.PrivateKey
	if cfg.Meepo.Identity.NoFile || cfg.Meepo.Identity.File == "" {
		pubk, prik, err = ed25519.GenerateKey(nil)
		if err != nil {
			summonLogger.WithError(err).Errorf("failed to generate ed25519 key")
			return err
		}
	} else {
		pubk, prik, err = crypto_core.LoadEd25519Key(cfg.Meepo.Identity.File)
		if err != nil {
			summonLogger.WithError(err).WithField("file", cfg.Meepo.Identity.File).Errorf("failed to load ed25519 key")
			return err
		}
	}

	signer := crypto_core.NewSigner(pubk, prik)
	cryptor := crypto_core.NewCryptor(pubk, prik, nil)
	mpAddr, err := addr.FromBytesWithoutMagicCode(pubk)
	if err != nil {
		summonLogger.WithError(err).Errorf("failed to parse ed25519 public key to addr")
		return err
	}
	summonLogger = summonLogger.WithField("addr", mpAddr.String())
	iceServers, err := mpo_webrtc.ParseICEServers(cfg.Meepo.Webrtc.IceServers)
	if err != nil {
		summonLogger.WithError(err).Errorf("failed to parse ice servers")
		return err
	}
	webrtcConfiguration := webrtc.Configuration{
		ICEServers: iceServers,
	}

	var tks []tracker_core.Tracker
	for _, tkCfg := range cfg.Meepo.Trackers {
		var tk tracker_core.Tracker
		var name string

		if tkCfg.Name == "skip" {
			continue
		}

		tkAddr, err := addr.FromString(tkCfg.Addr)
		if err != nil {
			summonLogger.WithError(err).Errorf("failed to parse tracker addr")
			return err
		}

		newTkOpts := []tracker_core.NewTrackerOption{}
		switch tkCfg.Name {
		case "rpc":
			name = "rpc"
			var callerName string
			var newCallerOpts []rpc_core.NewCallerOption
			switch tkCfg.CallerName {
			case "http":
				callerName = "http"
				newCallerOpts = append(newCallerOpts,
					well_known_option.WithLogger(logger),
					crypto_core.WithSigner(signer),
					crypto_core.WithCryptor(cryptor),
					marshaler.WithMarshaler(marshaler_json.Marshaler),
					marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
					rpc_http.WithBaseURL("http://"+tkCfg.Host),
				)
			}
			caller, err := rpc.NewCaller(callerName, newCallerOpts...)
			if err != nil {
				summonLogger.WithError(err).Errorf("failed to new caller")
				return err
			}
			newTkOpts = append(newTkOpts,
				well_known_option.WithAddr(tkAddr),
				rpc_core.WithCaller(caller),
			)
		}
		tk, err = tracker.NewTracker(name, newTkOpts...)
		if err != nil {
			logger.WithError(err).Errorf("failed to new tracker")
			return err
		}
		tks = append(tks, tk)
	}

	acl_, err := acl.FromString(cfg.Meepo.Acl)
	if err != nil {
		return err
	}

	smuxCfg := config.Get().Meepo.Smux
	kcpCfg := config.Get().Meepo.Kcp
	nmOpts := []meepo_core.NewMeepoOption{
		well_known_option.WithAddr(mpAddr),
		well_known_option.WithLogger(logger),
		tracker_core.WithTrackers(tks...),
		crypto_core.WithSigner(signer),
		crypto_core.WithCryptor(cryptor),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
		well_known_option.WithWebrtcConfiguration(webrtcConfiguration),
		acl.WithAcl(acl_),
		well_known_option.WithEnableMux(!smuxCfg.Disable),
		well_known_option.WithEnableKcp(!kcpCfg.Disable),
	}
	if !smuxCfg.Disable {
		nmOpts = append(nmOpts,
			well_known_option.WithMuxVer(smuxCfg.Version),
			well_known_option.WithMuxBuf(smuxCfg.BufferSize),
			well_known_option.WithMuxStreamBuf(smuxCfg.StreamBufferSize),
			well_known_option.WithMuxNocomp(smuxCfg.Nocomp),
		)
	}
	if !kcpCfg.Disable {
		nmOpts = append(nmOpts,
			well_known_option.WithKcpPreset(kcpCfg.Preset),
			well_known_option.WithKcpCrypt(kcpCfg.Crypt),
			well_known_option.WithKcpKey(kcpCfg.Key),
			well_known_option.WithKcpMtu(kcpCfg.Mtu),
			well_known_option.WithKcpSndwnd(kcpCfg.Sndwnd),
			well_known_option.WithKcpRecvwnd(kcpCfg.Rcvwnd),
			well_known_option.WithKcpDataShard(kcpCfg.DataShard),
			well_known_option.WithKcpParityShard(kcpCfg.ParityShard),
		)
	}

	mp, err := meepo_core.NewMeepo(nmOpts...)
	if err != nil {
		summonLogger.WithError(err).Errorf("failed to new meepo")
		return err
	}
	defer mp.Close(ctx)

	var name string
	apiServerLogger := logger.WithFields(logging.Fields{
		"name": cfg.Meepo.API.Name,
	})
	newAPIOpts := []rpc_core.NewServerOption{
		rpc_core.WithHandler(mp.AsAPIHandler()),
		well_known_option.WithLogger(logger),
		marshaler.WithMarshaler(marshaler_json.Marshaler),
		marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
	}
	switch cfg.Meepo.API.Name {
	case "http":
		name = "simple_http"
		lis, err := net.Listen("tcp", cfg.Meepo.API.Host)
		if err != nil {
			summonLogger.WithError(err).Errorf("failed to listen")
			return err
		}
		newAPIOpts = append(newAPIOpts, well_known_option.WithListener(lis))
		apiServerLogger = apiServerLogger.WithFields(logging.Fields{
			"host": cfg.Meepo.API.Host,
		})
	}
	apiSrv, err := rpc.NewServer(name, newAPIOpts...)
	if err != nil {
		summonLogger.WithError(err).Errorf("failed to new api server")
		return err
	}
	go apiSrv.Serve(ctx)
	defer apiSrv.Terminate(ctx) // nolint:errcheck
	apiServerLogger.Infof("api server started")

	if cfg.Meepo.Socks5.Host != "" {
		socks5Logger := logger
		lis, err := net.Listen("tcp", cfg.Meepo.Socks5.Host)
		if err != nil {
			summonLogger.WithError(err).Errorf("failed to listen")
			return err
		}
		socks5Srv, err := socks5.NewSocks5Server(
			well_known_option.WithLogger(logger),
			meepo_interface.WithMeepo(mp),
			well_known_option.WithListener(lis),
		)
		if err != nil {
			summonLogger.WithError(err).Errorf("failed to new socks5 server")
			return err
		}
		go socks5Srv.Serve(ctx)
		defer socks5Srv.Terminate(ctx) // nolint:errcheck
		socks5Logger.WithFields(logging.Fields{
			"host": cfg.Meepo.Socks5.Host,
		}).Infof("socks5 server started")
	}

	for _, tkdCfg := range cfg.Meepo.Trackerds {
		tkdLogger := logger

		switch tkdCfg.Name {
		case "rpc":
			var name string
			newSrvOpts := []rpc_core.NewServerOption{
				well_known_option.WithLogger(logger),
				rpc_core.WithHandler(mp.AsTrackerdHandler()),
				crypto_core.WithSigner(signer),
				crypto_core.WithCryptor(cryptor),
				marshaler.WithMarshaler(marshaler_json.Marshaler),
				marshaler.WithUnmarshaler(marshaler_json.Unmarshaler),
			}
			switch tkdCfg.ServerName {
			case "http":
				name = "http"
				lis, err := net.Listen("tcp", tkdCfg.Host)
				if err != nil {
					logger.WithError(err).Errorf("failed to listen")
					return err
				}
				newSrvOpts = append(newSrvOpts, well_known_option.WithListener(lis))
				tkdLogger = tkdLogger.WithFields(logging.Fields{
					"name": name,
					"host": tkdCfg.Host,
				})
			}
			srv, err := rpc.NewServer(name, newSrvOpts...)
			if err != nil {
				logger.WithError(err).Errorf("failed to new rpc server")
				return err
			}
			go srv.Serve(ctx)
			defer srv.Terminate(ctx) // nolint:errcheck
			tkdLogger.Infof("trackerd server started")
		}
	}

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)

	s := <-c
	logger.WithField("signal", s).Infof("catch signal")

	go func() {
		c := make(chan os.Signal, 1)
		defer close(c)
		signal.Notify(c, os.Interrupt)

		s := <-c
		logger.WithField("signal", s).Warningf("catch signal, force stop")

		os.Exit(1)
	}()

	return nil
}

// nolint:errcheck
func init() {
	fs := serveCmd.Flags()

	fs.BoolVarP(&config.Get().Meepo.Daemon, "daemon", "d", true, "run as daemon")
	fs.StringVar(&config.Get().Meepo.Acl, "acl", "", "access control list")

	fs.StringVar(&config.Get().Meepo.Pprof, "pprof", "", "profile listen address")

	webrtcCfg := &config.Get().Meepo.Webrtc
	fs.Uint32Var(&webrtcCfg.RecvBufferSize, "sockBuf", C.WEBRTC_RECEIVE_BUFFER_SIZE, "receive buffer in bytes/per webrtc connection")

	idCfg := &config.Get().Meepo.Identity
	fs.BoolVarP(&idCfg.NoFile, "no-identity-file", "n", false, "no identity file")
	fs.StringVarP(&idCfg.File, "identity-file", "i", "", "identity file")

	smuxCfg := &config.Get().Meepo.Smux
	fs.BoolVar(&smuxCfg.Disable, "disableMux", false, "disable Mux mode")
	fs.IntVar(&smuxCfg.Version, "muxVer", C.SMUX_VERSION, "specify smux version [1, 2]")
	fs.IntVar(&smuxCfg.BufferSize, "muxBuf", C.SMUX_BUFFER_SIZE, "the overall de-mux buffer in bytes")
	fs.IntVar(&smuxCfg.StreamBufferSize, "muxStreamBuf", C.SMUX_STREAM_BUFFER_SIZE, "per stream receive buffer in bytes, smux v2+")
	fs.BoolVar(&smuxCfg.Nocomp, "muxNocomp", C.SMUX_NOCOMP, "disable compression")

	kcpCfg := &config.Get().Meepo.Kcp
	fs.BoolVar(&kcpCfg.Disable, "disableKcp", false, "disable Kcp mode")
	fs.StringVar(&kcpCfg.Preset, "kcpPreset", C.KCP_PRESET, "presets: fast3, fast2, fast, normal")
	fs.StringVar(&kcpCfg.Crypt, "kcpCrypt", C.KCP_CRYPT, "crypt algorithms [aes, aes-128, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, sm4, none]")
	fs.StringVar(&kcpCfg.Key, "kcpKey", C.KCP_KEY, "pre-shared secret between client and server")
	fs.IntVar(&kcpCfg.Mtu, "kcpMtu", C.KCP_MTU, "set maximum transmission unit for packets")
	fs.IntVar(&kcpCfg.Sndwnd, "kcpSndwnd", C.KCP_SNDWND, "set send window size(num of packets)")
	fs.IntVar(&kcpCfg.Rcvwnd, "kcpRcvwnd", C.KCP_RCVWND, "set receive window size(num of packets)")
	fs.IntVar(&kcpCfg.DataShard, "kcpDataShard", C.KCP_DATA_SHARD, "set reed-solomon erasure coding - datashard")
	fs.IntVar(&kcpCfg.ParityShard, "kcpParityShard", C.KCP_PARITY_SHARD, "set reed-solomon erasure coding - parityshard")

	for _, c := range []struct {
		configKey  string
		commandKey string
	}{
		{"meepo.daemon", "daemon"},
		{"meepo.acl", "acl"},

		{"meepo.pprof", "pprof"},

		{"meepo.identity.no_file", "no-identity-file"},
		{"meepo.identity.file", "identity-file"},

		{"meepo.smux.disable", "disableMux"},
		{"meepo.smux.version", "muxVer"},
		{"meepo.smux.bufferSize", "muxBuf"},
		{"meepo.smux.streamBufferSize", "muxStreamBuf"},
		{"meepo.smux.keepalive", "muxKeepalive"},
		{"meepo.smux.nocomp", "mucNocomp"},

		{"meepo.kcp.disable", "disableKcp"},
		{"meepo.kcp.preset", "kcpPreset"},
		{"meepo.kcp.crypt", "kcpCrypt"},
		{"meepo.kcp.key", "kcpKey"},
		{"meepo.kcp.mtu", "kcpMtu"},
		{"meepo.kcp.sndwnd", "kcpSndwnd"},
		{"meepo.kcp.rcvwnd", "kcpRcvwnd"},
		{"meepo.kcp.dataShard", "kcpDataShard"},
		{"meepo.kcp.parityShard", "kcpParityShard"},
	} {
		viper.BindPFlag(c.configKey, fs.Lookup(c.commandKey))
	}

	rootCmd.AddCommand(serveCmd)
}
