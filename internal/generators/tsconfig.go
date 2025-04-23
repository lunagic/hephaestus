package generators

import "github.com/lunagic/hephaestus/internal/state"

type TSConfig struct{}

func (generator TSConfig) Output(s *state.State) error {
	return nil
}
