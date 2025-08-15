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
		goItems := []string{
			"__debug_bin*",
		}

		if _, err := os.Open(".goreleaser.yaml"); err == nil {
			goItems = append(goItems,
				"dist",
			)
		}

		gitignore.Sections = append(gitignore.Sections, &formats.GitIgnoreSection{
			Title: "Go",
			Items: goItems,
		})
	}

	if s.Node.Enabled() {
		npmIgnores := []string{
			"/node_modules/",
		}

		if s.Node.HasDependency("typescript") {
			npmIgnores = append(
				npmIgnores,
				"tsconfig.tsbuildinfo",
			)
		}

		if s.Node.HasDependency("next") {
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

	if len(s.Hephaestus.FrontendOutPaths) > 0 {
		gitignore.Sections = append(gitignore.Sections, &formats.GitIgnoreSection{
			Title: "Frontend Out Paths",
			Items: s.Hephaestus.FrontendOutPaths,
		})
	}

	if len(s.Hephaestus.Gitignore) > 0 {
		gitignore.Sections = append(gitignore.Sections, &formats.GitIgnoreSection{
			Title: "Project Specific",
			Items: s.Hephaestus.Gitignore,
		})
	}

	file, err := os.Create(".gitignore")
	if err != nil {
		return err
	}

	return gitignore.Build(file)
}
