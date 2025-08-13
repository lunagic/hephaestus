package formats

import (
	"fmt"
	"io"
)

type Dockerfile struct {
	Stages []*DockerStage
}

type DockerStage struct {
	Comment  string
	Image    string
	Tag      string
	Name     string
	Commands []string
}

func (b Dockerfile) Build(w io.Writer) error {
	for i, stage := range b.Stages {
		if stage.Comment != "" {
			if _, err := w.Write(fmt.Appendf(nil, "## %s\n", stage.Comment)); err != nil {
				return err
			}
		}

		nameOutput := ""
		if stage.Name != "" {
			nameOutput = fmt.Sprintf(" AS %s", stage.Name)
		}
		if _, err := fmt.Fprintf(w,
			"FROM %s:%s%s\n", stage.Image, stage.Tag, nameOutput,
		); err != nil {
			return err
		}

		for _, command := range stage.Commands {
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
