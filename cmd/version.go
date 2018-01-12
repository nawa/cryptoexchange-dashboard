package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	CommitHash = "N/A"
	BuildDate  = "N/A"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Arch:		 %s/%s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Println("Commit:		", CommitHash)
			fmt.Println("Build Time:	", BuildDate)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
