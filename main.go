package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "v3.29.0"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "forge",
		Short: "NHN forge build tool",
		Long:  "A Go-based build tool that manages development tools and generates Makefiles",
	}

	var syncCmd = &cobra.Command{
		Use:   "sync [project-name] [source-file] [output-file]",
		Short: "Sync project configuration and generate Makefile",
		Args:  cobra.ExactArgs(3),
		Run:   runSync,
	}

	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
