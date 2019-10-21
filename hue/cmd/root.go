package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	hostname string
	username string
	rootCmd  = &cobra.Command{
		Use:   "office",
		Short: "medien office 455 cli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Medien Office 455 Command Line Tool, by Changkun Ou")
		},
	}
)

// Execute defines the root cmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
