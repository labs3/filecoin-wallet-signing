package cmd

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin"
	"github.com/spf13/cobra"

	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send toAddress Amount",
	Short: "send",
	Run: func(cmd *cobra.Command, args []string) {
		send(cmd, args)
	},
}

func send(ccmd *cobra.Command, args []string) {
	if len(args) < 2 {
		_ = ccmd.Help()
		return
	}

	sAddr, err := address.NewFromString(args[0])
	if err != nil {
		fmt.Println("decode address failed:", err.Error())
		return
	}

	sFil, err := types.ParseFIL(args[1])
	if err != nil {
		fmt.Println("incorrect amount or format: ", err.Error())
		return
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		fmt.Println("decode private key failed: ", err)
		return
	}

	fmt.Printf("send from %v to %v amount %v \n", key.Address.String(), sAddr.String(), pkg.ToFloat64(abi.TokenAmount(sFil)))

	msg := types.Message{
		From:   key.Address,
		To:     sAddr,
		Value:  abi.TokenAmount(sFil),
		Method: builtin.MethodSend,
	}
	err = internal.PushSignedMsg(&msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}
}
