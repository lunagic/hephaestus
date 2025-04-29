package formats

type GitHubWorkflow struct {
	Name string                         `yaml:"name"`
	On   map[string]GitHubWorkflowEvent `yaml:"on"`
	Jobs map[string]GitHubWorkflowJob   `yaml:"jobs"`
}

type GitHubWorkflowJob struct {
	RunsOn string               `yaml:"runs-on"`
	Steps  []GitHubWorkflowStep `yaml:"steps"`
}

type GitHubWorkflowEvent struct {
	Branches []string `yaml:"branches,omitempty"`
}

type GitHubWorkflowStep struct {
	Name string            `yaml:"name,omitempty"`
	Uses string            `yaml:"uses,omitempty"`
	Run  string            `yaml:"run,omitempty"`
	With map[string]string `yaml:"with,omitempty"`
}
