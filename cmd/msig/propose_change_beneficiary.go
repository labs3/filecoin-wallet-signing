package msig

import (
	"context"
	"errors"
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	builtintypes "github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v9/miner"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/labs3/filecoin-wallet-signing/chain/actors"
	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
	"github.com/spf13/cobra"
)

// proposeChangeBeneficiaryCmd represents the changing of beneficiary
var proposeChangeBeneficiaryCmd = &cobra.Command{
	Use:   "change-beneficiary <beneficiaryAddress> <quota> <expiration>",
	Short: "change the miner's beneficiary",
	Run: func(cmd *cobra.Command, args []string) {
		err := changeBeneficiary(cmd, args, overwrite, multisigStr, minerStr)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func changeBeneficiary(ccmd *cobra.Command, args []string, overwrite bool, multisigStr, minerStr string) error {
	if len(args) != 3 {
		_ = ccmd.Help()
		return errors.New("not enough parameters")
	}

	ba, err := address.NewFromString(args[0])
	if err != nil {
		return fmt.Errorf("parsing beneficiary address: %w", err)
	}

	ctx := context.Background()
	newBeneficiary, err := internal.Lapi.StateLookupID(ctx, ba, types.EmptyTSK)
	if err != nil {
		return fmt.Errorf("looking up new beneficiary address: %w", err)
	}

	quota, err := types.ParseFIL(args[1])
	if err != nil {
		return fmt.Errorf("parsing quota: %w", err)
	}

	expiration, err := types.BigFromString(args[2])
	if err != nil {
		return fmt.Errorf("parsing expiration: %w", err)
	}

	multisigAddr, err := address.NewFromString(multisigStr)
	if err != nil {
		return fmt.Errorf("parsing multisig address: %w", err)
	}

	if multisigAddr.Protocol() != address.Actor && multisigAddr.Protocol() != address.ID {
		return errors.New("please input a correct multisig address, Actor/ID address")
	}

	minerAddr, err := address.NewFromString(minerStr)
	if err != nil {
		return fmt.Errorf("parsing miner address: %w", err)
	}

	if minerAddr.Protocol() != address.Actor && minerAddr.Protocol() != address.ID {
		return errors.New("please input a correct miner address, Actor/ID address")
	}

	mi, err := internal.Lapi.StateMinerInfo(ctx, minerAddr, types.EmptyTSK)
	if err != nil {
		return err
	}

	sender, err := pkg.ReadPrivteKey() // sender, the one of signers
	if err != nil {
		return fmt.Errorf("decode private key failed: %w", err)
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
		NewBeneficiary: newBeneficiary,
		NewQuota:       abi.TokenAmount(quota),
		NewExpiration:  abi.ChainEpoch(expiration.Int64()),
	}

	sp, err := actors.SerializeParams(params)
	if err != nil {
		return fmt.Errorf("serializing params: %w", err)
	}

	enc, actErr := actors.SerializeParams(&multisig.ProposeParams{
		To:     minerAddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtintypes.MethodsMiner.ChangeBeneficiary,
		Params: sp,
	})
	if actErr != nil {
		return fmt.Errorf("failed to serialize parameters: %w", actErr)
	}

	msg := types.Message{
		From:   sender.Address,
		To:     multisigAddr,
		Value:  abi.NewTokenAmount(0),
		Method: builtintypes.MethodsMultisig.Propose,
		Params: enc,
	}

	err = internal.PushSignedMsg(&msg, sender.PrivateKey)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
