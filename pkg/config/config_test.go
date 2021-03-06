/*
Copyright 2014 The Kubernetes Authors.

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

package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/testutil"
)

const (
	stringDelta        = "_changed"
	currentContextName = "def_ephemeral"
	defaultString      = "default"
)

func TestString(t *testing.T) {
	fSys := testutil.SetupTestFs(t, "testdata")

	tests := []struct {
		name     string
		stringer fmt.Stringer
	}{
		{
			name:     "config",
			stringer: testutil.DummyConfig(),
		},
		{
			name:     "context",
			stringer: testutil.DummyContext(),
		},
		{
			name:     "cluster",
			stringer: testutil.DummyCluster(),
		},
		{
			name:     "authinfo",
			stringer: testutil.DummyAuthInfo(),
		},
		{
			name:     "manifest",
			stringer: testutil.DummyManifest(),
		},
		{
			name:     "modules",
			stringer: testutil.DummyModules(),
		},
		{
			name:     "repository",
			stringer: testutil.DummyRepository(),
		},
		{
			name:     "repo-auth",
			stringer: testutil.DummyRepoAuth(),
		},
		{
			name:     "repo-checkout",
			stringer: testutil.DummyRepoCheckout(),
		},
		{
			name:     "bootstrap",
			stringer: testutil.DummyBootstrap(),
		},
		{
			name:     "bootstrap",
			stringer: testutil.DummyBootstrap(),
		},
		{
			name: "builder",
			stringer: &config.Builder{
				UserDataFileName:       "user-data",
				NetworkConfigFileName:  "netconfig",
				OutputMetadataFileName: "output-metadata.yaml",
			},
		},
		{
			name: "container",
			stringer: &config.Container{
				Volume:           "/dummy:dummy",
				Image:            "dummy_image:dummy_tag",
				ContainerRuntime: "docker",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			filename := fmt.Sprintf("/%s-string.yaml", tt.name)
			data, err := fSys.ReadFile(filename)
			require.NoError(t, err)

			assert.Equal(t, string(data), tt.stringer.String())
		})
	}
}

func TestPrettyString(t *testing.T) {
	fSys := testutil.SetupTestFs(t, "testdata")
	data, err := fSys.ReadFile("/prettycluster-string.yaml")
	require.NoError(t, err)

	cluster := testutil.DummyCluster()
	assert.EqualValues(t, cluster.PrettyString(), string(data))
}

func TestEqual(t *testing.T) {
	t.Run("config-equal", func(t *testing.T) {
		testConfig1 := config.NewConfig()
		testConfig2 := config.NewConfig()
		testConfig2.Kind = "Different"
		assert.True(t, testConfig1.Equal(testConfig1))
		assert.False(t, testConfig1.Equal(testConfig2))
		assert.False(t, testConfig1.Equal(nil))
	})

	t.Run("cluster-equal", func(t *testing.T) {
		testCluster1 := &config.Cluster{NameInKubeconf: "same"}
		testCluster2 := &config.Cluster{NameInKubeconf: "different"}
		assert.True(t, testCluster1.Equal(testCluster1))
		assert.False(t, testCluster1.Equal(testCluster2))
		assert.False(t, testCluster1.Equal(nil))
	})

	t.Run("context-equal", func(t *testing.T) {
		testContext1 := &config.Context{NameInKubeconf: "same"}
		testContext2 := &config.Context{NameInKubeconf: "different"}
		assert.True(t, testContext1.Equal(testContext1))
		assert.False(t, testContext1.Equal(testContext2))
		assert.False(t, testContext1.Equal(nil))
	})

	// TODO(howell): this needs to be fleshed out when the AuthInfo type is finished
	t.Run("authinfo-equal", func(t *testing.T) {
		testAuthInfo1 := &config.AuthInfo{}
		assert.True(t, testAuthInfo1.Equal(testAuthInfo1))
		assert.False(t, testAuthInfo1.Equal(nil))
	})

	t.Run("manifest-equal", func(t *testing.T) {
		testManifest1 := &config.Manifest{TargetPath: "same"}
		testManifest2 := &config.Manifest{TargetPath: "different"}
		assert.True(t, testManifest1.Equal(testManifest1))
		assert.False(t, testManifest1.Equal(testManifest2))
		assert.False(t, testManifest1.Equal(nil))
	})

	t.Run("repository-equal", func(t *testing.T) {
		testRepository1 := &config.Repository{URLString: "same"}
		testRepository2 := &config.Repository{URLString: "different"}
		assert.True(t, testRepository1.Equal(testRepository1))
		assert.False(t, testRepository1.Equal(testRepository2))
		assert.False(t, testRepository1.Equal(nil))
	})
	t.Run("auth-equal", func(t *testing.T) {
		testSpec1 := &config.RepoAuth{}
		testSpec2 := &config.RepoAuth{}
		testSpec2.Type = "ssh-key"
		assert.True(t, testSpec1.Equal(testSpec1))
		assert.False(t, testSpec1.Equal(testSpec2))
		assert.False(t, testSpec1.Equal(nil))
	})
	t.Run("checkout-equal", func(t *testing.T) {
		testSpec1 := &config.RepoCheckout{}
		testSpec2 := &config.RepoCheckout{}
		testSpec2.Branch = "Master"
		assert.True(t, testSpec1.Equal(testSpec1))
		assert.False(t, testSpec1.Equal(testSpec2))
		assert.False(t, testSpec1.Equal(nil))
	})

	t.Run("modules-equal", func(t *testing.T) {
		testModules1 := config.NewModules()
		testModules2 := config.NewModules()
		testModules2.BootstrapInfo["different"] = &config.Bootstrap{
			Container: &config.Container{Volume: "different"},
		}
		assert.True(t, testModules1.Equal(testModules1))
		assert.False(t, testModules1.Equal(testModules2))
		assert.False(t, testModules1.Equal(nil))
	})

	t.Run("bootstrap-equal", func(t *testing.T) {
		testBootstrap1 := &config.Bootstrap{
			Container: &config.Container{
				Image: "same",
			},
		}
		testBootstrap2 := &config.Bootstrap{
			Container: &config.Container{
				Image: "different",
			},
		}
		assert.True(t, testBootstrap1.Equal(testBootstrap1))
		assert.False(t, testBootstrap1.Equal(testBootstrap2))
		assert.False(t, testBootstrap1.Equal(nil))
	})

	t.Run("container-equal", func(t *testing.T) {
		testContainer1 := &config.Container{Image: "same"}
		testContainer2 := &config.Container{Image: "different"}
		assert.True(t, testContainer1.Equal(testContainer1))
		assert.False(t, testContainer1.Equal(testContainer2))
		assert.False(t, testContainer1.Equal(nil))
	})

	t.Run("builder-equal", func(t *testing.T) {
		testBuilder1 := &config.Builder{UserDataFileName: "same"}
		testBuilder2 := &config.Builder{UserDataFileName: "different"}
		assert.True(t, testBuilder1.Equal(testBuilder1))
		assert.False(t, testBuilder1.Equal(testBuilder2))
		assert.False(t, testBuilder1.Equal(nil))
	})
}

func TestLoadConfig(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	assert.Len(t, conf.Clusters, 5)
	require.Contains(t, conf.Clusters, "def")
	assert.Len(t, conf.Clusters["def"].ClusterTypes, 2)
	assert.Len(t, conf.Contexts, 3)
	assert.Len(t, conf.AuthInfos, 3)
}

func TestPersistConfig(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	err := conf.PersistConfig()
	require.NoError(t, err)

	// Check that the files were created
	assert.FileExists(t, conf.LoadedConfigPath())
	assert.FileExists(t, conf.KubeConfigPath())
	// Check that the invalid name was changed to a valid one
	assert.Contains(t, conf.KubeConfig().Clusters, "invalidName_target")

	// Check that the missing cluster was added to the airshipconfig
	assert.Contains(t, conf.Clusters, "onlyinkubeconf")

	// Check that the "stragglers" were removed from the airshipconfig
	assert.NotContains(t, conf.Clusters, "straggler")
}

func TestEnsureComplete(t *testing.T) {
	// This test is intentionally verbose. Since a user of EnsureComplete
	// does not need to know about the order of validation, each test
	// object passed into EnsureComplete should have exactly one issue, and
	// be otherwise valid
	tests := []struct {
		name        string
		config      config.Config
		expectedErr error
	}{
		{
			name: "no clusters defined",
			config: config.Config{
				Clusters:       map[string]*config.ClusterPurpose{},
				AuthInfos:      map[string]*config.AuthInfo{"testAuthInfo": {}},
				Contexts:       map[string]*config.Context{"testContext": {Manifest: "testManifest"}},
				Manifests:      map[string]*config.Manifest{"testManifest": {}},
				CurrentContext: "testContext",
			},
			expectedErr: config.ErrMissingConfig{What: "At least one cluster needs to be defined"},
		},
		{
			name: "no users defined",
			config: config.Config{
				Clusters:       map[string]*config.ClusterPurpose{"testCluster": {}},
				AuthInfos:      map[string]*config.AuthInfo{},
				Contexts:       map[string]*config.Context{"testContext": {Manifest: "testManifest"}},
				Manifests:      map[string]*config.Manifest{"testManifest": {}},
				CurrentContext: "testContext",
			},
			expectedErr: config.ErrMissingConfig{What: "At least one Authentication Information (User) needs to be defined"},
		},
		{
			name: "no contexts defined",
			config: config.Config{
				Clusters:       map[string]*config.ClusterPurpose{"testCluster": {}},
				AuthInfos:      map[string]*config.AuthInfo{"testAuthInfo": {}},
				Contexts:       map[string]*config.Context{},
				Manifests:      map[string]*config.Manifest{"testManifest": {}},
				CurrentContext: "testContext",
			},
			expectedErr: config.ErrMissingConfig{What: "At least one Context needs to be defined"},
		},
		{
			name: "no manifests defined",
			config: config.Config{
				Clusters:       map[string]*config.ClusterPurpose{"testCluster": {}},
				AuthInfos:      map[string]*config.AuthInfo{"testAuthInfo": {}},
				Contexts:       map[string]*config.Context{"testContext": {Manifest: "testManifest"}},
				Manifests:      map[string]*config.Manifest{},
				CurrentContext: "testContext",
			},
			expectedErr: config.ErrMissingConfig{What: "At least one Manifest needs to be defined"},
		},
		{
			name: "current context not defined",
			config: config.Config{
				Clusters:       map[string]*config.ClusterPurpose{"testCluster": {}},
				AuthInfos:      map[string]*config.AuthInfo{"testAuthInfo": {}},
				Contexts:       map[string]*config.Context{"testContext": {Manifest: "testManifest"}},
				Manifests:      map[string]*config.Manifest{"testManifest": {}},
				CurrentContext: "",
			},
			expectedErr: config.ErrMissingConfig{What: "Current Context is not defined"},
		},
		{
			name: "no context for current context",
			config: config.Config{
				Clusters:       map[string]*config.ClusterPurpose{"testCluster": {}},
				AuthInfos:      map[string]*config.AuthInfo{"testAuthInfo": {}},
				Contexts:       map[string]*config.Context{"DIFFERENT_CONTEXT": {Manifest: "testManifest"}},
				Manifests:      map[string]*config.Manifest{"testManifest": {}},
				CurrentContext: "testContext",
			},
			expectedErr: config.ErrMissingConfig{What: "Current Context (testContext) does not identify a defined Context"},
		},
		{
			name: "no manifest for current context",
			config: config.Config{
				Clusters:       map[string]*config.ClusterPurpose{"testCluster": {}},
				AuthInfos:      map[string]*config.AuthInfo{"testAuthInfo": {}},
				Contexts:       map[string]*config.Context{"testContext": {Manifest: "testManifest"}},
				Manifests:      map[string]*config.Manifest{"DIFFERENT_MANIFEST": {}},
				CurrentContext: "testContext",
			},
			expectedErr: config.ErrMissingConfig{What: "Current Context (testContext) does not identify a defined Manifest"},
		},
		{
			name: "complete config",
			config: config.Config{
				Clusters:       map[string]*config.ClusterPurpose{"testCluster": {}},
				AuthInfos:      map[string]*config.AuthInfo{"testAuthInfo": {}},
				Contexts:       map[string]*config.Context{"testContext": {Manifest: "testManifest"}},
				Manifests:      map[string]*config.Manifest{"testManifest": {}},
				CurrentContext: "testContext",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(subTest *testing.T) {
			actualErr := tt.config.EnsureComplete()
			assert.Equal(subTest, tt.expectedErr, actualErr)
		})
	}
}

func TestCurrentContextBootstrapInfo(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	clusterName := "def"
	clusterType := "ephemeral"

	bootstrapInfo, err := conf.CurrentContextBootstrapInfo()
	require.Error(t, err)
	assert.Nil(t, bootstrapInfo)

	conf.CurrentContext = currentContextName
	conf.Clusters[clusterName].ClusterTypes[clusterType].Bootstrap = defaultString
	conf.Contexts[currentContextName].Manifest = defaultString
	conf.Contexts[currentContextName].KubeContext().Cluster = clusterName

	bootstrapInfo, err = conf.CurrentContextBootstrapInfo()
	require.NoError(t, err)
	assert.Equal(t, conf.ModulesConfig.BootstrapInfo[defaultString], bootstrapInfo)
}

func TestPurge(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	// Store it
	err := conf.PersistConfig()
	assert.NoErrorf(t, err, "Unable to persist configuration expected at %v", conf.LoadedConfigPath())

	// Verify that the file is there
	_, err = os.Stat(conf.LoadedConfigPath())
	assert.Falsef(t, os.IsNotExist(err), "Test config was not persisted at %v, cannot validate Purge",
		conf.LoadedConfigPath())

	// Delete it
	err = conf.Purge()
	assert.NoErrorf(t, err, "Unable to Purge file at %v", conf.LoadedConfigPath())

	// Verify its gone
	_, err = os.Stat(conf.LoadedConfigPath())
	assert.Falsef(t, os.IsExist(err), "Purge failed to remove file at %v", conf.LoadedConfigPath())
}

func TestKClusterString(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	kClusters := conf.KubeConfig().Clusters
	for kClust := range kClusters {
		assert.NotEmpty(t, config.KClusterString(kClusters[kClust]))
	}
	assert.EqualValues(t, config.KClusterString(nil), "null\n")
}
func TestKContextString(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	kContexts := conf.KubeConfig().Contexts
	for kCtx := range kContexts {
		assert.NotEmpty(t, config.KContextString(kContexts[kCtx]))
	}
	assert.EqualValues(t, config.KClusterString(nil), "null\n")
}
func TestKAuthInfoString(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	kAuthInfos := conf.KubeConfig().AuthInfos
	for kAi := range kAuthInfos {
		assert.NotEmpty(t, config.KAuthInfoString(kAuthInfos[kAi]))
	}
	assert.EqualValues(t, config.KAuthInfoString(nil), "null\n")
}

func TestComplexName(t *testing.T) {
	cName := "aCluster"
	ctName := config.Ephemeral
	clusterName := config.NewClusterComplexName()
	clusterName.WithType(cName, ctName)
	assert.EqualValues(t, cName+"_"+ctName, clusterName.Name())
	assert.EqualValues(t, cName, clusterName.ClusterName())
	assert.EqualValues(t, ctName, clusterName.ClusterType())

	cName = "bCluster"
	clusterName.SetClusterName(cName)
	clusterName.SetDefaultType()
	ctName = clusterName.ClusterType()
	assert.EqualValues(t, cName+"_"+ctName, clusterName.Name())
	assert.EqualValues(t, "clusterName:"+cName+", clusterType:"+ctName, clusterName.String())
}

func TestValidClusterTypeFail(t *testing.T) {
	err := config.ValidClusterType("Fake")
	assert.Error(t, err)
}

func TestSetLoadedConfigPath(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	testPath := "/tmp/loadedconfig"

	assert.NotEqual(t, testPath, conf.LoadedConfigPath())
	conf.SetLoadedConfigPath(testPath)
	assert.Equal(t, testPath, conf.LoadedConfigPath())
}

func TestSetKubeConfigPath(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	testPath := "/tmp/kubeconfig"

	assert.NotEqual(t, testPath, conf.KubeConfigPath())
	conf.SetKubeConfigPath(testPath)
	assert.Equal(t, testPath, conf.KubeConfigPath())
}

func TestGetCluster(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	cluster, err := conf.GetCluster("def", config.Ephemeral)
	require.NoError(t, err)

	// Test Positives
	assert.EqualValues(t, cluster.NameInKubeconf, "def_ephemeral")
	assert.EqualValues(t, cluster.KubeCluster().Server, "http://5.6.7.8")

	// Test Wrong Cluster
	_, err = conf.GetCluster("unknown", config.Ephemeral)
	assert.Error(t, err)

	// Test Wrong Cluster Type
	_, err = conf.GetCluster("def", "Unknown")
	assert.Error(t, err)
}

func TestAddCluster(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	co := testutil.DummyClusterOptions()
	cluster, err := conf.AddCluster(co)
	require.NoError(t, err)

	assert.EqualValues(t, conf.Clusters[co.Name].ClusterTypes[co.ClusterType], cluster)
}

func TestModifyCluster(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	co := testutil.DummyClusterOptions()
	cluster, err := conf.AddCluster(co)
	require.NoError(t, err)

	co.Server += stringDelta
	co.InsecureSkipTLSVerify = true
	co.EmbedCAData = true
	mcluster, err := conf.ModifyCluster(cluster, co)
	require.NoError(t, err)
	assert.EqualValues(t, conf.Clusters[co.Name].ClusterTypes[co.ClusterType].KubeCluster().Server, co.Server)
	assert.EqualValues(t, conf.Clusters[co.Name].ClusterTypes[co.ClusterType], mcluster)

	// Error case
	co.CertificateAuthority = "unknown"
	_, err = conf.ModifyCluster(cluster, co)
	assert.Error(t, err)
}

func TestGetClusters(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	clusters := conf.GetClusters()
	assert.Len(t, clusters, 5)
}

func TestGetContexts(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	contexts := conf.GetContexts()
	assert.Len(t, contexts, 3)
}

func TestGetContext(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	context, err := conf.GetContext("def_ephemeral")
	require.NoError(t, err)

	// Test Positives
	assert.EqualValues(t, context.NameInKubeconf, "def_ephemeral")
	assert.EqualValues(t, context.KubeContext().Cluster, "def_ephemeral")

	// Test Wrong Cluster
	_, err = conf.GetContext("unknown")
	assert.Error(t, err)
}

func TestAddContext(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	co := testutil.DummyContextOptions()
	context := conf.AddContext(co)
	assert.EqualValues(t, conf.Contexts[co.Name], context)
}

func TestModifyContext(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	co := testutil.DummyContextOptions()
	context := conf.AddContext(co)

	co.Namespace += stringDelta
	co.Cluster += stringDelta
	co.AuthInfo += stringDelta
	co.Manifest += stringDelta
	conf.ModifyContext(context, co)
	assert.EqualValues(t, conf.Contexts[co.Name].KubeContext().Namespace, co.Namespace)
	assert.EqualValues(t, conf.Contexts[co.Name].KubeContext().Cluster, co.Cluster)
	assert.EqualValues(t, conf.Contexts[co.Name].KubeContext().AuthInfo, co.AuthInfo)
	assert.EqualValues(t, conf.Contexts[co.Name].Manifest, co.Manifest)
	assert.EqualValues(t, conf.Contexts[co.Name], context)
}

func TestGetCurrentContext(t *testing.T) {
	t.Run("getCurrentContext", func(t *testing.T) {
		conf, cleanup := testutil.InitConfig(t)
		defer cleanup(t)

		context, err := conf.GetCurrentContext()
		require.Error(t, err)
		assert.Nil(t, context)

		conf.CurrentContext = currentContextName
		conf.Contexts[currentContextName].Manifest = defaultString

		context, err = conf.GetCurrentContext()
		require.NoError(t, err)
		assert.Equal(t, conf.Contexts[currentContextName], context)
	})

	t.Run("getCurrentContextIncomplete", func(t *testing.T) {
		conf, cleanup := testutil.InitConfig(t)
		defer cleanup(t)

		context, err := conf.GetCurrentContext()
		require.Error(t, err)
		assert.Nil(t, context)

		conf.CurrentContext = currentContextName

		context, err = conf.GetCurrentContext()
		assert.Error(t, err)
		assert.Nil(t, context)
	})
}

func TestCurrentContextCluster(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	clusterName := "def"
	clusterType := "ephemeral"

	cluster, err := conf.CurrentContextCluster()
	require.Error(t, err)
	assert.Nil(t, cluster)

	conf.CurrentContext = currentContextName
	conf.Contexts[currentContextName].Manifest = defaultString
	conf.Contexts[currentContextName].KubeContext().Cluster = clusterName

	cluster, err = conf.CurrentContextCluster()
	require.NoError(t, err)
	assert.Equal(t, conf.Clusters[clusterName].ClusterTypes[clusterType], cluster)
}

func TestCurrentContextAuthInfo(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	authInfo, err := conf.CurrentContextAuthInfo()
	require.Error(t, err)
	assert.Nil(t, authInfo)

	conf.CurrentContext = currentContextName
	conf.Contexts[currentContextName].Manifest = defaultString

	authInfo, err = conf.CurrentContextAuthInfo()
	require.NoError(t, err)
	assert.Equal(t, conf.AuthInfos["k-admin"], authInfo)
}

func TestCurrentContextManifest(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	clusterName := "def"

	manifest, err := conf.CurrentContextManifest()
	require.Error(t, err)
	assert.Nil(t, manifest)

	conf.CurrentContext = currentContextName
	conf.Contexts[currentContextName].Manifest = defaultString
	conf.Contexts[currentContextName].KubeContext().Cluster = clusterName

	manifest, err = conf.CurrentContextManifest()
	require.NoError(t, err)
	assert.Equal(t, conf.Manifests[defaultString], manifest)
}

func TestCurrentContextEntryPoint(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	clusterName := "def"
	clusterType := "ephemeral"

	entryPoint, err := conf.CurrentContextEntryPoint(clusterType, defaultString)
	require.Error(t, err)
	assert.Equal(t, "", entryPoint)

	conf.CurrentContext = currentContextName
	conf.Contexts[currentContextName].Manifest = defaultString
	conf.Contexts[currentContextName].KubeContext().Cluster = clusterName

	entryPoint, err = conf.CurrentContextEntryPoint(clusterType, defaultString)
	require.NoError(t, err)
	assert.Nil(t, nil, entryPoint)
}

// AuthInfo Related

func TestGetAuthInfos(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	authinfos := conf.GetAuthInfos()
	assert.Len(t, authinfos, 3)
}

func TestGetAuthInfo(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	authinfo, err := conf.GetAuthInfo("def-user")
	require.NoError(t, err)

	// Test Positives
	assert.EqualValues(t, authinfo.KubeAuthInfo().Username, "dummy_username")

	// Test Wrong Cluster
	_, err = conf.GetAuthInfo("unknown")
	assert.Error(t, err)
}

func TestAddAuthInfo(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	co := testutil.DummyAuthInfoOptions()
	authinfo := conf.AddAuthInfo(co)
	assert.EqualValues(t, conf.AuthInfos[co.Name], authinfo)
}

func TestModifyAuthInfo(t *testing.T) {
	conf, cleanup := testutil.InitConfig(t)
	defer cleanup(t)

	co := testutil.DummyAuthInfoOptions()
	authinfo := conf.AddAuthInfo(co)

	co.Username += stringDelta
	co.Password += stringDelta
	co.ClientCertificate += stringDelta
	co.ClientKey += stringDelta
	co.Token += stringDelta
	conf.ModifyAuthInfo(authinfo, co)
	assert.EqualValues(t, conf.AuthInfos[co.Name].KubeAuthInfo().Username, co.Username)
	assert.EqualValues(t, conf.AuthInfos[co.Name].KubeAuthInfo().Password, co.Password)
	assert.EqualValues(t, conf.AuthInfos[co.Name].KubeAuthInfo().ClientCertificate, co.ClientCertificate)
	assert.EqualValues(t, conf.AuthInfos[co.Name].KubeAuthInfo().ClientKey, co.ClientKey)
	assert.EqualValues(t, conf.AuthInfos[co.Name].KubeAuthInfo().Token, co.Token)
	assert.EqualValues(t, conf.AuthInfos[co.Name], authinfo)
}
