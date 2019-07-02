package ceph

import (
	"fmt"

	"github.com/spf13/cobra"

	pluginenv "opendev.org/airship/airshipctl/pkg/ceph/environment"
	"opendev.org/airship/airshipctl/pkg/environment"
)

// PluginSettingsID is used as a key in the root settings map of plugin settings
const PluginSettingsID = "ceph"

// NewCephCommand creates a new command for ceph tools
func NewCephCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	pluginRootCmd := &cobra.Command{
		Use:   "ceph",
		Short: "ceph configuration and data viewer",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Under construction\n")
		},
	}

	// Standard airshipctl plugin setup
	pluginSettings := &pluginenv.Settings{
		AirshipCTLSettings: rootSettings,
	}
	pluginSettings.InitFlags(pluginRootCmd)

	return pluginRootCmd
}
