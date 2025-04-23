package generators

import "github.com/lunagic/hephaestus/internal/state"

type GitHubWorkflow struct{}

func (generator GitHubWorkflow) Output(s *state.State) error {
	return nil
}
