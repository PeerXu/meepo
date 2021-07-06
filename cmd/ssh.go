package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	msdk "github.com/PeerXu/meepo/pkg/sdk"
)

var (
	sshCmd = &cobra.Command{
		Use:     "ssh",
		Short:   "Run a ssh proxy to other Meepo ssh server",
		Example: "eval $(meepo ssh -- -p <port> ... <username>@<meepo-id>)",
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

// meepo ssh -- -p <ssh-port> ... <username>@<meepo-id>
func meepoSsh(cmd *cobra.Command, args []string) error {
	var err error

	sdk, err := NewHTTPSDK(cmd)
	if err != nil {
		return err
	}

	fs := cmd.Flags()
	binary, _ := fs.GetString("binary")
	id, _ := fs.GetString("id")
	username, _ := fs.GetString("username")
	la, _ := fs.GetString("laddr")
	ra, _ := fs.GetString("raddr")
	rh, rp, _ := net.SplitHostPort(ra)
	strict, _ := fs.GetBool("strict")
	secret, _ := fs.GetString("secret")

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

	remote := MustResolveTCPAddr("", ra)
	tpOpt := &msdk.TeleportOption{
		Name: fmt.Sprintf("ssh:%s:%s", id, rp),
	}
	if la != "" {
		tpOpt.Local = MustResolveTCPAddr("", la)
	}
	if secret != "" {
		tpOpt.Secret = secret
	}
	local, err := sdk.Teleport(id, remote, tpOpt)
	if err != nil {
		return err
	}

	lh, lp, _ := net.SplitHostPort(local.String())
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
	sshStr := strings.Join(buf, " ")

	fmt.Printf("%s\n", sshStr)

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
	sshCmd.PersistentFlags().String("secret", "", "New teleportation with secret")
}
