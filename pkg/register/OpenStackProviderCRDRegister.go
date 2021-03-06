// Code generated by pluginator on OpenStackProviderCRDRegister; DO NOT EDIT.
package register

import (
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha2"
	"sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha3"
	"sigs.k8s.io/kustomize/api/resmap"
)

// plugin loads the OpenStackProviderChart CRD scheme into kustomize
type OpenStackProviderCRDRegisterPlugin struct {
}

//nolint: golint

func (p *OpenStackProviderCRDRegisterPlugin) Config(
	_ *resmap.PluginHelpers, _ []byte) (err error) {

	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	err = v1alpha2.SchemeBuilder.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}
	return v1alpha3.SchemeBuilder.AddToScheme(scheme.Scheme)
}

func (p *OpenStackProviderCRDRegisterPlugin) Transform(m resmap.ResMap) error {
	return nil
}

func NewOpenStackProviderCRDRegisterPlugin() resmap.TransformerPlugin {
	return &OpenStackProviderCRDRegisterPlugin{}
}
