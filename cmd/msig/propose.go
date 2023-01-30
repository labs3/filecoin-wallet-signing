package msig

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/spf13/cobra"

	"github.com/labs3/filecoin-wallet-signing/chain/actors"
	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
)

// proposeCmd represents the msigpropose command
var proposeCmd = &cobra.Command{
	Use:   "propose  multisigAddr toAddr amount",
	Short: "make a proposal",
	Run: func(cmd *cobra.Command, args []string) {
		propose(cmd, args)
	},
}

func propose(ccmd *cobra.Command, args []string) {
	if len(args) < 3 {
		_ = ccmd.Help()
		return
	}

	mtsaddr, err := address.NewFromString(args[0])
	if err != nil {
		fmt.Println("decode multisigAddress failed:", err.Error())
		return
	}

	if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
		fmt.Println("please input a correct multisigAddress")
		return
	}

	acceptAddr, err := address.NewFromString(args[1])
	if err != nil {
		fmt.Println("decode miner address failed:", err.Error())
		return
	}

	sfil, err := types.ParseFIL(args[2])
	if err != nil {
		fmt.Println("The withdrawal amount is wrong or the format is wrong:", err.Error())
		return
	}

	proposeParams, err := actors.SerializeParams(&multisig.ProposeParams{
		To:     acceptAddr,
		Method: builtin.MethodSend,
		Value:  abi.TokenAmount(sfil),
		Params: []byte{},
	})
	if err != nil {
		fmt.Println("actors.SerializeParams &miner2.WithdrawBalanceParams failed: ", err)
		return
	}

	key, err := pkg.ReadPrivteKey()
	if err != nil {
		fmt.Println("decode private key failed: ", err)
		return
	}

	msg := types.Message{
		From:   key.Address,
		To:     mtsaddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtin.MethodsMultisig.Propose,
		Params: proposeParams,
	}

	err = internal.PushSignedMsg(&msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("send from %v to %v amount %v \n", mtsaddr.String(), acceptAddr.String(), pkg.ToFloat64(abi.TokenAmount(sfil)))
}
