package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "afk",
	Short:   "Set AKF slack status across workspaces",
	Long:    "afk is a CLI to configure slack status across multiple workspaces",
	Example: "afk eat 1h\nafk afk 67m",
	Version: "1.0.0",
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]
		var duration string

		if len(args) > 1 {
			duration = args[1]
		}

		fmt.Printf("Dynamic Command: %s\n", command)
		fmt.Printf("Dynamic Value:   %s\n", duration)
		// afk := afk.New()
		// afk.GetConfig()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
