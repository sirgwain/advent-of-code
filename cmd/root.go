package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/phsym/console-slog"
	"github.com/spf13/cobra"
)

var debug bool

func logPreRun(cmd *cobra.Command, args []string) error {
	level := slog.LevelWarn
	if debug {
		level = slog.LevelDebug
	}
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{Level: level, TimeFormat: time.TimeOnly}),
	)
	slog.SetDefault(logger)

	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "advent-of-code-2025",
	Short:             "advent-of-code solutions for 2025",
	PersistentPreRunE: logPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		// Show usage
		cmd.Help()
		os.Exit(1)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		println(err)
		os.Exit(1)
	}
}

func init() {
	// all commands have debug mode
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
}
