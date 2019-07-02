package prometheus

import (
	"fmt"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/environment"
	pluginenv "opendev.org/airship/airshipctl/pkg/prometheus/environment"
)

// PluginSettingsID is used as a key in the root settings map of plugin settings
const PluginSettingsID = "prometheus"

// NewPrometheusCommand creates a new command for prometheus access
func NewPrometheusCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	pluginRootCmd := &cobra.Command{
		Use:   "prometheus",
		Short: "prometheus configuration and data viewer",
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
