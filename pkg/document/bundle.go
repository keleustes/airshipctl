package document

import (
	// "errors"
	"fmt"
	"io"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/loader"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/kustomize/api/types"

	// "opendev.org/airship/airshipctl/pkg/log"
	// docplugins "opendev.org/airship/airshipctl/pkg/document/plugins"
	utilyaml "opendev.org/airship/airshipctl/pkg/util/yaml"
)

func init() {
	// NOTE (dukov) This is sort of a hack but it's the only way to add an
	// external 'builtin' plugin to Kustomize
	// builtinhelpers.TransformerFactories[builtinhelpers.Unknown] = docplugins.NewTransformerLoader
}

// KustomizeBuildOptions contain the options for running a Kustomize build on a bundle
type KustomizeBuildOptions struct {
	KustomizationPath string
	OutputPath        string
	LoadRestrictor    loader.LoadRestrictorFunc
	OutOrder          int
}

// BundleFactory contains the objects within a bundle
type BundleFactory struct {
	krusty.Options
	resmap.ResMap
	FileSystem
}

// Bundle interface provides the specification for a bundle implementation
type Bundle interface {
	Write(out io.Writer) error
	GetKustomizeResourceMap() resmap.ResMap
	SetKustomizeResourceMap(resmap.ResMap) error
	GetKrustyOptions() krusty.Options
	SetKrustyOptions(krusty.Options) error
	SetFileSystem(FileSystem) error
	GetFileSystem() FileSystem
	Select(selector Selector) ([]Document, error)
	SelectBundle(selector Selector) (Bundle, error)
	SelectByFieldValue(string, func(interface{}) bool) (Bundle, error)
	GetByGvk(string, string, string) ([]Document, error)
	GetByName(string) (Document, error)
	GetByAnnotation(annotationSelector string) ([]Document, error)
	GetByLabel(labelSelector string) ([]Document, error)
	GetAllDocuments() ([]Document, error)
}

// NewBundleByPath helper function that returns new document.Bundle interface based on clusterType and
// phase, example: helpers.NewBunde(airConfig, "ephemeral", "initinfra")
func NewBundleByPath(rootPath string) (Bundle, error) {
	return NewBundle(NewDocumentFs(), rootPath, "")
}

