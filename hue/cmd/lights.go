package cmd

import (
	"fmt"
	"os"

	"github.com/changkun/gobase/hue/lights"
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
	&cobra.Command{
		Use:   "on",
		Short: "Turn on all lights",
		Run:   turn,
	},
	&cobra.Command{
		Use:   "off",
		Short: "Turn off all lights",
		Run:   turn,
	},
}

func turn(cmd *cobra.Command, args []string) {
	l := lights.NewBridge(hostname, username)

	ls, err := l.GetLights()
	if err != nil {
		panic(err)
	}

	for _, l := range ls {
		if cmd.Use == "on" {
			l.Turn(true)
		} else {
			l.Turn(false)
		}
	}
}
