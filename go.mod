module opendev.org/airship/airshipctl

go 1.13

require (
	github.com/Masterminds/semver v1.5.0
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/Nordix/go-redfish v0.0.0-20191016124000-fd2ad07270c9
	github.com/Nordix/go-redfish/client v0.0.0-20191016124000-fd2ad07270c9
	github.com/argoproj/argo v0.0.0-00010101000000-000000000000
	github.com/argoproj/argo-rollouts v0.5.0
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/docker v1.13.1
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/ghodss/yaml v1.0.1-0.20180820084758-c7ce16629ff4
	github.com/golang/protobuf v1.3.2
	github.com/gosuri/uitable v0.0.3
	github.com/keleustes/armada-crd v1.16.2-keleustes.20191102
	github.com/keleustes/capi-yaml-gen v1.16.2-keleustes.20191102
	github.com/openshift/api v0.0.0-20190322043348-8741ff068a47
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0
	google.golang.org/grpc v1.24.0
	istio.io/istio v0.0.0-20191009042236-b4e86c385016
	k8s.io/api v0.0.0
	k8s.io/apiextensions-apiserver v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/cli-runtime v0.0.0
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/helm v2.14.3+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.0.0
	k8s.io/kubernetes v1.16.2
	k8s.io/utils v0.0.0-20190923111123-69764acb6e8e
	sigs.k8s.io/cluster-api v0.3.0
	sigs.k8s.io/cluster-api-provider-baremetal v0.0.0-00010101000000-000000000000
	sigs.k8s.io/cluster-api-provider-openstack v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.3.0
	sigs.k8s.io/kustomize/api v1.16.2
	sigs.k8s.io/kustomize/v3 v3.1.1-0.20190821175718-4b67a6de1296
	sigs.k8s.io/yaml v1.1.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.0.0+incompatible
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest v13.0.1-0.20190815170712-85d9c035382e+incompatible
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest v13.0.1-0.20190816223009-243e2a3d5e75+incompatible
	github.com/Azure/go-autorest/autorest/date => github.com/Azure/go-autorest v13.0.1-0.20190816215130-5bd9621f41a0+incompatible
	github.com/Azure/go-autorest/autorest/mocks => github.com/Azure/go-autorest v13.0.1-0.20190816215130-5bd9621f41a0+incompatible

	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.1

	github.com/argoproj/argo => github.com/keleustes/argo v1.16.2-keleustes.20191102
	github.com/colinmarc/hdfs => github.com/colinmarc/hdfs v0.0.0-20180802165501-48eb8d6c34a9
	// Pin docker to some obscure version
	github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/ghodss/yaml => github.com/ghodss/yaml v0.0.0-20180820084758-c7ce16629ff4

	k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191016113550-5357c4baaf65
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191016112112-5190913f932d
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191016115326-20453efc2458
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191016115129-c07a134afb42
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191004115455-8e001e5d1894
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191016111319-039242c015a9
	k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190828162817-608eb1dad4ac
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191016115521-756ffa5af0bd
	k8s.io/helm => github.com/keleustes/helm v1.16.2-keleustes.20191102
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191016112429-9587704a8ad4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191016114939-2b2b218dc1df
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191016114407-2e83b6f20229
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191016114748-65049c67a58b
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191016114556-7841ed97f1b2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191016115753-cf0698c3a16b
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191016113814-3b1a734dba6e
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191016112829-06bb3c9d77c9

	sigs.k8s.io/cluster-api => github.com/keleustes/cluster-api v1.16.2-keleustes.20191102
	sigs.k8s.io/cluster-api-provider-aws => github.com/keleustes/cluster-api-provider-aws v1.16.2-keleustes.20191102
	sigs.k8s.io/cluster-api-provider-baremetal => github.com/keleustes/cluster-api-provider-baremetal v1.16.2-keleustes.20191102
	sigs.k8s.io/cluster-api-provider-openstack => github.com/keleustes/cluster-api-provider-openstack v1.16.2-keleustes.20191102
	sigs.k8s.io/controller-runtime => github.com/keleustes/controller-runtime v1.16.2-keleustes.20191102
	sigs.k8s.io/kustomize/api => github.com/keleustes/kustomize/api v1.16.2-keleustes.20191117
)
