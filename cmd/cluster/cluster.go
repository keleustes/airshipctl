package cluster

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	pluginenv "opendev.org/airship/airshipctl/pkg/clusterctl/environment"
	"opendev.org/airship/airshipctl/pkg/environment"
)

var (
	// ClusterUse subcommand string
	ClusterUse = "cluster"

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a cluster API resource",
		Long:  `Create a cluster API resource with one command`,
	}
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a cluster API resource",
		Long:  `Delete a cluster API resource with one command`,
	}
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate an API resource created by cluster API.",
		Long:  `Validate an API resource created by cluster API. See subcommands for supported API resources.`,
	}
	alphaCmd = &cobra.Command{
		Use:   "alpha",
		Short: "Alpha/Experimental features",
		Long:  `Alpha/Experimental features`,
	}
)

// NewClusterCommand returns cobra command object of the airshipctl cluster and adds it's subcommands.
func NewClusterCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	clusterRootCmd := &cobra.Command{
		Use: ClusterUse,
		// TODO: (kkalynovskyi) Add more description when more subcommands are added
		Short: "Control Kubernetes cluster",
		Long:  "Interactions with Kubernetes cluster, such as get status, deploy initial infrastructure",
	}

	clusterRootCmd.AddCommand(NewCmdInitInfra(rootSettings))

	// Standard airshipctl plugin setup
	pluginSettings := &pluginenv.Settings{
		AirshipCTLSettings: rootSettings,
	}
	pluginSettings.InitFlags(clusterRootCmd)

	clusterRootCmd.AddCommand(createCmd)
	clusterRootCmd.AddCommand(deleteCmd)
	clusterRootCmd.AddCommand(validateCmd)
	clusterRootCmd.AddCommand(alphaCmd)
	return clusterRootCmd
}

func exitWithHelp(cmd *cobra.Command, err string) {
	fmt.Fprintln(os.Stderr, err)
	cmd.Help()
	os.Exit(1)
}
