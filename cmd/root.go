package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cryptoexchange-wallet-info",
	Short: "CryptoExchange Wallet Info",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
