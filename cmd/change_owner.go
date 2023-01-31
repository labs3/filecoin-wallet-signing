package cmd

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/labs3/filecoin-wallet-signing/chain/actors"
	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
	"github.com/spf13/cobra"
)

// changeOwnerCmd represents the change-owner command
var changeOwnerCmd = &cobra.Command{
	Use:   "change-owner <miner> <newAddr>",
	Short: "change owner of miner",
	Run: func(cmd *cobra.Command, args []string) {
		changeOwner(cmd, args)
	},
}

func changeOwner(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		_ = cmd.Help()
		return
	}

	minerAddr, err := address.NewFromString(args[0])
	if err != nil {
		fmt.Println("decode address failed:", err.Error())
		return
	}

	if minerAddr.Protocol() != address.Actor && minerAddr.Protocol() != address.ID {
		fmt.Println("please input a correct minerAddress")
		return
	}

	newAddr, err := address.NewFromString(args[1])
	if err != nil {
		fmt.Println("decode new  address failed:", err.Error())
		return
	}
	if newAddr.Protocol() != address.ID {
		newAddr, err = internal.Lapi.StateLookupID(context.Background(), newAddr, types.EmptyTSK)
		if err != nil {
			fmt.Println("query new address ID failed:", err.Error())
			return
		}
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		fmt.Println("decode private key failed: ", err)
		return
	}

	params, err := actors.SerializeParams(&newAddr)
	if err != nil {
		fmt.Println("Serialize miner.WithdrawBalanceParams failed: ", err)
		return
	}

	msg := types.Message{
		From:   key.Address,
		To:     minerAddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtin.MethodsMiner.ChangeOwnerAddress,
		Params: params,
	}
	err = internal.PushSignedMsg(&msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}

}
