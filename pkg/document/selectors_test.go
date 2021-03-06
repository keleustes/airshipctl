package document_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/testutil"
)

func TestSelectorsPositive(t *testing.T) {
	bundle := testutil.NewTestBundle(t, "testdata/selectors/valid")

	t.Run("TestEphemeralCloudDataSelector", func(t *testing.T) {
		doc, err := bundle.Select(document.NewEphemeralCloudDataSelector())
		require.NoError(t, err)
		assert.Len(t, doc, 1)
	})

	t.Run("TestEphemeralNetworkDataSelector", func(t *testing.T) {
		docs, err := bundle.Select(document.NewEphemeralBMHSelector())
		require.NoError(t, err)
		assert.Len(t, docs, 1)
		bmhDoc := docs[0]
		selector, err := document.NewEphemeralNetworkDataSelector(bmhDoc)
		require.NoError(t, err)
		assert.Equal(t, "validName", selector.Name)
	})

	t.Run("TestEphemeralCloudDataSelector", func(t *testing.T) {
		doc, err := bundle.Select(document.NewEphemeralCloudDataSelector())
		require.NoError(t, err)
		assert.Len(t, doc, 1)
	})
}

func TestSelectorsNegative(t *testing.T) {
	// These two tests take bundle with two malformed documents
	// each of the documents will fail at different locations providing higher
	// test coverage
	bundle := testutil.NewTestBundle(t, "testdata/selectors/invalid")

	t.Run("TestNewEphemeralNetworkDataSelectorErr", func(t *testing.T) {
		docs, err := bundle.Select(document.NewEphemeralBMHSelector())
		require.NoError(t, err)
		assert.Len(t, docs, 2)
		bmhDoc := docs[0]
		_, err = document.NewEphemeralNetworkDataSelector(bmhDoc)
		assert.Error(t, err)
	})

	t.Run("TestEphemeralNetworkDataSelectorErr", func(t *testing.T) {
		docs, err := bundle.Select(document.NewEphemeralBMHSelector())
		require.NoError(t, err)
		assert.Len(t, docs, 2)
		bmhDoc := docs[1]
		_, err = document.NewEphemeralNetworkDataSelector(bmhDoc)
		assert.Error(t, err)
	})
}
