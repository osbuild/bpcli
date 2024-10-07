package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootLogLevel string

func rootPreRunE(cmd *cobra.Command, _ []string) error {
	if rootLogLevel == "" {
		logrus.SetLevel(logrus.ErrorLevel)
		return nil
	}

	level, err := logrus.ParseLevel(rootLogLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	return nil
}

func cmdVersion(_ *cobra.Command, _ []string) error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("cannot read build info")
	}
	var gitRev string
	for _, bs := range info.Settings {
		if bs.Key == "vcs.revision" {
			gitRev = bs.Value
			break
		}
	}
	if gitRev != "" {
		fmt.Printf("revision: %s\n", gitRev[:7])
	} else {
		fmt.Printf("revision: unknown\n")
	}
	return nil
}

func buildCobraCmdline() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:               "bp",
		Long:              "Read, write, and modify osbuild blueprint configuration files",
		PersistentPreRunE: rootPreRunE,
		SilenceErrors:     true,
	}

	rootCmd.PersistentFlags().StringVar(&rootLogLevel, "log-level", "", "logging level (debug, info, error); default error")

	versionCmd := &cobra.Command{
		Use:          "version",
		SilenceUsage: true,
		Hidden:       false,
		RunE:         cmdVersion,
	}
	rootCmd.AddCommand(versionCmd)

	setCmd := &cobra.Command{
		Use:          "set",
		SilenceUsage: true,
		Hidden:       false,
		Args:         cobra.ExactArgs(1),
		RunE:         cmdSet,
	}
	rootCmd.AddCommand(setCmd)

	return rootCmd, nil
}

func main() {
	cmd, err := buildCobraCmdline()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
