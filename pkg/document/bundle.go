package document

import (
	"fmt"
	"io"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/api/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/api/k8sdeps/validator"
	fLdr "sigs.k8s.io/kustomize/api/loader"
	"sigs.k8s.io/kustomize/api/pgmconfig"
	pLdr "sigs.k8s.io/kustomize/api/plugins/loader"
	"sigs.k8s.io/kustomize/api/resid"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/kustomize/api/target"
	"sigs.k8s.io/kustomize/api/types"

	utilyaml "opendev.org/airship/airshipctl/pkg/util/yaml"
)

// KustomizeBuildOptions contain the options for running a Kustomize build on a bundle
type KustomizeBuildOptions struct {
	KustomizationPath string
	OutputPath        string
	LoadRestrictor    fLdr.LoadRestrictorFunc
	OutOrder          int
}

// BundleFactory contains the objects within a bundle
type BundleFactory struct {
	KustomizeBuildOptions
	resmap.ResMap
	filesys.FileSystem
}

// Bundle interface provides the specification for a bundle implementation
type Bundle interface {
	Write(out io.Writer) error
	GetKustomizeResourceMap() resmap.ResMap
	SetKustomizeResourceMap(resmap.ResMap) error
	GetKustomizeBuildOptions() KustomizeBuildOptions
	SetKustomizeBuildOptions(KustomizeBuildOptions) error
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

	var options = KustomizeBuildOptions{
		KustomizationPath: kustomizePath,
		OutputPath:        outputPath,
		LoadRestrictor:    fLdr.RestrictionRootOnly,
		OutOrder:          0,
	}

	// init an empty bundle factory
	var bundle Bundle = &BundleFactory{}

	// set the fs and build options we will use
	bundle.SetFileSystem(fSys)
	bundle.SetKustomizeBuildOptions(options)

	// boiler plate to allow us to run Kustomize build
	uf := kunstruct.NewKunstructuredFactoryImpl()
	pf := transformer.NewFactoryImpl()
	rf := resmap.NewFactory(resource.NewFactory(uf), pf)
	v := validator.NewKustValidator()

	pluginConfig, err := pgmconfig.EnabledPluginConfig()
	if err != nil {
		return bundle, err
	}
	pl := pLdr.NewLoader(pluginConfig, rf)

	ldr, err := fLdr.NewLoader(
		bundle.GetKustomizeBuildOptions().LoadRestrictor, bundle.GetKustomizeBuildOptions().KustomizationPath, fSys)
	if err != nil {
		return bundle, err
	}
	defer ldr.Cleanup()
	kt, err := target.NewKustTarget(ldr, v, rf, pf, pl)
	if err != nil {
		return bundle, err
	}

	// build a resource map of kustomize rendered objects
	m, err := kt.MakeCustomizedResMap()
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

// GetKustomizeBuildOptions returns the build options object used to generate the resource map
// for this bundle
func (b *BundleFactory) GetKustomizeBuildOptions() KustomizeBuildOptions {
	return b.KustomizeBuildOptions
}

// SetKustomizeBuildOptions sets the build options to be used for this bundle. In
// the future, it may perform some basic validations.
func (b *BundleFactory) SetKustomizeBuildOptions(k KustomizeBuildOptions) error {
	b.KustomizeBuildOptions = k
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
