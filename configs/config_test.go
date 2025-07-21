package configs_test

import (
	"os"
	"testing"

	"github.com/LiquidCats/rater/configs"
	"github.com/go-playground/sensitive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSecretFromFile(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		// Create a temporary file with a known secret
		content := "my-secret-value"
		tmpFile, err := os.CreateTemp(t.TempDir(), "secret-*")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(content)
		require.NoError(t, err)
		defer tmpFile.Close()

		cfg := configs.CoinMarketCapConfig{
			SecretFile: tmpFile.Name(),
		}

		secret, err := cfg.GetSecret()
		require.NoError(t, err)
		assert.Equal(t, content, string(secret))
	})

	t.Run("file does not exist", func(t *testing.T) {
		cfg := configs.CoinMarketCapConfig{
			SecretFile: t.TempDir() + "/nonexistent.file",
		}

		secret, err := cfg.GetSecret()
		require.Error(t, err)
		assert.Equal(t, sensitive.String(""), secret)
	})
}
