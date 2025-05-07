package packagist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillAuthStruct(t *testing.T) {
	// Test with empty struct
	auth := &ComposerAuth{}
	filledAuth := fillAuthStruct(auth)

	assert.NotNil(t, filledAuth.HTTPBasicAuth)
	assert.NotNil(t, filledAuth.BearerAuth)
	assert.NotNil(t, filledAuth.GitlabAuth)
	assert.NotNil(t, filledAuth.GithubOAuth)
	assert.NotNil(t, filledAuth.BitbucketOauth)

	// Test with partially filled struct
	auth = &ComposerAuth{
		HTTPBasicAuth: map[string]ComposerAuthHttpBasic{
			"example.org": {
				Username: "user",
				Password: "pass",
			},
		},
	}
	filledAuth = fillAuthStruct(auth)

	assert.NotNil(t, filledAuth.HTTPBasicAuth)
	assert.Equal(t, "user", filledAuth.HTTPBasicAuth["example.org"].Username)
	assert.Equal(t, "pass", filledAuth.HTTPBasicAuth["example.org"].Password)
	assert.NotNil(t, filledAuth.BearerAuth)
	assert.NotNil(t, filledAuth.GitlabAuth)
	assert.NotNil(t, filledAuth.GithubOAuth)
	assert.NotNil(t, filledAuth.BitbucketOauth)
}

func TestComposerAuthSave(t *testing.T) {
	tempDir := t.TempDir()
	authFile := filepath.Join(tempDir, "auth.json")

	auth := &ComposerAuth{
		path: authFile,
		HTTPBasicAuth: map[string]ComposerAuthHttpBasic{
			"example.org": {
				Username: "user",
				Password: "pass",
			},
		},
		BearerAuth: map[string]string{
			"api.example.org": "token123",
		},
	}

	err := auth.Save()
	assert.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(authFile)
	assert.NoError(t, err)

	// Read and verify content
	content, err := os.ReadFile(authFile)
	assert.NoError(t, err)

	var savedAuth ComposerAuth
	err = json.Unmarshal(content, &savedAuth)
	assert.NoError(t, err)

	assert.Equal(t, auth.HTTPBasicAuth["example.org"].Username, savedAuth.HTTPBasicAuth["example.org"].Username)
	assert.Equal(t, auth.HTTPBasicAuth["example.org"].Password, savedAuth.HTTPBasicAuth["example.org"].Password)
	assert.Equal(t, auth.BearerAuth["api.example.org"], savedAuth.BearerAuth["api.example.org"])
}

func TestReadComposerAuth(t *testing.T) {
	// Test with existing file
	t.Run("existing file", func(t *testing.T) {
		tempDir := t.TempDir()
		authFile := filepath.Join(tempDir, "auth.json")

		testAuth := ComposerAuth{
			HTTPBasicAuth: map[string]ComposerAuthHttpBasic{
				"example.org": {
					Username: "user",
					Password: "pass",
				},
			},
			BearerAuth: map[string]string{
				"api.example.org": "token123",
			},
		}

		content, err := json.MarshalIndent(testAuth, "", "  ")
		assert.NoError(t, err)
		err = os.WriteFile(authFile, content, 0o644)
		assert.NoError(t, err)

		auth, err := ReadComposerAuth(authFile, false)
		assert.NoError(t, err)
		assert.Equal(t, authFile, auth.path)
		assert.Equal(t, "user", auth.HTTPBasicAuth["example.org"].Username)
		assert.Equal(t, "pass", auth.HTTPBasicAuth["example.org"].Password)
		assert.Equal(t, "token123", auth.BearerAuth["api.example.org"])
	})

	// Test with non-existing file, with fallback
	t.Run("non-existing file with fallback", func(t *testing.T) {
		tempDir := t.TempDir()
		authFile := filepath.Join(tempDir, "nonexistent.json")

		auth, err := ReadComposerAuth(authFile, true)
		assert.NoError(t, err)
		assert.Equal(t, authFile, auth.path)
		assert.NotNil(t, auth.HTTPBasicAuth)
		assert.NotNil(t, auth.BearerAuth)
		assert.NotNil(t, auth.GitlabAuth)
		assert.NotNil(t, auth.GithubOAuth)
		assert.NotNil(t, auth.BitbucketOauth)
	})

	// Test with non-existing file, without fallback
	t.Run("non-existing file without fallback", func(t *testing.T) {
		tempDir := t.TempDir()
		authFile := filepath.Join(tempDir, "nonexistent.json")

		auth, err := ReadComposerAuth(authFile, false)
		assert.Error(t, err)
		assert.Nil(t, auth)
	})

	// Test with invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		tempDir := t.TempDir()
		authFile := filepath.Join(tempDir, "invalid.json")

		err := os.WriteFile(authFile, []byte("{invalid json}"), 0o644)
		assert.NoError(t, err)

		auth, err := ReadComposerAuth(authFile, false)
		assert.Error(t, err)
		assert.Nil(t, auth)
	})
}
