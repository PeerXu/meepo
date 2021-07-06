package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mikesmitty/edkey"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	mfs "github.com/PeerXu/meepo/pkg/util/fs"
)

var (
	keygenCmd = &cobra.Command{
		Use:   "keygen [-t ed25519] [-f identity_file]",
		Short: "Generate an identity key",
		RunE:  meepoKeygen,
	}
)

func meepoKeygen(cmd *cobra.Command, args []string) error {
	fs := cmd.Flags()

	filename, _ := fs.GetString("filename")
	overwrite, _ := fs.GetBool("overwrite")
	typ, _ := fs.GetString("type")

	_, err := os.Stat(filename)
	if err == nil {
		if !overwrite {
			return os.ErrExist
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if filename, err = homedir.Expand(filename); err != nil {
		return err
	}

	if err = mfs.EnsureDirectoryExist(filename); err != nil {
		return err
	}

	switch typ {
	case "ed25519":
	default:
		return fmt.Errorf("Unsupported type: %v", typ)
	}

	_, prik, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	buf := pem.EncodeToMemory(&pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(prik),
	})

	if err = ioutil.WriteFile(filename, buf, 0600); err != nil {
		return err
	}

	fmt.Println("Key generated!")

	return nil
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	keygenCmd.PersistentFlags().StringP("type", "t", "ed25519", "Specifies the type of key to create")
	keygenCmd.PersistentFlags().StringP("filename", "f", "meepo.pem", "Specifies the filename of the key file")
	keygenCmd.PersistentFlags().Bool("overwrite", false, "Overwrite an exists key file")
}
