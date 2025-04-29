package generators

import (
	"os"

	"github.com/lunagic/hephaestus/internal/formats"
	"github.com/lunagic/hephaestus/internal/state"
)

type GitIgnore struct {
	file formats.GitIgnore
}

func (generator GitIgnore) Output(s *state.State) error {
	generator.file.Sections = append(generator.file.Sections, &formats.GitIgnoreSection{
		Title: "System",
		Items: []string{
			".DS_Store",
		},
	})
	generator.file.Sections = append(generator.file.Sections, &formats.GitIgnoreSection{
		Title: "Temporary Files",
		Items: []string{
			"/tmp/",
		},
	})

	generator.file.Sections = append(generator.file.Sections, &formats.GitIgnoreSection{
		Title: "Secrets",
		Items: []string{
			".env.local",
			".env.*.local",
		},
	})

	if s.Go.Enabled() {
		generator.file.Sections = append(generator.file.Sections, &formats.GitIgnoreSection{
			Title: "Go",
			Items: []string{
				"__debug_bin*",
			},
		})
	}

	if s.NPM.Enabled() {
		npmIgnores := []string{
			"/node_modules/",
		}

		if s.NPM.HasDependency("next") {
			npmIgnores = append(
				npmIgnores,
				"out",
				".next",
				"next-env.d.ts",
			)
		}
		generator.file.Sections = append(generator.file.Sections, &formats.GitIgnoreSection{
			Title: "npm",
			Items: npmIgnores,
		})

	}

	file, err := os.Create(".gitignore")
	if err != nil {
		return err
	}

	return generator.file.Build(file)
}
