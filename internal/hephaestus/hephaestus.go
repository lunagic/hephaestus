package hephaestus

import (
	"github.com/lunagic/hephaestus/internal/generators"
	"github.com/lunagic/hephaestus/internal/state"
)

type Generator interface {
	Output(s *state.State) error
}

type State interface {
	Output(s *state.State) error
}

func Run() error {
	generatorList := []Generator{
		generators.Dockerfile{},
		generators.GitHubWorkflow{},
		generators.GitIgnore{},
		generators.Hera{},
		generators.Makefile{},
		generators.MarkDown{},
		generators.TSConfig{},
		generators.TSConfig{},
		generators.VSCodeExtensions{},
		generators.VSCodeLaunch{},
		generators.VSCodeSettings{},
	}

	s, err := state.New()
	if err != nil {
		return err
	}

	for _, generator := range generatorList {
		if err := generator.Output(s); err != nil {
			return err
		}
	}

	return nil
}
