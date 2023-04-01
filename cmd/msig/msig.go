package msig

import (
	"github.com/spf13/cobra"
)

var overwrite bool
var multisigStr, minerStr string

// Cmd represents the msig command
var Cmd = &cobra.Command{
	Use:   "msig",
	Short: "multisig address tools",
}

func init() {

	Cmd.AddCommand(approveCmd)
	Cmd.AddCommand(inspectCmd)
	proposeCmd.AddCommand(proposeWhithdrawCmd)
	proposeCmd.AddCommand(proposeChangeOwnerCmd)
	proposeCmd.AddCommand(proposeChangeBeneficiaryCmd)
	proposeChangeBeneficiaryCmd.Flags().BoolVar(&overwrite, "overwrite-pending-change", false, "Overwrite the current beneficiary change proposal")
	proposeChangeBeneficiaryCmd.Flags().StringVar(&multisigStr, "msig-addr", "", "The multi signer address of the miner's owner")
	proposeChangeBeneficiaryCmd.Flags().StringVar(&minerStr, "miner-addr", "", "The address of the miner")
	Cmd.AddCommand(proposeCmd)
}
