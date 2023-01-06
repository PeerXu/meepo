package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"
)

var (
	ncatCmd = &cobra.Command{
		Use:     "ncat",
		Short:   "ncat-like program, support SOCKS5 with authentication",
		Example: "ncat [--proxy <addr>:<port>] [--proxy-type <type>] [--proxy-auth <auth>] <hostname> <port>",
		RunE:    meepoNcat,
		Args:    cobra.ExactArgs(2),
	}

	ncatOptions struct {
		Address string
		Type    string
		Auth    string
	}
)

func meepoNcat(cmd *cobra.Command, args []string) error {
	var auth *proxy.Auth
	var err error

	hostname := args[0]
	port := args[1]

	if ncatOptions.Type != "socks5" {
		return fmt.Errorf("unsupported proxy type: %s", ncatOptions.Type)
	}

	if ncatOptions.Auth != "" {
		ss := strings.SplitN(ncatOptions.Auth, ":", 2)
		if len(ss) != 2 {
			return fmt.Errorf("invalid proxy authentication")
		}

		auth = &proxy.Auth{
			User:     ss[0],
			Password: ss[1],
		}
	}

	dialer, err := proxy.SOCKS5("tcp", ncatOptions.Address, auth, proxy.Direct)
	if err != nil {
		return err
	}

	conn, err := dialer.Dial("tcp", net.JoinHostPort(hostname, port))
	if err != nil {
		return err
	}

	done1 := make(chan struct{})
	go func() {
		defer close(done1)
		io.Copy(conn, os.Stdin) // nolint:errcheck
	}()

	done2 := make(chan struct{})
	go func() {
		defer close(done2)
		io.Copy(os.Stdout, conn) // nolint:errcheck
	}()

	select {
	case <-done1:
	case <-done2:
	}

	return nil
}

func init() {
	fs := ncatCmd.Flags()
	fs.StringVar(&ncatOptions.Address, "proxy", "127.0.0.1:12341", "Specify address of host to proxy through")
	fs.StringVar(&ncatOptions.Type, "proxy-type", "socks5", "Specify proxy type [socks5]")
	fs.StringVar(&ncatOptions.Auth, "proxy-auth", "", "Authenticate with SOCKS proxy server")

	rootCmd.AddCommand(ncatCmd)
}
