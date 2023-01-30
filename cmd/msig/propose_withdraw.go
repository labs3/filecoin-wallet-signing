package msig

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	builtintypes "github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/spf13/cobra"

	"github.com/labs3/filecoin-wallet-signing/chain/actors"
	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
)

// proposeCmd represents the msigpropose command
var proposeWhithdrawCmd = &cobra.Command{
	Use:   "withdraw <multisigAddress> <minerAddress> <amount> ",
	Short: "propose withdraw from miner ",
	Run: func(cmd *cobra.Command, args []string) {
		proposeWithdraw(cmd, args)
	},
}

func proposeWithdraw(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		_ = cmd.Help()
		return
	}

	mtsaddr, err := address.NewFromString(args[0])
	if err != nil {
		fmt.Println("decode address failed:", err.Error())
		return
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return
	}

	mnersaddr, err := address.NewFromString(args[1])
	if err != nil {
		fmt.Println("decode miner address failed:", err.Error())
		return
	}

	if mnersaddr.Protocol() != address.Actor && mnersaddr.Protocol() != address.ID {
		fmt.Println("please input a correct miner address")
		return
	}

	wdfil, err := types.ParseFIL(args[2])
	if err != nil {
		fmt.Println("The withdrawal amount is wrong or the format is wrong: ", err.Error())
		return
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		fmt.Println("decode private key failed: ", err)
		return
	}

	withdrawBalanceParams, err := actors.SerializeParams(&miner.WithdrawBalanceParams{
		AmountRequested: abi.TokenAmount(wdfil), // Default to attempting to withdraw all the extra funds in the miner actor
	})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return
	}

	proposeParams, err := actors.SerializeParams(&multisig.ProposeParams{
		To:     mnersaddr,
		Method: builtintypes.MethodsMiner.WithdrawBalance,
		Value:  abi.NewTokenAmount(0),
		Params: withdrawBalanceParams,
	})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return
	}

	msg := types.Message{
		From:   key.Address,
		To:     mtsaddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtintypes.MethodsMultisig.Propose,
		Params: proposeParams,
	}

	err = internal.PushSignedMsg(&msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("withdraw %v FIL from %v \n", pkg.ToFloat64(abi.TokenAmount(wdfil)), mtsaddr.String())

}
