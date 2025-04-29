package generators

import (
	"os"

	"github.com/lunagic/hephaestus/internal/formats"
	"github.com/lunagic/hephaestus/internal/state"
	"github.com/lunagic/hephaestus/internal/utils"
)

type GitHubWorkflow struct{}

func (generator GitHubWorkflow) Output(s *state.State) error {
	if err := os.MkdirAll(".github/workflows", 0755); err != nil {
		return err
	}

	validateWorkflow := formats.GitHubWorkflow{
		Name: "Validate",
		On: map[string]formats.GitHubWorkflowEvent{
			"pull_request": {},
			"push": {
				Branches: []string{
					"main",
				},
			},
		},
		Jobs: map[string]formats.GitHubWorkflowJob{
			"validate": {
				RunsOn: "ubuntu-latest",
				Steps: func() []formats.GitHubWorkflowStep {
					steps := []formats.GitHubWorkflowStep{
						{
							Uses: "actions/checkout@v4",
						},
					}

					if s.Go.Enabled() {
						steps = append(
							steps,
							formats.GitHubWorkflowStep{
								Uses: "actions/setup-go@v5",
								With: map[string]string{
									"go-version": s.Go.Version(),
								},
							},
						)
					}

					if s.NPM.Enabled() {
						steps = append(
							steps,
							formats.GitHubWorkflowStep{
								Uses: "actions/setup-node@v4",
								With: map[string]string{
									"node-version": "20",
								},
							},
						)
					}

					steps = append(
						steps,
						formats.GitHubWorkflowStep{
							Run: "make",
						},
					)

					return steps
				}(),
			},
		},
	}

	validateWorkflowFile, err := os.Create(".github/workflows/validate.yml")
	if err != nil {
		return err
	}

	if err := utils.YAML(validateWorkflowFile, validateWorkflow); err != nil {
		return nil
	}

	return nil
}
