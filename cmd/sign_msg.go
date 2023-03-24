package cmd

import (
	"errors"
	"fmt"

	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
	"github.com/spf13/cobra"
)

// signCmd represents the sign message command
var signCmd = &cobra.Command{
	Use:   "sign <message>",
	Short: "sign any string message",
	Run: func(cmd *cobra.Command, args []string) {
		err := sign(cmd, args)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func sign(ccmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		_ = ccmd.Help()
		return errors.New("not enough parameters")
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		return fmt.Errorf("decode private key failed: %s", err.Error())
	}
	fmt.Printf("address: %s\n", key.Address.String())

	fmt.Printf("sign msg: %s\n", args[0])
	msg := internal.AnyMessage{
		From: key.Address,
		Msg:  args[0],
	}

	signature, err := msg.SignMsg(key.KeyInfo.PrivateKey)
	if err != nil {
		return fmt.Errorf("sign failed: %s", err.Error())
	}

	fmt.Printf("signature: %s\n", signature)
	return nil
}