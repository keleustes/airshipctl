package document

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/register"
)

// NewDocumentCommand creates a new command for managing airshipctl documents
func NewDocumentCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	// TODO(jeb): Find the appropriate place to do that without slowing down the cli
	register.RegisterCoreKinds()

	documentRootCmd := &cobra.Command{
		Use:   "document",
		Short: "manages deployment documents",
	}

	documentRootCmd.AddCommand(NewDocumentPullCommand(rootSettings))
	documentRootCmd.AddCommand(NewRenderCommand(rootSettings))

	return documentRootCmd
}
