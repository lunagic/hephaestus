package state

import (
	"os"

	"gopkg.in/yaml.v3"
)

func NewHephaestus() (*Hephaestus, error) {
	hephaestusConfigFile, err := os.Open(".config/hephaestus.yaml")
	if err != nil {
		return &Hephaestus{}, nil
	}

	hephaestusConfig := &Hephaestus{}
	if err := yaml.NewDecoder(hephaestusConfigFile).Decode(hephaestusConfig); err != nil {
		return nil, err
	}

	return hephaestusConfig, nil
}

type Hephaestus struct {
	DefaultPort    int                 `yaml:"default_port"`
	Gitignore      []string            `yaml:"gitignore"`
	HeraWatchPaths map[string][]string `yaml:"hera_watch_paths"`
}

func (hephaestus *Hephaestus) GitCleanExcludes() []string {
	return []string{
		".env.local",
		".env.*.local",
	}
}

func (hephaestus *Hephaestus) HTTPEnabled() bool {
	return hephaestus.DefaultPort != 0
}
