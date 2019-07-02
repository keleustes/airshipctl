package openstack

import (
	"fmt"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/environment"
	pluginenv "opendev.org/airship/airshipctl/pkg/openstack/environment"
)

// PluginSettingsID is used as a key in the root settings map of plugin settings
const PluginSettingsID = "openstack"

// NewOpenStackCommand creates a new command for openstack access
func NewOpenStackCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	pluginRootCmd := &cobra.Command{
		Use:   "openstack",
		Short: "openstack configuration and data viewer",
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