// NewBundle is a convenience function to create a new bundle
// Over time, it will evolve to support allowing more control
// for kustomize plugins
func NewBundle(fSys FileSystem, kustomizePath string, outputPath string) (Bundle, error) {
	var opts = &krusty.Options{
		RerorderTransformer: "none",
		LoadRestrictions:    types.LoadRestrictionsRootOnly,
		DoPrune:             false,
		// PluginConfig:        konfig.DisabledPluginConfig(),
	}

	// init an empty bundle factory
	bundle := &BundleFactory{}

	// set the fs and build options we will use
	bundle.SetFileSystem(fSys)
	bundle.SetKrustyOptions(*opts)

	// build a resource map of kustomize rendered objects
	k := krusty.MakeKustomizer(fSys, opts)
	m, err := k.Run(kustomizePath)
	bundle.SetKustomizeResourceMap(m)
	if err != nil {
		return bundle, err
	}
	err = bundle.SetKustomizeResourceMap(m)

	return bundle, err
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
func (b *BundleFactory) SetFileSystem(fSys FileSystem) error {
	b.FileSystem = fSys
	return nil
}

// GetFileSystem gets the filesystem that will be used by this bundle
func (b *BundleFactory) GetFileSystem() FileSystem {
	return b.FileSystem
}

// GetAllDocuments returns all documents in this bundle
func (b *BundleFactory) GetAllDocuments() ([]Document, error) {
	docSet := make([]Document, len(b.ResMap.Resources()))
	for i, res := range b.ResMap.Resources() {
		// Construct Bundle document for each resource returned
		doc, err := NewDocument(res)
		if err != nil {
			return docSet, err
		}
		docSet[i] = doc
	}
	return docSet, nil
}

// GetByName finds a document by name, error if more than one document found
// or if no documents found
func (b *BundleFactory) GetByName(name string) (Document, error) {
	resSet := make([]*resource.Resource, 0, len(b.ResMap.Resources()))
	for _, res := range b.ResMap.Resources() {
		if res.GetName() == name {
			resSet = append(resSet, res)
		}
	}
	// alanmeadows(TODO): improve this and other error potentials by
	// by adding strongly typed errors
	switch found := len(resSet); {
	case found == 0:
		return &Factory{}, fmt.Errorf("no documents found with name %s", name)
	case found > 1:
		return &Factory{}, fmt.Errorf("more than one document found with name %s", name)
	default:
		return NewDocument(resSet[0])
	}
}

// Select offers an interface to pass a Selector, built on top of kustomize Selector
// to the bundle returning Documents that match the criteria
func (b *BundleFactory) Select(selector Selector) ([]Document, error) {
	// use the kustomize select method
	resources, err := b.ResMap.Select(selector.Selector)
	if err != nil {
		return []Document{}, err
	}

	// Construct Bundle document for each resource returned
	docSet := make([]Document, len(resources))
	for i, res := range resources {
		var doc Document
		doc, err = NewDocument(res)
		if err != nil {
			return docSet, err
		}
		docSet[i] = doc
	}
	return docSet, err
}

// SelectBundle offers an interface to pass a Selector, built on top of kustomize Selector
// to the bundle returning a new Bundle that matches the criteria.  This is useful
// where you want to actually prune the underlying bundle you are working with
// rather then getting back the matching documents for scenarios like
// test cases where you want to pass in custom "filtered" bundles
// specific to the test case
func (b *BundleFactory) SelectBundle(selector Selector) (Bundle, error) {
	// use the kustomize select method
	resources, err := b.ResMap.Select(selector.Selector)
	if err != nil {
		return nil, err
	}

	// create a blank resourcemap and append the found resources
	// into the new resource map
	resourceMap := resmap.New()
	for _, res := range resources {
		if err = resourceMap.Append(res); err != nil {
			return nil, err
		}
	}

	// return a new bundle with the same options and filesystem
	// as this one but with a reduced resourceMap
	return &BundleFactory{
		Options:    b.GetKrustyOptions(),
		ResMap:     resourceMap,
		FileSystem: b.FileSystem,
	}, nil
}

// SelectByFieldValue returns new Bundle with filtered resource documents.
// Method iterates over all resources in the bundle. If resource has field
// (i.e. key) specified in JSON path, and the comparison function returns
// 'true' for value referenced by JSON path, then resource is added to
// resulting bundle.
// Example:
// The bundle contains 3 documents
//
//     ---
//     apiVersion: v1
//     kind: DocKind1
//     metadata:
//       name: doc1
//     spec:
//       somekey:
//         somefield: "someValue"
//     ---
//     apiVersion: v1
//     kind: DocKind2
//     metadata:
//       name: doc2
//     spec:
//       somekey:
//         somefield: "someValue"
//     ---
//     apiVersion: v1
//     kind: DocKind1
//     metadata:
//       name: doc3
//     spec:
//       somekey:
//         somefield: "someOtherValue"
//
// Execution of bundleInstance.SelectByFieldValue(
//		"spec.somekey.somefield",
//		func(v interface{}) { return v == "someValue" })
// will return a new Bundle instance containing 2 documents:
//     ---
//     apiVersion: v1
//     kind: DocKind1
//     metadata:
//       name: doc1
//     spec:
//       somekey:
//         somefield: "someValue"
//     ---
//     apiVersion: v1
//     kind: DocKind2
//     metadata:
//       name: doc2
//     spec:
//       somekey:
//         somefield: "someValue"
func (b *BundleFactory) SelectByFieldValue(path string, condition func(interface{}) bool) (Bundle, error) {
	result := &BundleFactory{
		Options:    b.GetKrustyOptions(),
		FileSystem: b.FileSystem,
	}
	resourceMap := resmap.New()
	for _, res := range b.Resources() {
		val, err := res.GetFieldValue(path)
		if err != nil {
			// if errors.As(err, &types.NoFieldError{}) {
			// this resource doesn't have the specified field - skip it
			// continue
			// } else {
			return nil, err
			// }
		}

		if condition(val) {
			if err = resourceMap.Append(res); err != nil {
				return nil, err
			}
		}
	}

	if err := result.SetKustomizeResourceMap(resourceMap); err != nil {
		return nil, err
	}
	return result, nil
}

// GetByAnnotation is a convenience method to get documents for a particular annotation
func (b *BundleFactory) GetByAnnotation(annotationSelector string) ([]Document, error) {
	// Construct kustomize annotation selector
	selector := NewSelector().ByAnnotation(annotationSelector)
	// pass it to the selector
	return b.Select(selector)
}

// GetByLabel is a convenience method to get documents for a particular label
func (b *BundleFactory) GetByLabel(labelSelector string) ([]Document, error) {
	// Construct kustomize label selector
	selector := NewSelector().ByLabel(labelSelector)
	// pass it to the selector
	return b.Select(selector)
}

// GetByGvk is a convenience method to get documents for a particular Gvk tuple
func (b *BundleFactory) GetByGvk(group, version, kind string) ([]Document, error) {
	// Construct kustomize gvk object
	selector := NewSelector().ByGvk(group, version, kind)

	// pass it to the selector
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
