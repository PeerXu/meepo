package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	xproxy "golang.org/x/net/proxy"
)

var (
	ncatCmd = &cobra.Command{
		Use:     "ncat",
		Short:   "ncat-like program, support SOCKS5 with authentication",
		Example: "ncat [--proxy <addr>:<port>] [--proxy-type <type>] [--proxy-auth <auth>] <hostname> <port>",
		RunE:    meepoNcat,
		Args:    cobra.ExactArgs(2),
	}
)

func meepoNcat(cmd *cobra.Command, args []string) error {
	var socks5Auth *xproxy.Auth
	var err error

	if len(args) != 2 {
		return fmt.Errorf("require hostname and port")
	}

	fs := cmd.Flags()
	hostname := args[0]
	port := args[1]
	proxy, _ := fs.GetString("proxy")
	proxyType, _ := fs.GetString("proxy-type")
	proxyAuth, _ := fs.GetString("proxy-auth")

	if proxyType != "socks5" {
		return fmt.Errorf("unsupported proxy type: %s", proxyType)
	}

	if proxyAuth != "" {
		ss := strings.SplitN(proxyAuth, ":", 2)
		if len(ss) != 2 {
			return fmt.Errorf("invalid proxy authentication")
		}

		socks5Auth = &xproxy.Auth{
			User:     ss[0],
			Password: ss[1],
		}
	}

	dialer, err := xproxy.SOCKS5("tcp", proxy, socks5Auth, xproxy.Direct)
	if err != nil {
		return err
	}

	conn, err := dialer.Dial("tcp", net.JoinHostPort(hostname, port))
	if err != nil {
		return err
	}

	done := make(chan struct{})
	var doneOnce sync.Once
	go func() {
		defer doneOnce.Do(func() { close(done) })
		io.Copy(conn, os.Stdin)

	}()
	go func() {
		defer doneOnce.Do(func() { close(done) })
		io.Copy(os.Stdout, conn)
	}()

	<-done

	return nil
}

func init() {
	rootCmd.AddCommand(ncatCmd)

	ncatCmd.PersistentFlags().String("proxy", "127.0.0.1:12341", "Specify address of host to proxy through")
	ncatCmd.PersistentFlags().String("proxy-type", "socks5", "Specify proxy type (\"socks5\")")
	ncatCmd.PersistentFlags().String("proxy-auth", "", "Authenticate with SOCKS proxy server")
}
