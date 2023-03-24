package cmd

import (
	"errors"
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify signature command
var verifyCmd = &cobra.Command{
	Use:   "verify <address> <message> <signature>",
	Short: "verify the signature of any string message",
	Run: func(cmd *cobra.Command, args []string) {
		err := verify(cmd, args)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func verify(ccmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		_ = ccmd.Help()
		return errors.New("not enough parameters")
	}

	address, err := address.NewFromString(args[0])
	if err != nil {
		return err
	}

	msg := internal.AnyMessage{
		From: address,
		Msg:  args[1],
		Sig:  args[2],
	}

	err = msg.VerifyMsg()
	if err != nil {
		return err
	}

	fmt.Println("verify success !")
	return nil
}