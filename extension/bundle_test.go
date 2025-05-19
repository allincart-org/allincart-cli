package extension

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBundleEmptyFolder(t *testing.T) {
	dir := t.TempDir()

	bundle, err := newAllincartBundle(dir)
	assert.Error(t, err)
	assert.Nil(t, bundle)
}

func TestCreateBundleInvalidComposerType(t *testing.T) {
	dir := t.TempDir()

	// Create composer.json
	composer := []byte(`{
		"name": "allincart/invalid",
		"type": "invalid"
	}
	`)
	_ = os.WriteFile(path.Join(dir, "composer.json"), composer, 0o644)

	bundle, err := newAllincartBundle(dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "composer.json type is not allincart-bundle")
	assert.Nil(t, bundle)
}

func TestCreateBundleMissingName(t *testing.T) {
	dir := t.TempDir()

	// Create composer.json
	composer := []byte(`{
		"name": "allincart/invalid",
		"type": "allincart-bundle"
	}
	`)
	_ = os.WriteFile(path.Join(dir, "composer.json"), composer, 0o644)

	bundle, err := newAllincartBundle(dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "composer.json does not contain allincart-bundle-name")
	assert.Nil(t, bundle)
}

func TestCreateBundle(t *testing.T) {
	dir := t.TempDir()

	// Create composer.json
	composer := []byte(`{
		"name": "allincart/invalid",
		"version": "1.0.0",
		"type": "allincart-bundle",
		"extra": {
			"allincart-bundle-name": "TestBundle"
		},
		"autoload": {
			"psr-4": {
				"TestBundle\\": "src/"
			}
		}
	}
	`)
	_ = os.WriteFile(path.Join(dir, "composer.json"), composer, 0o644)

	bundle, err := newAllincartBundle(dir)
	assert.NoError(t, err)

	name, err := bundle.GetName()
	assert.NoError(t, err)

	assert.Equal(t, "TestBundle", name)
	assert.Equal(t, path.Join(dir, "src"), bundle.GetRootDir())
	assert.Equal(t, dir, bundle.GetPath())
	assert.Equal(t, path.Join(dir, "src", "Resources"), bundle.GetResourcesDir())
	assert.Equal(t, path.Join(dir, "src", "Resources"), bundle.GetResourcesDirs()[0])
	assert.Equal(t, TypeAllincartBundle, bundle.GetType())

	_, err = bundle.GetChangelog()
	// changelog is missing
	assert.Error(t, err)

	version, err := bundle.GetVersion()
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", version.String())

	// does notthing
	bundle.Validate(getTestContext(), &ValidationContext{})

	assert.Equal(t, "FALLBACK", bundle.GetMetaData().Description.German)
}
