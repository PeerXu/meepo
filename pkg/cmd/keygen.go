package cmd

import (
	"crypto/ed25519"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mikesmitty/edkey"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	keygenCmd = &cobra.Command{
		Use:   "keygen [-i <identity-file>]",
		Short: "Generate an identity file",
		RunE:  meepoKeygen,
		Args:  cobra.NoArgs,
	}
)

func meepoKeygen(cmd *cobra.Command, args []string) error {
	filename, err := cmd.Flags().GetString("identity-file")
	if err != nil {
		return err
	}

	filename, err = homedir.Expand(filename)
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)
	if err == nil {
		return fmt.Errorf("%s exists, remove or backup the key first", filename)
	}

	_, prik, err := ed25519.GenerateKey(nil)
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
	fmt.Printf("Your identity file has been saved in %s\n", filename)

	return nil
}

func init() {
	keygenCmd.PersistentFlags().StringP("identity-file", "i", "mpo_id_ed25519", "Specifies the filename of the key file")
	rootCmd.AddCommand(keygenCmd)
}
