package generators

import "github.com/lunagic/hephaestus/internal/state"

type VSCodeExtensions struct{}

func (generator VSCodeExtensions) Output(s *state.State) error {
	payload := map[string]any{}

	recommendations := []string{
		"streetsidesoftware.code-spell-checker",
	}

	if s.Go.Enabled() {
		recommendations = append(recommendations, "golang.go")
	}

	if s.Node.Enabled() {
		if s.Node.HasDependency("@biomejs/biome") {
			recommendations = append(recommendations, "biomejs.biome")
		}
	}

	if len(recommendations) > 0 {
		payload["recommendations"] = recommendations
	}

	if len(payload) == 0 {
		return nil
	}

	return writeJSONFile(".vscode/extensions.json", payload)
}
