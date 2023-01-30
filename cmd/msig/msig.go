package msig

import (
	"github.com/spf13/cobra"
)

// Cmd represents the msig command
var Cmd = &cobra.Command{
	Use:   "msig",
	Short: "multisig address tool",
}

func init() {

	Cmd.AddCommand(approveCmd)
	Cmd.AddCommand(inspectCmd)
	proposeCmd.AddCommand(proposeWhithdrawCmd)
	Cmd.AddCommand(proposeCmd)
}
