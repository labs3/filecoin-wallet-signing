package cmd

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/miner"
	"github.com/spf13/cobra"

	"github.com/labs3/filecoin-wallet-signing/chain/actors"
	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
)

// withdrawCmd represents the withdraw command
var withdrawCmd = &cobra.Command{
	Use:   "withdraw <miner> <amount>",
	Short: "withdraw from miner",
	Run: func(cmd *cobra.Command, args []string) {
		withdraw(cmd, args)
	},
}

func withdraw(cmd *cobra.Command, args []string) {

	if len(args) != 2 {
		fmt.Println("Parameter error, please check the parameter")
		return
	}

	mineraddr, err := address.NewFromString(args[0])
	if err != nil {
		fmt.Println("invalid address: ", err.Error())
		return
	}

	wdfil, err := types.ParseFIL(args[1])
	if err != nil {
		fmt.Println("The withdrawal amount is wrong or the format is wrong: ", err.Error())
		return
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		fmt.Println("decode private key failed: ", err)
		return
	}

	params, err := actors.SerializeParams(&miner.WithdrawBalanceParams{
		AmountRequested: abi.TokenAmount(wdfil), // Default to attempting to withdraw all the extra funds in the miner actor
	})
	if err != nil {
		fmt.Println("Serialize miner.WithdrawBalanceParams failed: ", err)
		return
	}

	msg := types.Message{
		From:   key.Address,
		To:     mineraddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtin.MethodsMiner.WithdrawBalance,
		Params: params,
	}
	err = internal.PushSignedMsg(&msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}
}
