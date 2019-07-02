package calico

import (
	"fmt"

	"github.com/spf13/cobra"

	pluginenv "opendev.org/airship/airshipctl/pkg/calico/environment"
	"opendev.org/airship/airshipctl/pkg/environment"
)

// PluginSettingsID is used as a key in the root settings map of plugin settings
const PluginSettingsID = "calico"

// NewCalicoCommand creates a new command for calico tools
func NewCalicoCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	pluginRootCmd := &cobra.Command{
		Use:   "calico",
		Short: "calico configuration and data viewer",
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
