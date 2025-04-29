package generators

import "github.com/lunagic/hephaestus/internal/state"

type GitHubWorkflow struct {
	Name string   `yaml:"name"`
	On   struct{} `yaml:"on"`
	Push struct{} `yaml:"push"`
	Jobs struct{} `yaml:"jobs"`
}

func (generator GitHubWorkflow) Output(s *state.State) error {
	return nil
}
