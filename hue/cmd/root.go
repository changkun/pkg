package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

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

// Execute runs the root command. The context is cancelled on interrupt, so a
// request to the bridge stops with the process instead of running to its
// timeout.
func Execute() {
	os.Exit(run())
}

// run is separate from Execute so that the signal handler is released before
// the process exits; os.Exit does not run deferred calls.
func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
