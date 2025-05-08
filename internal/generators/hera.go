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

	if s.Go.Enabled() && s.Hephaestus.HTTPEnabled() {
		heraConfig.Services["backend"] = &hera.Service{
			Command: "make dev-go",
			Watch: []string{
				"Makefile",
				"go.mod",
				"main.go",
				".env.local",
			},
			Exclude: []string{},
		}
	}

	if s.NPM.Enabled() && s.Hephaestus.HTTPEnabled() {
		heraConfig.Services["frontend"] = &hera.Service{
			Command: "make dev-npm",
			Watch: []string{
				"Makefile",
				"package.json",
				"package-lock.json",
				"tsconfig.json",
				"vite.config.ts",
			},
			Exclude: []string{},
		}
	}

	if s.NPM.HasScript("storybook") {
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
