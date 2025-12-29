package sklairConfig

import (
	"os"
	"path/filepath"
)

// TODO: check TODO.md for more info about this

type GlobalConfig struct {
	CheckForUpdates bool `json:"checkForUpdates,omitempty"`
}

var defaultGlobalConfig = GlobalConfig{
	CheckForUpdates: true,
}

func GlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".sklair/config.json"), nil
}
