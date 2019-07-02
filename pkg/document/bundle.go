package document

import (
	"fmt"
	"io"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/konfig"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resid"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/kustomize/api/types"

	utilyaml "opendev.org/airship/airshipctl/pkg/util/yaml"
)

// BundleFactory contains the objects within a bundle
type BundleFactory struct {
	krusty.Options
	resmap.ResMap
	filesys.FileSystem
}

// Bundle interface provides the specification for a bundle implementation
type Bundle interface {
	Write(out io.Writer) error
	GetKustomizeResourceMap() resmap.ResMap
	SetKustomizeResourceMap(resmap.ResMap) error
	GetKrustyOptions() krusty.Options
	SetKrustyOptions(krusty.Options) error
	SetFileSystem(filesys.FileSystem) error
	GetFileSystem() filesys.FileSystem
	Select(selector types.Selector) ([]Document, error)
	GetByGvk(string, string, string) ([]Document, error)
	GetByName(string) (Document, error)
	GetByAnnotation(string) ([]Document, error)
	GetByLabel(string) ([]Document, error)
	GetAllDocuments() ([]Document, error)
}

// NewBundle is a convenience function to create a new bundle
// Over time, it will evolve to support allowing more control
// for kustomize plugins
func NewBundle(fSys filesys.FileSystem, kustomizePath string, outputPath string) (Bundle, error) {

	var opts = &krusty.Options{
		RerorderTransformer: "none",
		LoadRestrictions:    types.LoadRestrictionsRootOnly,
		DoPrune:             false,
		PluginConfig:        konfig.DisabledPluginConfig(),
	}

	// init an empty bundle factory
	var bundle Bundle = &BundleFactory{}

	// set the fs and build options we will use
	bundle.SetFileSystem(fSys)
	bundle.SetKrustyOptions(*opts)

	// boiler plate to allow us to run Kustomize build
	// build a resource map of kustomize rendered objects
	k := krusty.MakeKustomizer(fSys, opts)
	m, err := k.Run(kustomizePath)
	bundle.SetKustomizeResourceMap(m)
	if err != nil {
		return bundle, err
	}

	return bundle, nil

}

// GetKustomizeResourceMap returns a Kustomize Resource Map for this bundle
func (b *BundleFactory) GetKustomizeResourceMap() resmap.ResMap {
	return b.ResMap
}

// SetKustomizeResourceMap allows us to set the populated resource map for this bundle.  In
// the future, it may modify it before saving it.
func (b *BundleFactory) SetKustomizeResourceMap(r resmap.ResMap) error {
	b.ResMap = r
	return nil
}

// GetKrustyOptions returns the build options object used to generate the resource map
// for this bundle
func (b *BundleFactory) GetKrustyOptions() krusty.Options {
	return b.Options
}

// SetKrustyOptions sets the build options to be used for this bundle. In
// the future, it may perform some basic validations.
func (b *BundleFactory) SetKrustyOptions(k krusty.Options) error {
	b.Options = k
	return nil
}

// SetFileSystem sets the filesystem that will be used by this bundle
func (b *BundleFactory) SetFileSystem(fSys filesys.FileSystem) error {
	b.FileSystem = fSys
	return nil
}

// GetFileSystem gets the filesystem that will be used by this bundle
func (b *BundleFactory) GetFileSystem() filesys.FileSystem {
	return b.FileSystem
}

// GetAllDocuments returns all documents in this bundle
func (b *BundleFactory) GetAllDocuments() ([]Document, error) {
	docSet := []Document{}
	for _, res := range b.ResMap.Resources() {
		// Construct Bundle document for each resource returned
		doc, err := NewDocument(res)
		if err != nil {
			return docSet, err
		}
		docSet = append(docSet, doc)
	}
	return docSet, nil
}

// GetByName finds a document by name, error if more than one document found
// or if no documents found
func (b *BundleFactory) GetByName(name string) (Document, error) {
	resSet := []*resource.Resource{}
	for _, res := range b.ResMap.Resources() {
		if res.GetName() == name {
			resSet = append(resSet, res)
		}
	}
	// alanmeadows(TODO): improve this and other error potentials by
	// by adding strongly typed errors
	switch found := len(resSet); {
	case found == 0:
		return &DocumentFactory{}, fmt.Errorf("No documents found with name %s", name)
	case found > 1:
		return &DocumentFactory{}, fmt.Errorf("More than one document found with name %s", name)
	default:
		return NewDocument(resSet[0])
	}
}

// Select offers a direct interface to pass a Kustomize Selector to the bundle
// returning Documents that match the criteria
func (b *BundleFactory) Select(selector types.Selector) ([]Document, error) {

	// use the kustomize select method
	resources, err := b.ResMap.Select(selector)
	if err != nil {
		return []Document{}, err
	}

	// Construct Bundle document for each resource returned
	docSet := []Document{}
	for _, res := range resources {
		doc, err := NewDocument(res)
		if err != nil {
			return docSet, err
		}
		docSet = append(docSet, doc)
	}
	return docSet, err
}

// GetByAnnotation is a convenience method to get documents for a particular annotation
func (b *BundleFactory) GetByAnnotation(annotation string) ([]Document, error) {

	// Construct kustomize annotation selector
	selector := types.Selector{AnnotationSelector: annotation}

	// pass it to the selector
	return b.Select(selector)

}

// GetByLabel is a convenience method to get documents for a particular label
func (b *BundleFactory) GetByLabel(label string) ([]Document, error) {

	// Construct kustomize annotation selector
	selector := types.Selector{LabelSelector: label}

	// pass it to the selector
	return b.Select(selector)

}

// GetByGvk is a convenience method to get documents for a particular Gvk tuple
func (b *BundleFactory) GetByGvk(group, version, kind string) ([]Document, error) {

	// Construct kustomize gvk object
	g := resid.Gvk{Group: group, Version: version, Kind: kind}

	// pass it to the selector
	selector := types.Selector{Gvk: g}
	return b.Select(selector)

}

// Write will write out the entire bundle resource map
func (b *BundleFactory) Write(out io.Writer) error {
	for _, res := range b.ResMap.Resources() {
		err := utilyaml.WriteOut(out, res)
		if err != nil {
			return err
		}
	}
	return nil
}
