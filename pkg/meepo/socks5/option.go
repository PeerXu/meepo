package meepo_socks5

import (
	"github.com/PeerXu/meepo/pkg/lib/option"
	"github.com/PeerXu/meepo/pkg/lib/well_known_option"
)

type NewSocks5ServerOption = option.ApplyOption

func defaultNewSocks5ServerOptions() option.Option {
	return option.NewOption(map[string]any{
		well_known_option.OPTION_HOST: "127.0.0.1",
		well_known_option.OPTION_PORT: "12341",
	})
}
