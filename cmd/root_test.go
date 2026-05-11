package cmd

import (
	"testing"

	"github.com/franciscorojas27/HT/internal/ui"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestValidateFlagsQuickMode(t *testing.T) {
	err := validateFlags(ui.Flags{})
	require.Error(t, err)

	err = validateFlags(ui.Flags{URL: "https://example.com"})
	require.NoError(t, err)

	err = validateFlags(ui.Flags{URL: "https://example.com", LS: true})
	require.Error(t, err)

	err = validateFlags(ui.Flags{URL: "https://example.com", HeadOnly: true, Method: "GET"})
	require.Error(t, err)

	err = validateFlags(ui.Flags{URL: "https://example.com", HeadOnly: true, Data: "a=1"})
	require.Error(t, err)
}

func TestValidateFlagsUserFormat(t *testing.T) {
	err := validateFlags(ui.Flags{URL: "https://example.com", User: "bad"})
	require.Error(t, err)

	err = validateFlags(ui.Flags{URL: "https://example.com", User: "user:pass"})
	require.NoError(t, err)
}

func TestValidateFlagsYAMLMode(t *testing.T) {
	err := validateFlags(ui.Flags{YML: "file.yml"})
	require.Error(t, err)

	err = validateFlags(ui.Flags{YML: "file.yml", Name: "get"})
	require.NoError(t, err)

	err = validateFlags(ui.Flags{YML: "file.yml", LS: true})
	require.NoError(t, err)
}

func TestBindRootFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "ht"}
	rootFlags = ui.Flags{}
	bindRootFlags(cmd)

	flags := cmd.Flags()
	require.NotNil(t, flags.Lookup("yml"))
	require.NotNil(t, flags.Lookup("ls"))
	require.NotNil(t, flags.Lookup("name"))
	require.NotNil(t, flags.Lookup("dry-run"))
	require.NotNil(t, flags.Lookup("verbose"))
	require.NotNil(t, flags.Lookup("url"))
	require.NotNil(t, flags.Lookup("method"))
	require.NotNil(t, flags.Lookup("header"))
	require.NotNil(t, flags.Lookup("data"))
	require.NotNil(t, flags.Lookup("head"))
	require.NotNil(t, flags.Lookup("insecure"))
	require.NotNil(t, flags.Lookup("user"))
}
