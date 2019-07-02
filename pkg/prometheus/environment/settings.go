package environment

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/environment"
)

// Settings is a container for all of the settings needed by workflows
type Settings struct {
	*environment.AirshipCTLSettings

	// Initialized denotes whether the settings have been initialized or not. It is useful for unit-testing
	Initialized bool
}

// InitFlags adds the default settings flags to cmd
func (s *Settings) InitFlags(cmd *cobra.Command) {
	// flags := cmd.PersistentFlags()
	// flags.StringVar(&s.KubeConfigFilePath, "kubeconfig", "", "path to kubeconfig")
	s.Initialized = true
}
