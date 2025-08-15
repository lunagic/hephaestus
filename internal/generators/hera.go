package generators

import (
	"os"

	"github.com/lunagic/hephaestus/internal/state"
	"github.com/lunagic/hephaestus/internal/utils"
	"github.com/lunagic/hera/hera"
)

type Hera struct{}

func (generator Hera) Output(s *state.State) error {
	heraConfig := hera.Config{
		Services: map[string]*hera.Service{},
	}

	fileChecker := func(paths []string, filePath string) []string {
		if _, err := os.Open(filePath); err == nil {
			paths = append(paths, filePath)
		}

		return paths
	}

	if s.Go.Enabled() && s.Go.AbleToBuild() && s.Hephaestus.HTTPEnabled() {
		watchPaths := []string{
			".env",
			".env.local",
		}

		watchPaths = fileChecker(watchPaths, "Makefile")
		watchPaths = fileChecker(watchPaths, "go.mod")
		watchPaths = fileChecker(watchPaths, "go.sum")
		watchPaths = fileChecker(watchPaths, "main.go")

		watchPaths = append(watchPaths, s.Hephaestus.HeraWatchPaths["backend"]...)

		heraConfig.Services["backend"] = &hera.Service{
			Command: "make dev-go",
			Watch:   watchPaths,
			Exclude: []string{},
		}
	}

	if s.Node.Enabled() && s.Hephaestus.HTTPEnabled() {
		watchPaths := []string{}
		watchPaths = fileChecker(watchPaths, "Makefile")
		watchPaths = fileChecker(watchPaths, "package.json")
		watchPaths = fileChecker(watchPaths, "package-lock.json")
		watchPaths = fileChecker(watchPaths, "tsconfig.json")
		watchPaths = fileChecker(watchPaths, "vite.config.ts")
		watchPaths = fileChecker(watchPaths, "main.tsx")

		watchPaths = append(watchPaths, s.Hephaestus.HeraWatchPaths["frontend"]...)

		heraConfig.Services["frontend"] = &hera.Service{
			Command: "make dev-npm",
			Watch:   watchPaths,
			Exclude: []string{},
		}
	}

	if s.Node.HasScript("storybook") {
		heraConfig.Services["storybook"] = &hera.Service{
			Command: "npm run storybook",
			Watch: []string{
				".storybook",
			},
			Exclude: []string{},
		}
	}

	if len(heraConfig.Services) == 0 {
		return nil
	}

	if err := os.MkdirAll(".config", 0755); err != nil {
		return err
	}
	heraConfigFile, err := os.Create(".config/hera.yaml")
	if err != nil {
		return err
	}

	if err := utils.YAML(heraConfigFile, heraConfig); err != nil {
		return nil
	}

	return nil
}
