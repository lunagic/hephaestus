package generators

import (
	"os"
	"path"

	"github.com/lunagic/hephaestus/internal/state"
	"github.com/lunagic/hephaestus/internal/utils"
)

type VSCodeSettings struct{}

func (generator VSCodeSettings) Output(s *state.State) error {
	payload := map[string]any{}
	if s.Go.Enabled() {
		payload["go.lintTool"] = "golangci-lint"

	}
	if s.Node.Enabled() {
		payload["javascript.preferences.importModuleSpecifier"] = "non-relative"
		payload["typescript.preferences.importModuleSpecifier"] = "non-relative"

		if s.Node.HasDependency("@biomejs/biome") {
			for _, lang := range []string{
				"[css]",
				"[javascript]",
				"[javascriptreact]",
				"[json]",
				"[jsonc]",
				"[scss]",
				"[typescript]",
				"[typescriptreact]",
			} {
				payload[lang] = map[string]any{
					"editor.defaultFormatter": "biomejs.biome",
				}
			}
			payload["editor.codeActionsOnSave"] = map[string]any{
				"source.action.useSortedKeys.biome": "explicit",
				"source.fixAll.biome":               "explicit",
			}
			payload["biome.enabled"] = true
		}
	}

	if len(payload) == 0 {
		return nil
	}

	return writeJSONFile(".vscode/settings.json", payload)
}

func writeJSONFile(targetPath string, payload any) error {
	dir := path.Dir(targetPath)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	jsonString, err := utils.MarshalSortedJSON(payload)
	if err != nil {
		return err
	}

	if _, err := file.Write(jsonString); err != nil {
		return err
	}

	// jsonBytes, err := json.MarshalIndent(payload, "", "\t")
	// if err != nil {
	// 	return err
	// }

	// jsonBytes = append(jsonBytes, []byte("\n")...)

	// if _, err := file.Write(jsonBytes); err != nil {
	// 	return err
	// }

	return nil
}
