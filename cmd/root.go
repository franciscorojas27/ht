/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/franciscorojas27/HT/internal/ht"
	"github.com/franciscorojas27/HT/internal/ui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ht",
	Short: "HTTP command-line client.",
	Long:  "HT is a command-line HTTP client for quick testing, debugging and reproducible HTTP flows. Supports YAML files and a quick cURL-like mode.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFlags(rootFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return ht.Run(rootFlags)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var rootFlags ui.Flags

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	bindRootFlags(rootCmd)
}

func bindRootFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&rootFlags.YML, "yml", "", "Path to YAML configuration file")
	flags.BoolVar(&rootFlags.LS, "ls", false, "List requests defined in the YAML file")
	flags.StringVar(&rootFlags.Name, "name", "", "Request name to execute from the YAML file")
	flags.BoolVar(&rootFlags.DryRun, "dry-run", false, "Print the request without sending it")
	flags.BoolVar(&rootFlags.Verbose, "verbose", false, "Enable verbose output (timings)")
	flags.StringVar(&rootFlags.URL, "url", "", "Direct URL for quick request mode")
	flags.StringVarP(&rootFlags.Method, "method", "X", "", "HTTP method for quick request mode")
	flags.StringArrayVarP(&rootFlags.Headers, "header", "H", []string{}, "Extra header (repeatable)")
	flags.StringVarP(&rootFlags.Data, "data", "d", "", "Request body for quick request mode")
	flags.BoolVarP(&rootFlags.HeadOnly, "head", "I", false, "Fetch headers only (HEAD)")
	flags.BoolVarP(&rootFlags.Insecure, "insecure", "k", false, "Allow insecure TLS connections")
	flags.StringVarP(&rootFlags.User, "user", "u", "", "Basic auth in format user:pass")
}

func validateFlags(flags ui.Flags) error {
	quickMode := flags.YML == "" || flags.URL != "" || flags.Method != ""
	if quickMode {
		if flags.LS {
			return fmt.Errorf("-ls requires -yml")
		}
		if flags.URL == "" {
			return fmt.Errorf("missing --url for quick request")
		}
		if flags.Method != "" && flags.HeadOnly {
			return fmt.Errorf("--head cannot be combined with --method")
		}
		if flags.HeadOnly && flags.Data != "" {
			return fmt.Errorf("--head cannot be combined with --data")
		}
		if flags.User != "" {
			user, _, ok := strings.Cut(flags.User, ":")
			if !ok || user == "" {
				return fmt.Errorf("invalid --user value, expected user:pass")
			}
		}
		return nil
	}

	if flags.YML == "" {
		return fmt.Errorf("missing --yml")
	}
	if flags.LS {
		return nil
	}
	if flags.Name == "" {
		return fmt.Errorf("missing --name")
	}

	return nil
}
