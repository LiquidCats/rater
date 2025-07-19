package docker

import (
	"os"
	"strings"

	"github.com/rotisserie/eris"
)

func GetSecret(name string) (string, error) {
	if strings.HasPrefix(name, "/") {
		data, err := os.ReadFile(name)
		if err != nil {
			return "", eris.Wrap(err, "failed to open secret file")
		}

		return strings.TrimSpace(string(data)), nil
	}

	return name, nil
}
