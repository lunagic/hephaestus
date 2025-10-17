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
		Name: "Main",
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
								With: map[string]any{
									"go-version": s.Go.Version(),
								},
							},
						)
					}

					if s.Node.Enabled() {
						steps = append(
							steps,
							formats.GitHubWorkflowStep{
								Uses: "actions/setup-node@v4",
								With: map[string]any{
									"node-version": s.Node.Version(),
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

	if s.Hephaestus.DockerImage != "" {
		validateWorkflow.Jobs["publish"] = formats.GitHubWorkflowJob{
			Needs:  "validate",
			If:     "github.ref == 'refs/heads/main'",
			RunsOn: "ubuntu-latest",
			Steps: []formats.GitHubWorkflowStep{
				{
					Uses: "actions/checkout@v4",
				},

				{
					Uses: "docker/setup-buildx-action@v3",
				},
				{
					Uses: "docker/login-action@v3",
					With: map[string]any{
						"registry": "ghcr.io",
						"username": "${{ github.actor }}",
						"password": "${{ secrets.GITHUB_TOKEN }}",
					},
				},
				{
					Uses: "docker/build-push-action@v5",
					With: map[string]any{
						"context":   ".",
						"push":      true,
						"platforms": "linux/amd64,linux/arm64",
						"tags":      "ghcr.io/${{ github.repository }}:latest",
					},
				},
			},
		}
	}

	validateWorkflowFile, err := os.Create(".github/workflows/main.yml")
	if err != nil {
		return err
	}

	if err := utils.YAML(validateWorkflowFile, validateWorkflow); err != nil {
		return nil
	}

	return nil
}
