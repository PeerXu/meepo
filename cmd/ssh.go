package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	http_api "github.com/PeerXu/meepo/pkg/api/http"
)

var (
	sshCmd = &cobra.Command{
		Use:     "ssh",
		Short:   "Run a ssh proxy to other Meepo ssh server",
		Example: "eval $(meepo ssh -- <username>@<meepo-id> -p <port> ...)",
		RunE:    meepoSsh,
	}
)

func splitUsernameHost(s string) (string, string, error) {
	if !strings.ContainsRune(s, '@') {
		return "", s, nil
	}

	ss := strings.SplitN(s, "@", 2)
	if len(ss) > 2 {
		return "", "", fmt.Errorf("bad destination")
	}

	return ss[0], ss[1], nil
}

// meepo ssh -- <username>@<meepo-id> -p <ssh-port> ...
func meepoSsh(cmd *cobra.Command, args []string) error {
	var err error

	fs := cmd.Flags()

	binary, _ := fs.GetString("binary")
	host, _ := fs.GetString("host")
	id, _ := fs.GetString("id")
	username, _ := fs.GetString("username")
	la, _ := fs.GetString("laddr")
	ra, _ := fs.GetString("raddr")
	rh, rp, _ := net.SplitHostPort(ra)
	strict, _ := fs.GetBool("strict")

	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			if len(args[i]) >= 2 && args[i][:2] == "-p" {
				if len(args[i]) > 3 {
					if args[i][2] == '=' {
						rp = args[i][3:]
					} else {
						rp = args[i][2:]
					}
				} else {
					rp = args[i+1]
					args[i+1] = ""
				}
				args[i] = ""
				break
			}
		}
		username, id, _ = splitUsernameHost(args[len(args)-1])
	}
	ra = net.JoinHostPort(rh, rp)

	targetUrl, err := url.Parse(host)
	if err != nil {
		return err
	}

	targetUrl.Path = "/v1/actions/teleport"

	client := resty.New()
	res, err := client.R().
		SetBody(map[string]interface{}{
			"id":            id,
			"remoteNetwork": "tcp",
			"remoteAddress": ra,
			"localNetwork":  "tcp",
			"localAddress":  la,
		}).
		Post(targetUrl.String())
	if err != nil {
		return err
	}

	var tr http_api.TeleportResponse
	if err = json.NewDecoder(bytes.NewReader(res.Body())).Decode(&tr); err != nil {
		return err
	}

	lh, lp, _ := net.SplitHostPort(tr.LocalAddress)
	if username != "" {
		username += "@"
	}
	buf := []string{binary, fmt.Sprintf("-p%s", lp)}
	if !strict {
		buf = append(buf, "-o 'StrictHostKeyChecking=no'", "-o 'UserKnownHostsFile=/dev/null'")
	}

	if len(args) > 1 {
		buf = append(buf, args[:len(args)-2]...)
	}
	buf = append(buf, fmt.Sprintf("%s%s", username, lh))

	fmt.Println(strings.Join(buf, " "))

	return nil
}

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.PersistentFlags().String("binary", "ssh", "Binary of ssh")
	sshCmd.PersistentFlags().String("id", "", "Meepo ID")
	sshCmd.PersistentFlags().String("username", "", "Username for ssh")
	sshCmd.PersistentFlags().String("laddr", "127.0.0.1:0", "Local listen address, if not set, random port will be listen")
	sshCmd.PersistentFlags().String("raddr", "127.0.0.1:22", "Remote ssh server address, port overwrite by rest option(-p)")
	sshCmd.PersistentFlags().Bool("strict", false, "Check host key and write to known hosts file")
}
