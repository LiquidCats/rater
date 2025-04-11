package configs

import (
	"os"
	"testing"

	"github.com/go-playground/sensitive"
	"github.com/stretchr/testify/assert"
)

func TestGetSecretFromFile(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		// Create a temporary file with a known secret
		content := "my-secret-value"
		tmpFile, err := os.CreateTemp("", "secret-*")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name()) // Clean up

		_, err = tmpFile.Write([]byte(content))
		assert.NoError(t, err)
		defer tmpFile.Close()

		// Call the function
		secret, err := getSecretFromFile(tmpFile.Name())
		assert.NoError(t, err)
		assert.Equal(t, sensitive.String(content), secret)
	})

	t.Run("file does not exist", func(t *testing.T) {
		secret, err := getSecretFromFile("nonexistent.file")
		assert.Error(t, err)
		assert.Equal(t, sensitive.String(""), secret)
	})
}
