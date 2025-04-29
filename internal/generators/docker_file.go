package generators

import "github.com/lunagic/hephaestus/internal/state"

type Dockerfile struct{}

func (generator Dockerfile) Output(s *state.State) error {
	return nil
}
