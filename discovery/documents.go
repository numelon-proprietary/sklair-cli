package discovery

import (
	"os"
	"path/filepath"
	"strings"
)

type DocumentLists struct {
	HtmlFiles   []string
	StaticFiles []string
}

func DocumentDiscovery(root string) (*DocumentLists, error) {
	lists := &DocumentLists{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && filepath.Base(path) == "components" {
			return filepath.SkipDir
		}

		// skip directories since they will be walked by filepath.Walk later anyways
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext == ".html" || ext == ".shtml" {
			lists.HtmlFiles = append(lists.HtmlFiles, path)
		} else {
			lists.StaticFiles = append(lists.StaticFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return lists, nil
}
