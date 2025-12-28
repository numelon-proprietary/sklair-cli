package sklairConfig

import (
	"encoding/json"
	"os"
)

// TODO: expand this struct when JS obfuscation is added
type ObfuscateJS struct {
	Enabled bool `json:"enabled,omitempty"`
}

type PreventFOUC struct {
	Enabled bool   `json:"enabled,omitempty"`
	Colour  string `json:"colour,omitempty"`
}

//type ResourceHints struct {
//	Enabled    bool   `json:"enabled,omitempty"`
//	SiteOrigin string `json:"siteOrigin,omitempty"`
//}

type ProjectConfig struct {
	Input      string `json:"input,omitempty"`
	Components string `json:"components,omitempty"`

	Exclude        []string `json:"exclude,omitempty"`
	ExcludeCompile []string `json:"excludeCompile,omitempty"`

	Output string `json:"output,omitempty"`

	Minify      bool         `json:"minify,omitempty"`
	ObfuscateJS *ObfuscateJS `json:"obfuscateJS,omitempty"`

	PreventFOUC *PreventFOUC `json:"PreventFOUC,omitempty"`
	//ResourceHints *ResourceHints `json:"resourceHints,omitempty"` // TODO: in sklair init, add ResourceHints to the questionnaire
}

var DefaultConfig = ProjectConfig{
	Input:      "./",
	Components: "./components",

	Exclude:        []string{},
	ExcludeCompile: []string{},

	Output: "./build",

	Minify: false,
	ObfuscateJS: &ObfuscateJS{
		Enabled: false,
	},

	PreventFOUC: &PreventFOUC{
		Enabled: false,
		Colour:  "#202020",
	},
	//ResourceHints: &ResourceHints{
	//	Enabled:    false,
	//	SiteOrigin: "https://sklair.numelon.com", // TODO: maybe just make it empty by default
	//},
}

func LoadProject(path string) (*ProjectConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
