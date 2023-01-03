package cmd

import (
	"fmt"

	C "github.com/PeerXu/meepo/pkg/lib/constant"
	"github.com/PeerXu/meepo/pkg/lib/dialer"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
	"github.com/PeerXu/meepo/pkg/lib/addr"
	simple_sdk "github.com/PeerXu/meepo/pkg/lib/cmd/contrib/simple/sdk"
	"github.com/PeerXu/meepo/pkg/lib/config"
	sdk_interface "github.com/PeerXu/meepo/pkg/meepo/sdk/interface"
	"github.com/spf13/cobra"
)

var (
	teleportCmd = &cobra.Command{
		Use:   "teleport [-m ...] [--source-network <local-network>] [-l <local-address>] [--mode <mode>] <addr> <remote-address>",
		Short: "New transport and teleportation at same time",
		RunE:  meepoTeleport,
		Args:  cobra.ExactArgs(2),
	}

	teleportOptions struct {
		Manual        bool
		SourceNetwork string
		SourceAddress string
		Mode          string
	}
)

func meepoTeleport(cmd *cobra.Command, args []string) error {
	var err error

	targetStr := args[0]
	target, err := addr.FromString(targetStr)
	if err != nil {
		return err
	}

	sinkAddress := args[1]
	sourceNetwork := teleportOptions.SourceNetwork
	sourceAddress := teleportOptions.SourceAddress
	mode := newTeleportationOptions.Mode

	sdk, err := simple_sdk.GetSDK()
	if err != nil {
		return err
	}

	switch sourceNetwork {
	case "socks5":
		sinkAddress = "*"
	}

	opts := []sdk_interface.TeleportOption{
		well_known_option.WithManual(teleportOptions.Manual),
	}
	if teleportOptions.Manual {
		smuxCfg := config.Get().Meepo.Smux
		kcpCfg := config.Get().Meepo.Kcp
		opts = []sdk_interface.NewTransportOption{
			well_known_option.WithEnableMux(!smuxCfg.Disable),
			well_known_option.WithEnableKcp(!kcpCfg.Disable),
		}
		if !smuxCfg.Disable {
			opts = append(opts,
				well_known_option.WithMuxVer(smuxCfg.Version),
				well_known_option.WithMuxBuf(smuxCfg.BufferSize),
				well_known_option.WithMuxStreamBuf(smuxCfg.StreamBufferSize),
			)
		}
		if !kcpCfg.Disable {
			opts = append(opts,
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
	}

	tpv, err := sdk.Teleport(target, dialer.NewAddr(sourceNetwork, sourceAddress), dialer.NewAddr("tcp", sinkAddress), mode, opts...)
	if err != nil {
		return err
	}

	fmt.Printf("teleportation %s created, listen on %s\n", tpv.ID, tpv.SourceAddress)

	return nil
}

func init() {
	fs := teleportCmd.Flags()

	fs.BoolVarP(&teleportOptions.Manual, "manual", "m", false, "specify new transport arguments manually")

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

	fs.StringVar(&teleportOptions.SourceNetwork, "source-network", "tcp", "Source network")
	fs.StringVarP(&teleportOptions.SourceAddress, "listen", "l", "127.0.0.1:0", "Listen address")
	fs.StringVar(&teleportOptions.Mode, "mode", "mux", "Teleportation mode [raw, mux, kcp]")

	bindFlags(fs, []BindFlagsStruct{
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
	})

	rootCmd.AddCommand(teleportCmd)
}
