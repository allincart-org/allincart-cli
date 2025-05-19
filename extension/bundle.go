package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/shyim/go-version"
)

type AllincartBundle struct {
	path     string
	Composer allincartBundleComposerJson
	config   *Config
}

func newAllincartBundle(path string) (*AllincartBundle, error) {
	composerJsonFile := fmt.Sprintf("%s/composer.json", path)
	if _, err := os.Stat(composerJsonFile); err != nil {
		return nil, err
	}

	jsonFile, err := os.ReadFile(composerJsonFile)
	if err != nil {
		return nil, fmt.Errorf("newAllincartBundle: %v", err)
	}

	var composerJson allincartBundleComposerJson
	err = json.Unmarshal(jsonFile, &composerJson)
	if err != nil {
		return nil, fmt.Errorf("newAllincartBundle: %v", err)
	}

	if composerJson.Type != "allincart-bundle" {
		return nil, fmt.Errorf("newAllincartBundle: composer.json type is not allincart-bundle")
	}

	if composerJson.Extra.BundleName == "" {
		return nil, fmt.Errorf("composer.json does not contain allincart-bundle-name in extra")
	}

	cfg, err := readExtensionConfig(path)
	if err != nil {
		return nil, fmt.Errorf("newAllincartBundle: %v", err)
	}

	extension := AllincartBundle{
		Composer: composerJson,
		path:     path,
		config:   cfg,
	}

	return &extension, nil
}

type composerAutoload struct {
	Psr4 map[string]string `json:"psr-4"`
}

type allincartBundleComposerJson struct {
	Name     string                           `json:"name"`
	Type     string                           `json:"type"`
	License  string                           `json:"license"`
	Version  string                           `json:"version"`
	Require  map[string]string                `json:"require"`
	Extra    allincartBundleComposerJsonExtra `json:"extra"`
	Suggest  map[string]string                `json:"suggest"`
	Autoload composerAutoload                 `json:"autoload"`
}

type allincartBundleComposerJsonExtra struct {
	BundleName string `json:"allincart-bundle-name"`
}

func (p AllincartBundle) GetComposerName() (string, error) {
	return p.Composer.Name, nil
}

// GetRootDir returns the src directory of the bundle.
func (p AllincartBundle) GetRootDir() string {
	return path.Join(p.path, "src")
}

func (p AllincartBundle) GetSourceDirs() []string {
	var result []string

	for _, val := range p.Composer.Autoload.Psr4 {
		result = append(result, path.Join(p.path, val))
	}

	return result
}

// GetResourcesDir returns the resources directory of the allincart bundle.
func (p AllincartBundle) GetResourcesDir() string {
	return path.Join(p.GetRootDir(), "Resources")
}

func (p AllincartBundle) GetResourcesDirs() []string {
	var result []string

	for _, val := range p.GetSourceDirs() {
		result = append(result, path.Join(val, "Resources"))
	}

	return result
}

func (p AllincartBundle) GetName() (string, error) {
	return p.Composer.Extra.BundleName, nil
}

func (p AllincartBundle) GetExtensionConfig() *Config {
	return p.config
}

func (p AllincartBundle) GetAllincartVersionConstraint() (*version.Constraints, error) {
	if p.config != nil && p.config.Build.AllincartVersionConstraint != "" {
		constraint, err := version.NewConstraint(p.config.Build.AllincartVersionConstraint)
		if err != nil {
			return nil, err
		}

		return &constraint, nil
	}

	allincartConstraintString, ok := p.Composer.Require["allincart/core"]

	if !ok {
		return nil, fmt.Errorf("require.allincart/core is required")
	}

	allincartConstraint, err := version.NewConstraint(allincartConstraintString)
	if err != nil {
		return nil, err
	}

	return &allincartConstraint, err
}

func (AllincartBundle) GetType() string {
	return TypeAllincartBundle
}

func (p AllincartBundle) GetVersion() (*version.Version, error) {
	return version.NewVersion(p.Composer.Version)
}

func (p AllincartBundle) GetChangelog() (*ExtensionChangelog, error) {
	return parseExtensionMarkdownChangelog(p)
}

func (p AllincartBundle) GetLicense() (string, error) {
	return p.Composer.License, nil
}

func (p AllincartBundle) GetPath() string {
	return p.path
}

func (p AllincartBundle) GetMetaData() *extensionMetadata {
	return &extensionMetadata{
		Label: extensionTranslated{
			German:  "FALLBACK",
			English: "FALLBACK",
		},
		Description: extensionTranslated{
			German:  "FALLBACK",
			English: "FALLBACK",
		},
	}
}

func (p AllincartBundle) Validate(c context.Context, ctx *ValidationContext) {
}
