package generators

import (
	"os"

	"github.com/lunagic/hephaestus/internal/formats"
	"github.com/lunagic/hephaestus/internal/state"
)

type GitIgnore struct{}

func (generator GitIgnore) Output(s *state.State) error {
	gitignore := formats.GitIgnore{}

	gitignore.Sections = append(gitignore.Sections, &formats.GitIgnoreSection{
		Title: "System",
		Items: []string{
			".DS_Store",
		},
	})
	gitignore.Sections = append(gitignore.Sections, &formats.GitIgnoreSection{
		Title: "Temporary Files",
		Items: []string{
			"/tmp/",
		},
	})

	gitignore.Sections = append(gitignore.Sections, &formats.GitIgnoreSection{
		Title: "Secrets",
		Items: []string{
			".env.local",
			".env.*.local",
		},
	})

	if s.Go.Enabled() {
		gitignore.Sections = append(gitignore.Sections, &formats.GitIgnoreSection{
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
		gitignore.Sections = append(gitignore.Sections, &formats.GitIgnoreSection{
			Title: "npm",
			Items: npmIgnores,
		})
	}

	file, err := os.Create(".gitignore")
	if err != nil {
		return err
	}

	return gitignore.Build(file)
}
