package cmd

import (
	"afk/internal/afk"
	"afk/internal/clients/slack"
	"afk/internal/config"
	"afk/internal/secrets"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var setup bool
var clear bool

var rootCmd = &cobra.Command{
	Use:           "afk",
	Short:         "Set AKF slack status across workspaces",
	Long:          "afk is a CLI to configure slack status across multiple workspaces",
	Example:       "afk eat 1h\nafk afk 67m",
	Version:       "1.0.0",
	SilenceUsage:  true,
	SilenceErrors: true,
	// Args:          cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if --setup flag was explicitly called
		if cmd.Flags().Changed("setup") {
			fmt.Println("Running setup wizard...")
			// 1. Save client and secret
			// 2. Run authorizer
			// 3. Ask user to build a config
			return nil
		}

		s := secrets.New()
		config, err := config.New()
		if err != nil {
			return err
		}

		slackClient, err := slack.New(s)
		if err != nil {
			return err
		}
		afk := afk.New(config, slackClient)

		if cmd.Flags().Changed("clear") {
			afk.ClearStatus()
			return nil
		}

		if len(args) < 1 || len(args) > 2 {
			return fmt.Errorf("accepts between 1 and 2 arg(s), received %d", len(args))
		}

		preset := args[0]
		var duration string

		if len(args) > 1 {
			duration = args[1]
		}

		err = afk.UpdateStatus(preset, duration)
		if err != nil {
			return err
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

func init() {
	rootCmd.Flags().BoolVarP(&setup, "setup", "s", false, "Setup afk and providers")
	rootCmd.Flags().BoolVarP(&clear, "clear", "c", false, "Clear any/all AFK status")
}
