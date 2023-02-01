package msig

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	builtintypes "github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/spf13/cobra"

	"github.com/labs3/filecoin-wallet-signing/chain/actors"
	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
)

// proposeChangeOwnerCmd represents the msigpropose command
var proposeChangeOwnerCmd = &cobra.Command{
	Use:   "change-owner <multisigAddress> <minerAddress> <newOwnerAddr> ",
	Short: "propose change miner owner ",
	Run: func(cmd *cobra.Command, args []string) {
		proposeChangeOwner(cmd, args)
	},
}

func proposeChangeOwner(cmd *cobra.Command, args []string) {
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

	newOwnerAddr, err := address.NewFromString(args[2])
	if err != nil {
		fmt.Println("decode miner address failed:", err.Error())
		return
	}
	if newOwnerAddr.Protocol() != address.ID {
		newOwnerAddr, err = internal.Lapi.StateLookupID(context.Background(), newOwnerAddr, types.EmptyTSK)
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

	params, err := actors.SerializeParams(&newOwnerAddr)
	if err != nil {
		fmt.Println("Serialize miner.newOwnerAddr failed: ", err)
		return
	}

	proposeParams, err := actors.SerializeParams(&multisig.ProposeParams{
		To:     mnersaddr,
		Method: builtintypes.MethodsMiner.ChangeOwnerAddress,
		Value:  abi.NewTokenAmount(0),
		Params: params,
	})
	if err != nil {
		fmt.Println(err)
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

	fmt.Printf("change miner %v  owner is  %v \n", mnersaddr.String(), newOwnerAddr.String())

}
