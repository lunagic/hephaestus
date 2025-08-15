package generators

import (
	"fmt"
	"os"

	"github.com/lunagic/hephaestus/internal/formats"
	"github.com/lunagic/hephaestus/internal/state"
)

type Dockerfile struct{}

func (generator Dockerfile) Output(s *state.State) error {
	dockerfile := formats.Dockerfile{}

	frontendCopyCommands := []string{}

	if s.NPM.Enabled() {
		dockerfile.Stages = append(dockerfile.Stages, &formats.DockerStage{
			Name:  "frontend_builder",
			Image: "node",
			Tag:   fmt.Sprintf("%s-alpine", s.NPM.Version()),
			Commands: []string{
				"WORKDIR /workspace",
				"COPY . .",
				"RUN npm install",
				"RUN npm run build",
			},
		})
	}

	if s.Hephaestus.StaticSitePath != "" {
		dockerfile.Stages = append(dockerfile.Stages, &formats.DockerStage{
			Image: "ghcr.io/lunagic/poseidon",
			Tag:   "latest",
			Commands: []string{
				fmt.Sprintf("COPY --from=frontend_builder /workspace/%s .", s.Hephaestus.StaticSitePath),
			},
		})
	}

	if s.Go.Enabled() && s.Hephaestus.HTTPEnabled() {
		buildCommands := []string{
			"RUN apk add --no-cache git gcc g++",
			"WORKDIR /workspace",
			"COPY . .",
			"RUN git clean -Xdff",
		}

		buildCommands = append(buildCommands, frontendCopyCommands...)

		buildCommands = append(buildCommands,
			"RUN CGO_ENABLED=1 go build -ldflags='-s -w' -o /usr/local/bin/build .",
		)

		dockerfile.Stages = append(dockerfile.Stages, &formats.DockerStage{
			Name:     "backend_builder",
			Image:    "golang",
			Tag:      fmt.Sprintf("%s-alpine", s.Go.Version()),
			Commands: buildCommands,
		})

		dockerfile.Stages = append(dockerfile.Stages, &formats.DockerStage{
			Image: "alpine",
			Tag:   "latest",
			Commands: []string{
				"WORKDIR /workspace",
				"COPY --from=backend_builder /usr/local/bin/build /usr/local/bin/build",
				`CMD [ "build" ]`,
				"ENV HOST=0.0.0.0",
				fmt.Sprintf("ENV PORT=%d", s.Hephaestus.DefaultPort),
				fmt.Sprintf("EXPOSE %d", s.Hephaestus.DefaultPort),
			},
		})
	}

	if len(dockerfile.Stages) == 0 {
		return nil
	}

	file, err := os.Create("Dockerfile")
	if err != nil {
		return err
	}

	return dockerfile.Build(file)
}
