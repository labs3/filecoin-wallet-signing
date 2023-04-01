package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/labs3/filecoin-wallet-signing/cmd/msig"
)

var cfgFile string
var overwrite bool
var minerActor string
var existingBeneficiary bool
var newBeneficiary bool

// rootCmd represents the base command when called without any subcommands generate
var rootCmd = &cobra.Command{
	Use:   "wallet command [options]",
	Short: "Filecoin wallet sign tools",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.filecoin-wallet-signing.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(withdrawCmd)
	rootCmd.AddCommand(changeOwnerCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(verifyCmd)

	rootCmd.AddCommand(changeBeneficiaryCmd)
	changeBeneficiaryCmd.Flags().BoolVar(&overwrite, "overwrite-pending-change", false, "Overwrite the current beneficiary change proposal")
	changeBeneficiaryCmd.Flags().StringVar(&minerActor, "actor", "", "specify the address of miner actor")

	rootCmd.AddCommand(confirmChangeBeneficiaryCmd)
	confirmChangeBeneficiaryCmd.Flags().BoolVar(&existingBeneficiary, "existing-beneficiary", false, "send confirmation from the existing beneficiary address")
	confirmChangeBeneficiaryCmd.Flags().BoolVar(&newBeneficiary, "new-beneficiary", false, "send confirmation from the new beneficiary address")

	rootCmd.AddCommand(msig.Cmd)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".github.com/labs3/filecoin-wallet-signing" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".filecoin-wallet-signing")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
