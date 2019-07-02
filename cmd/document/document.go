package document

import (
	"flag"
	"os"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/cmd/document/secret"
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/register"
)

// NewDocumentCommand creates a new command for managing airshipctl documents
func NewDocumentCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	// TODO(jeb): Find the appropriate place to do that without slowing down the cli
	register.RegisterCoreKinds()

	stdOut := os.Stdout

	documentRootCmd := &cobra.Command{
		Use:   "document",
		Short: "manages deployment documents",
		Long: `
Manages declarative configuration of Kubernetes.
See https://sigs.k8s.io/kustomize
`,
	}

	// Add the document build command
	documentRootCmd.AddCommand(NewCmdBuild(stdOut))

	// Add the document secret command
	documentRootCmd.AddCommand(secret.NewSecretCommand(rootSettings))

	// Add the document generate command
	documentRootCmd.AddCommand(NewGenerateCommand())

	documentRootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})
	return documentRootCmd
}
