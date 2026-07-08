package cmd

import (
	"afk/internal/afk"
	"afk/internal/clients/slack"
	"afk/internal/config"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "afk",
	Short:         "Set AKF slack status across workspaces",
	Long:          "afk is a CLI to configure slack status across multiple workspaces",
	Example:       "afk eat 1h\nafk afk 67m",
	Version:       "1.0.0",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		command := args[0]
		var duration string

		if len(args) > 1 {
			duration = args[1]
		}

		// fmt.Printf("Dynamic Command: %s\n", command)
		// fmt.Printf("Dynamic Value:   %s\n", duration)

		config, err := config.New()
		if err != nil {
			return err
		}
		// TODO find tokens from secret
		slackClient := slack.New(map[string]string{})
		afk := afk.New(config, slackClient)

		if strings.EqualFold(command, "clear") {
			afk.ClearStatus()
		} else {
			afk.UpdateStatus(command, duration)
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
