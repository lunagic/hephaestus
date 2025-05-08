package generators

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/lunagic/hephaestus/internal/formats"
	"github.com/lunagic/hephaestus/internal/state"
)

type Makefile struct{}

func (generator Makefile) Output(s *state.State) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	projectName := path.Base(pwd)

	m := formats.Makefile{
		Variables: []string{
			"SHELL=/bin/bash -o pipefail",
			"$(shell git config core.hooksPath ops/git-hooks)",
		},
	}

	if s.Go.Enabled() {
		m.Variables = append(
			m.Variables,
			"GO_PATH := $(shell go env GOPATH 2> /dev/null)",
			"PATH := /usr/local/bin:$(GO_PATH)/bin:$(PATH)",
		)
	}

	// Default Target
	defaultTarget := &formats.MakefileTarget{
		Name:          "full",
		BeforeTargets: []string{},
	}
	m.Targets = append(m.Targets, defaultTarget)

	{ // Clean Target
		cleanIgnores := []string{}
		for _, cleanIgnore := range s.Hephaestus.GitCleanExcludes() {
			cleanIgnores = append(cleanIgnores, fmt.Sprintf(`--exclude="!%s"`, cleanIgnore))
		}
		parentTarget := &formats.MakefileTarget{
			Comment: "Clean the project of temporary files",
			Name:    "clean",
			Commands: []string{
				fmt.Sprintf(`git clean -Xdff %s`, strings.Join(cleanIgnores, " ")),
			},
		}
		m.Targets = append(m.Targets, parentTarget)
		defaultTarget.BeforeTargets = append(defaultTarget.BeforeTargets, parentTarget.Name)
	}

	{ // Lint Targets
		parentTarget := &formats.MakefileTarget{
			Comment:       "Lint the project",
			Name:          "lint",
			BeforeTargets: []string{},
		}
		m.Targets = append(m.Targets, parentTarget)
		defaultTarget.BeforeTargets = append(defaultTarget.BeforeTargets, parentTarget.Name)

		if s.NPM.Enabled() && s.NPM.HasScript("lint") {
			localTarget := &formats.MakefileTarget{
				Name: "lint-npm",
				Commands: []string{
					"npm install",
					"npm run lint",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}

		if s.Go.Enabled() {
			localTarget := &formats.MakefileTarget{
				Name: "lint-go",
				Commands: []string{
					"@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.5",
					"go mod tidy",
					"golangci-lint run ./...",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}
	}

	{ // Fix Targets
		parentTarget := &formats.MakefileTarget{
			Comment:       "Fix the project",
			Name:          "fix",
			BeforeTargets: []string{},
		}
		m.Targets = append(m.Targets, parentTarget)

		if s.NPM.Enabled() && s.NPM.HasScript("fix") {
			localTarget := &formats.MakefileTarget{
				Name: "fix-npm",
				Commands: []string{
					"npm install",
					"npm run fix",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}

		if s.Go.Enabled() {
			localTarget := &formats.MakefileTarget{
				Name: "fix-go",
				Commands: []string{
					"go mod tidy",
					"gofmt -s -w .",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}
	}

	{ // Test Targets
		parentTarget := &formats.MakefileTarget{
			Comment:       "Test the project",
			Name:          "test",
			BeforeTargets: []string{},
		}
		m.Targets = append(m.Targets, parentTarget)
		defaultTarget.BeforeTargets = append(defaultTarget.BeforeTargets, parentTarget.Name)

		if s.NPM.Enabled() && s.NPM.HasScript("test") {
			localTarget := &formats.MakefileTarget{
				Name: "test-npm",
				Commands: []string{
					"npm install",
					"npm run test",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}

		if s.Go.Enabled() {
			localTarget := &formats.MakefileTarget{
				Name: "test-go",
				Commands: []string{
					"@go install github.com/boumenot/gocover-cobertura@latest",
					"@mkdir -p tmp/coverage/go/",
					"go test -cover -coverprofile tmp/coverage/go/profile.txt ./...",
					`@go tool cover -func tmp/coverage/go/profile.txt | awk '/^total/{print $$1 " " $$3}'`,
					"@go tool cover -html tmp/coverage/go/profile.txt -o tmp/coverage/go/coverage.html",
					"@gocover-cobertura < tmp/coverage/go/profile.txt > tmp/coverage/go/cobertura-coverage.xml",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}
	}

	{ // Build Targets
		parentTarget := &formats.MakefileTarget{
			Comment:       "Build the project",
			Name:          "build",
			BeforeTargets: []string{},
		}
		m.Targets = append(m.Targets, parentTarget)
		defaultTarget.BeforeTargets = append(defaultTarget.BeforeTargets, parentTarget.Name)

		if s.NPM.Enabled() && s.NPM.HasScript("build") {
			localTarget := &formats.MakefileTarget{
				Name: "build-npm",
				Commands: []string{
					"npm install",
					"npm run build",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}

		if s.Go.Enabled() && s.Go.AbleToBuild() {
			localTarget := &formats.MakefileTarget{
				Name: "build-go",
				Commands: []string{
					"go generate",
					fmt.Sprintf("go build -ldflags='-s -w' -o tmp/build/%s .", projectName),
					"go install .",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}
	}

	{ // Dev Targets
		if s.NPM.Enabled() && s.NPM.HasScript("build") {
			localTarget := &formats.MakefileTarget{
				Name: "dev-npm",
				BeforeTargets: []string{
					"build-npm",
				},
			}

			m.Targets = append(m.Targets, localTarget)
		}

		if s.Go.Enabled() && (s.Go.AbleToBuild() || s.Hephaestus.HTTPEnabled()) {
			localTarget := &formats.MakefileTarget{
				Name: "dev-go",
				Commands: []string{
					"go run . | jq",
				},
			}

			m.Targets = append(m.Targets, localTarget)
		}
	}

	{ // Watch Targets
		parentTarget := &formats.MakefileTarget{
			Comment:       "Watch the project",
			Name:          "watch",
			BeforeTargets: []string{},
		}
		m.Targets = append(m.Targets, parentTarget)

		if s.NPM.Enabled() {

			if s.Go.Enabled() {
				localTarget := &formats.MakefileTarget{
					Name: "watch-npm",
					Commands: []string{
						"npm install",
						"hera frontend",
					},
				}

				m.Targets = append(m.Targets, localTarget)
				parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
			} else if s.NPM.HasScript("watch") {
				localTarget := &formats.MakefileTarget{
					Name: "watch-npm",
					Commands: []string{
						"npm install",
						"npm run watch",
					},
				}

				m.Targets = append(m.Targets, localTarget)
				parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
			}
		}

		if s.Go.Enabled() && s.Go.AbleToBuild() && s.Hephaestus.HTTPEnabled() {

			localTarget := &formats.MakefileTarget{
				Name: "watch-go",
				Commands: []string{
					"hera backend",
				},
			}

			m.Targets = append(m.Targets, localTarget)
			parentTarget.BeforeTargets = append(parentTarget.BeforeTargets, localTarget.Name)
		}

		if len(parentTarget.BeforeTargets) > 0 {
			watchCommand := fmt.Sprintf(
				"make -j%d %s",
				len(parentTarget.BeforeTargets),
				strings.Join(parentTarget.BeforeTargets, " "),
			)

			if s.Go.Enabled() {
				watchCommand = "hera"
			}

			parentTarget.Commands = []string{
				watchCommand,
			}
			parentTarget.BeforeTargets = []string{}
		}
	}

	{ // Clean Target
		if s.Go.Enabled() {
			m.Targets = append(m.Targets, &formats.MakefileTarget{
				Comment: "Run the docs server for the project",
				Name:    "docs-go",
				Commands: []string{
					"@go install golang.org/x/tools/cmd/godoc@latest",
					fmt.Sprintf(`@echo "listening on http://127.0.0.1:6060/pkg/%s"`, s.Go.Path()),
					`@godoc -http=127.0.0.1:6060`,
				},
			})
		}
	}

	file, err := os.Create("Makefile")
	if err != nil {
		return err
	}

	return m.Build(file)
}
