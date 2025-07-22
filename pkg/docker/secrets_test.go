package docker_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LiquidCats/rater/pkg/docker"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSecret(t *testing.T) {
	t.Run("returns name when not starting with slash", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "simple string",
				input:    "my-secret",
				expected: "my-secret",
			},
			{
				name:     "empty string",
				input:    "",
				expected: "",
			},
			{
				name:     "string with spaces",
				input:    "my secret value",
				expected: "my secret value",
			},
			{
				name:     "string with special characters",
				input:    "secret@123!",
				expected: "secret@123!",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := docker.GetSecret(tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("reads file when starting with slash", func(t *testing.T) {
		// Create a temporary directory for test files
		tempDir := t.TempDir()

		tests := []struct {
			name         string
			fileContent  string
			expectedRead string
		}{
			{
				name:         "file with simple content",
				fileContent:  "secret-value",
				expectedRead: "secret-value",
			},
			{
				name:         "file with trailing whitespace",
				fileContent:  "secret-value\n\n",
				expectedRead: "secret-value",
			},
			{
				name:         "file with leading and trailing whitespace",
				fileContent:  "  \n secret-value \t\n",
				expectedRead: "secret-value",
			},
			{
				name:         "file with only whitespace",
				fileContent:  "   \n\t\n   ",
				expectedRead: "",
			},
			{
				name:         "empty file",
				fileContent:  "",
				expectedRead: "",
			},
			{
				name:         "file with multiline content",
				fileContent:  "line1\nline2\nline3\n",
				expectedRead: "line1\nline2\nline3",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create a test file
				filePath := filepath.Join(tempDir, "test-secret")
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				require.NoError(t, err)

				// Test reading the file
				result, err := docker.GetSecret(filePath)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedRead, result)

				// Clean up
				os.Remove(filePath)
			})
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		nonExistentPath := "/tmp/non-existent-secret-file-12345"

		result, err := docker.GetSecret(nonExistentPath)
		require.Error(t, err)
		assert.Empty(t, result)

		// Check that the error is wrapped with eris
		assert.True(t, eris.Is(err, os.ErrNotExist))
		assert.Contains(t, err.Error(), "failed to open secret file")
	})

	t.Run("returns error for unreadable file", func(t *testing.T) {
		// This test requires specific permissions handling
		// Skip on Windows as permission handling is different
		if os.Getenv("GOOS") == "windows" {
			t.Skip("Skipping permission test on Windows")
		}

		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "unreadable-secret")

		// Create a file
		err := os.WriteFile(filePath, []byte("secret"), 0644)
		require.NoError(t, err)

		// Make it unreadable
		err = os.Chmod(filePath, 0000)
		require.NoError(t, err)

		// Ensure we restore permissions for cleanup
		defer os.Chmod(filePath, 0644)

		result, err := docker.GetSecret(filePath)
		require.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to open secret file")
	})

	t.Run("handles directory path", func(t *testing.T) {
		tempDir := t.TempDir()

		result, err := docker.GetSecret(tempDir)
		require.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to open secret file")
	})
}
