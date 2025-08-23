package generators

import "github.com/lunagic/hephaestus/internal/state"

type VSCodeLaunch struct{}

func (generator VSCodeLaunch) Output(s *state.State) error {
	payload := map[string]any{
		"version": "0.2.0",
	}

	configurations := []any{}

	if s.Go.Enabled() {
		configurations = append(configurations, map[string]any{
			"name":    "Launch Main",
			"type":    "go",
			"request": "launch",
			"mode":    "auto",
			"program": ".",
		})
	}

	if len(configurations) == 0 {
		return nil
	}

	payload["configurations"] = configurations

	return writeJSONFile(".vscode/launch.json", payload)
}
