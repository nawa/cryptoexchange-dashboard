package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "cryptoexchange-dashboard",
	Short: "CryptoExchange Dashboard",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug level")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
