package cmd

import (
	"afk/internal/afk"
	"afk/internal/clients/slack"
	"afk/internal/config"
	"afk/internal/secrets"
	"fmt"
	"os"
	"strings"

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
		var isSetup bool

		if cmd.Flags().Changed("setup") {
			fmt.Println("Looks like this is your first AFK run, running setup wizard")
			isSetup = true
		}

		s := secrets.New()
		config, err := config.New()
		if err != nil {
			return err
		}

		slackClient, err := slack.New(s, isSetup)
		if err != nil {
			return err
		}
		afk := afk.New(config, slackClient)

		if cmd.Flags().Changed("clear") {
			afk.ClearStatus()
			return nil
		}

		if !isSetup && len(args) < 1 || len(args) > 2 {
			var names []string
			for name := range config.Presets {
				names = append(names, name)
			}

			return fmt.Errorf("afk invoked without a preset name. \n\nYou have following presets configured: [%s]. \nChoose one or create one by running 'afk mypreset'", strings.Join(names, ", "))
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
