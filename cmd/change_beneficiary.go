package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v9/miner"
	"github.com/labs3/filecoin-wallet-signing/chain/actors"
	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
	"github.com/spf13/cobra"
)

// changeBeneficiaryCmd represents the changing of beneficiary
var changeBeneficiaryCmd = &cobra.Command{
	Use:   "change-beneficiary <beneficiaryAddress> <quota> <expiration>",
	Short: "change the miner's beneficiary",
	Run: func(cmd *cobra.Command, args []string) {
		err := changeBeneficiary(cmd, args, overwrite, minerActor)
		if err != nil {
			fmt.Println(err)
		}
	},
}


func changeBeneficiary(ccmd *cobra.Command, args []string, overwrite bool, minerActor string) error {
	if len(args) < 3 {
		_ = ccmd.Help()
		return errors.New("not enough parameters")
	}

	ba, err := address.NewFromString(args[0])
	if err != nil {
		return fmt.Errorf("parsing beneficiary address: %w", err)
	}

	ctx := context.Background()
	newAddr, err := internal.Lapi.StateLookupID(ctx, ba, types.EmptyTSK)
	if err != nil {
		return fmt.Errorf("looking up new beneficiary address: %w", err)
	}

	quota, err := types.ParseFIL(args[1])
	if err != nil {
		return fmt.Errorf("parsing quota: %w", err)
	}

	expiration, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return fmt.Errorf("parsing expiration: %w", err)
	}

	mActor, err := address.NewFromString(minerActor)
	if err != nil {
		return errors.New("must specify miner actor address")
	}

	key, err := pkg.ReadPrivteKey() // owner propose
	if err != nil {
		return fmt.Errorf("decode private key failed: %w", err)
	}

	mi, err := internal.Lapi.StateMinerInfo(ctx, mActor, types.EmptyTSK)
	if err != nil {
		return fmt.Errorf("getting miner info: %w", err)
	}

	if mi.Beneficiary == mi.Owner && newAddr == mi.Owner {
		return fmt.Errorf("beneficiary %s already set to owner address", mi.Beneficiary)
	}

	if mi.PendingBeneficiaryTerm != nil {
		fmt.Println("WARNING: replacing Pending Beneficiary Term of:")
		fmt.Println("Beneficiary: ", mi.PendingBeneficiaryTerm.NewBeneficiary)
		fmt.Println("Quota:", mi.PendingBeneficiaryTerm.NewQuota)
		fmt.Println("Expiration Epoch:", mi.PendingBeneficiaryTerm.NewExpiration)

		if !overwrite {
			return fmt.Errorf("must pass --overwrite-pending-change to replace current pending beneficiary change. Please review CAREFULLY")
		}
	}

	params := &miner.ChangeBeneficiaryParams{
		NewBeneficiary: newAddr,
		NewQuota:       abi.TokenAmount(quota),
		NewExpiration:  abi.ChainEpoch(expiration),
	}

	sp, err := actors.SerializeParams(params)
	if err != nil {
		return fmt.Errorf("serializing params: %w", err)
	}

	msg := types.Message{
		From:   key.Address,
		To:     mActor,
		Value:  abi.NewTokenAmount(0),
		Method: builtin.MethodsMiner.ChangeBeneficiary,
		Params: sp,
	}

	err = internal.PushSignedMsg(&msg, key.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
