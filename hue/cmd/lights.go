package cmd

import (
	"fmt"
	"os"

	"changkun.de/x/pkg/hue/lights"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lightsCmd)

	lightsCmd.Flags().StringVarP(&hostname, "hostname", "H", "", "bridge hostname, or from OFFICE_HOST")
	lightsCmd.Flags().StringVarP(&username, "username", "U", "", "bridge username, or from OFFICE_USER")

	if hostname == "" {
		v, ok := os.LookupEnv("OFFICE_HOST")
		if ok {
			hostname = v
		}
	}
	if username == "" {
		v, ok := os.LookupEnv("OFFICE_USER")
		if ok {
			username = v
		}
	}

	// turn off all lights
	lightsCmd.AddCommand(lightsTurnCmd...)
}

var lightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "Lights control",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Lights control")
	},
}

var lightsTurnCmd = []*cobra.Command{
	{
		Use:   "on",
		Short: "Turn on all lights",
		RunE:  turn,
	},
	{
		Use:   "off",
		Short: "Turn off all lights",
		RunE:  turn,
	},
}

// turn switches every light on the bridge. It uses the command's context, so
// an interrupt cancels the requests in flight, and reports failures through
// RunE rather than panicking.
func turn(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	bridge := lights.NewBridge(hostname, username)

	ls, err := bridge.GetLights(ctx)
	if err != nil {
		return fmt.Errorf("cannot list lights: %w", err)
	}

	on := cmd.Use == "on"
	for _, l := range ls {
		if _, err := l.Turn(ctx, on); err != nil {
			return fmt.Errorf("cannot turn light %d: %w", l.ID, err)
		}
	}
	return nil
}
