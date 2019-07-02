package cmd

import (
	"io"
	"math/rand"
	"os"
	"time"

	argo "github.com/argoproj/argo/cmd/argo/commands"
	"github.com/spf13/cobra"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/cmd"
	kubectl "k8s.io/kubernetes/pkg/kubectl/cmd"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"opendev.org/airship/airshipctl/cmd/bootstrap"
	"opendev.org/airship/airshipctl/cmd/calico"
	"opendev.org/airship/airshipctl/cmd/ceph"
	"opendev.org/airship/airshipctl/cmd/cluster"
	"opendev.org/airship/airshipctl/cmd/completion"
	"opendev.org/airship/airshipctl/cmd/document"
	"opendev.org/airship/airshipctl/cmd/helm"
	"opendev.org/airship/airshipctl/cmd/openstack"
	"opendev.org/airship/airshipctl/cmd/prometheus"
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/log"
)

// NewAirshipCTLCommand creates a root `airshipctl` command with the default commands attached
func NewAirshipCTLCommand(out io.Writer) (*cobra.Command, *environment.AirshipCTLSettings, error) {
	rootCmd, settings, err := NewRootCmd(out)
	return AddDefaultAirshipCTLCommands(rootCmd, settings), settings, err
}

// NewRootCmd creates the root `airshipctl` command. All other commands are
// subcommands branching from this one
func NewRootCmd(out io.Writer) (*cobra.Command, *environment.AirshipCTLSettings, error) {
	rand.Seed(time.Now().UnixNano())

	settings := &environment.AirshipCTLSettings{}
	rootCmd := &cobra.Command{
		Use:           "airshipctl",
		Short:         "airshipctl is a unified entrypoint to various airship components",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Init(settings.Debug, cmd.OutOrStderr())
		},
		// BashCompletionFunction: bashCompletionFunc,
	}
	rootCmd.SetOutput(out)
	rootCmd.AddCommand(NewVersionCommand())

	settings.InitFlags(rootCmd)

	return rootCmd, settings, nil
}

// AddDefaultAirshipCTLCommands is a convenience function for adding all of the
// default commands to airshipctl
func AddDefaultAirshipCTLCommands(cmd *cobra.Command, settings *environment.AirshipCTLSettings) *cobra.Command {
	cmd.AddCommand(argo.NewCommand())
	cmd.AddCommand(bootstrap.NewBootstrapCommand(settings))
	cmd.AddCommand(cluster.NewClusterCommand(settings))
	cmd.AddCommand(completion.NewCompletionCommand())
	cmd.AddCommand(document.NewDocumentCommand(settings))
	cmd.AddCommand(kubectl.NewDefaultKubectlCommand())
	cmd.AddCommand(helm.NewHelmCommand(settings))
	cmd.AddCommand(openstack.NewOpenStackCommand(settings))
	cmd.AddCommand(prometheus.NewPrometheusCommand(settings))
	cmd.AddCommand(ceph.NewCephCommand(settings))
	cmd.AddCommand(calico.NewCalicoCommand(settings))

	// Should we use cmd.OutOrStdout?
	cmd.AddCommand(kubeadm.NewKubeadmCommand(os.Stdin, os.Stdout, os.Stderr))

	return cmd
}
