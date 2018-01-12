package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	mongoURL string

	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Syncs your wallet data",
		Long:  "Syncs your wallet data with database in background",
		RunE:  commandSyncRun,
	}
)

func init() {
	syncCmd.Flags().StringVarP(&mongoURL, "db-url", "u", "", "Url to MongoDB")

	rootCmd.AddCommand(syncCmd)
}

func commandSyncRun(cmd *cobra.Command, args []string) error {
	fmt.Println(mongoURL)
	return nil
}
