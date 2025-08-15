package state

import (
	"errors"
	"os"

	"github.com/lunagic/hephaestus/internal/formats"
)

func NewNode() (*Node, error) {
	s := &Node{
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

type Node struct {
	packageJSON *formats.PackageJSON
}

func (s *Node) Enabled() bool {
	return s.packageJSON.ReadFromDisk() == nil
}

func (s *Node) Version() string {
	return s.packageJSON.NodeVersion()
}

func (s *Node) HasScript(script string) bool {
	return s.packageJSON.GetScript(script) != ""
}

func (s *Node) HasDependency(packageName string) bool {
	return s.packageJSON.HasPackage(packageName)
}
