package state

import (
	"errors"
	"os"

	"github.com/lunagic/hephaestus/internal/formats"
)

func NewNPM() (*NPM, error) {
	s := &NPM{
		packageJSON: &formats.PackageJSON{},
	}

	if err := s.packageJSON.ReadFromDisk(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}

	return s, nil
}

type NPM struct {
	packageJSON *formats.PackageJSON
}

func (s *NPM) Enabled() bool {
	return s.packageJSON.ReadFromDisk() == nil
}

func (s *NPM) Version() string {
	return s.packageJSON.NodeVersion()
}

func (s *NPM) HasScript(script string) bool {
	return s.packageJSON.GetScript(script) != ""
}

func (s *NPM) HasDependency(packageName string) bool {
	return s.packageJSON.HasPackage(packageName)
}
