/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	clusterclient "opendev.org/airship/airshipctl/pkg/client/cluster"
	"opendev.org/airship/airshipctl/pkg/phases"
	"sigs.k8s.io/cluster-api/util/yaml"
)

type AlphaPhaseApplyClusterOptions struct {
	Kubeconfig string
	Cluster    string
}

var paco = &AlphaPhaseApplyClusterOptions{}

var alphaPhaseApplyClusterCmd = &cobra.Command{
	Use:   "apply-cluster",
	Short: "Apply Cluster",
	Long:  `Apply Cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if paco.Cluster == "" {
			exitWithHelp(cmd, "Please provide yaml file for cluster definition.")
		}

		if paco.Kubeconfig == "" {
			exitWithHelp(cmd, "Please provide a kubeconfig file.")
		}

		if err := RunAlphaPhaseApplyCluster(paco); err != nil {
			klog.Exit(err)
		}
	},
}

func RunAlphaPhaseApplyCluster(paco *AlphaPhaseApplyClusterOptions) error {
	kubeconfig, err := ioutil.ReadFile(paco.Kubeconfig)
	if err != nil {
		return err
	}

	out, err := yaml.Parse(yaml.ParseInput{File: paco.Cluster})
	if err != nil {
		return err
	}
	if len(out.Clusters) == 0 {
		return errors.Errorf("no Cluster object found in file %q", paco.Cluster)
	}

	clientFactory := clusterclient.NewFactory()
	client, err := clientFactory.NewClientFromKubeconfig(string(kubeconfig))
	if err != nil {
		return errors.Wrap(err, "unable to create cluster client")
	}

	cluster := out.Clusters[0]
	if err := phases.ApplyCluster(
		client,
		cluster,
		yaml.ExtractClusterReferences(out, cluster)...,
	); err != nil {
		return errors.Wrap(err, "unable to apply cluster")
	}

	return nil
}

func init() {
	// Required flags
	alphaPhaseApplyClusterCmd.Flags().StringVarP(&paco.Kubeconfig, "kubeconfig", "", "", "Path for the kubeconfig file to use")
	alphaPhaseApplyClusterCmd.Flags().StringVarP(&paco.Cluster, "cluster", "c", "", "A yaml file containing cluster object definition")
	alphaPhasesCmd.AddCommand(alphaPhaseApplyClusterCmd)
}
