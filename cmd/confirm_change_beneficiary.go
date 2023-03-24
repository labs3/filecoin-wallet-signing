package cmd

import (
	"context"
	"errors"
	"fmt"

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
var confirmChangeBeneficiaryCmd = &cobra.Command{
	Use:   "confirm-change-beneficiary <minerAddress>",
	Short: "confirm change the miner's beneficiary",
	Run: func(cmd *cobra.Command, args []string) {
		err := confirmChangeBeneficiary(cmd, args, existingBeneficiary, newBeneficiary)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func confirmChangeBeneficiary(ccmd *cobra.Command, args []string, existingBeneficiary bool, newBeneficiary bool) error {
	if len(args) < 1 {
		_ = ccmd.Help()
		return errors.New("not enough parameters")
	}

	maddr, err := address.NewFromString(args[0])
	if err != nil {
		return fmt.Errorf("parsing beneficiary address: %w", err)
	}
	ctx := context.Background()
	mi, err := internal.Lapi.StateMinerInfo(ctx, maddr, types.EmptyTSK)
	if err != nil {
		return fmt.Errorf("getting miner info: %w", err)
	}

	if mi.PendingBeneficiaryTerm == nil {
		return fmt.Errorf("no pending beneficiary term found for miner %s", maddr)
	}

	if (existingBeneficiary && newBeneficiary) || (!existingBeneficiary && !newBeneficiary) {
		return fmt.Errorf("must pass exactly one of --existing-beneficiary or --new-beneficiary")
	}

	key, err := pkg.ReadPrivteKey() // existing-beneficiary or new-beneficiary
	if err != nil {
		return fmt.Errorf("decode private key failed: %w", err)
	}

	fmt.Println("Sign address: ", key.Address.String())

	//var fromAddr address.Address
	if existingBeneficiary {
		if mi.PendingBeneficiaryTerm.ApprovedByBeneficiary {
			return fmt.Errorf("beneficiary change already approved by current beneficiary")
		}
		//fromAddr = mi.Beneficiary
	} else {
		if mi.PendingBeneficiaryTerm.ApprovedByNominee {
			return fmt.Errorf("beneficiary change already approved by new beneficiary")
		}
		//fromAddr = mi.PendingBeneficiaryTerm.NewBeneficiary
	}

	fmt.Println("Confirming Pending Beneficiary Term of:")
	fmt.Println("Beneficiary: ", mi.PendingBeneficiaryTerm.NewBeneficiary)
	fmt.Println("Quota:", mi.PendingBeneficiaryTerm.NewQuota)
	fmt.Println("Expiration Epoch:", mi.PendingBeneficiaryTerm.NewExpiration)

	params := &miner.ChangeBeneficiaryParams{
		NewBeneficiary: mi.PendingBeneficiaryTerm.NewBeneficiary,
		NewQuota:       mi.PendingBeneficiaryTerm.NewQuota,
		NewExpiration:  mi.PendingBeneficiaryTerm.NewExpiration,
	}

	sp, err := actors.SerializeParams(params)
	if err != nil {
		return fmt.Errorf("serializing params: %w", err)
	}

	msg := types.Message{
		From:   key.Address,
		To:     maddr,
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
