package formats

import (
	"fmt"
	"io"
)

type Dockerfile struct {
	Stages []*DockerStage
}

type DockerStage struct {
	Comment     string
	Commands    []string
	CopyCommand string
}

func (b Dockerfile) Build(w io.Writer) error {
	for i, stage := range b.Stages {
		if stage.Comment != "" {
			if _, err := w.Write(fmt.Appendf(nil, "## %s\n", stage.Comment)); err != nil {
				return err
			}
		}

		for _, command := range stage.Commands {
			if command == "## COPY_PLACEHOLDER" {
				if stage.CopyCommand == "" {
					continue
				}

				command = stage.CopyCommand
			}

			if _, err := w.Write(fmt.Appendf(nil, "%s\n", command)); err != nil {
				return err
			}
		}

		if i < len(b.Stages)-1 {
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}

	return nil
}
