package discovery

import (
	"os"
	"path/filepath"
	"strings"
)

func DiscoverComponents(source string) (map[string]string, error) {
	dir, err := os.ReadDir(source)
	if err != nil {
		return nil, err
	}

	components := make(map[string]string)

	for _, file := range dir {
		if !file.IsDir() {
			name := file.Name()
			ext := filepath.Ext(name)

			trimmed := strings.TrimSuffix(name, ext)
			components[strings.ToLower(trimmed)] = name
		}
	}

	return components, nil
}
