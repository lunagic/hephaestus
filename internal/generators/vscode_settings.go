package generators

import "github.com/lunagic/hephaestus/internal/state"

type VSCodeSettings struct{}

func (generator VSCodeSettings) Output(s *state.State) error {
	return nil
}
