package formats

import (
	"fmt"
	"io"
	"strings"
)

type Makefile struct {
	Variables []string
	Targets   []*MakefileTarget
}

type MakefileTarget struct {
	Comment       string
	Name          string
	BeforeTargets []string
	Commands      []string
}

func (m Makefile) Build(w io.Writer) error {
	// PHONY
	allTargets := []string{}
	for _, target := range m.Targets {
		allTargets = append(allTargets, target.Name)
	}

	_, _ = w.Write(fmt.Appendf(nil, ".PHONY: %s\n\n", strings.Join(allTargets, " ")))

	// Variables
	for _, variable := range m.Variables {
		_, _ = w.Write(fmt.Appendf(nil, "%s\n", variable))
	}

	_, _ = w.Write([]byte("\n"))

	// Targets
	for i, target := range m.Targets {
		// Comment line
		if target.Comment != "" {
			_, _ = w.Write(fmt.Appendf(nil, "## %s\n", target.Comment))
		}

		// Target Line
		beforeTargets := ""
		if len(target.BeforeTargets) > 0 {
			beforeTargets = fmt.Sprintf(" %s", strings.Join(target.BeforeTargets, " "))
		}

		_, _ = w.Write(fmt.Appendf(nil, "%s:%s\n", target.Name, beforeTargets))

		// Command lines
		for _, command := range target.Commands {
			_, _ = w.Write(fmt.Appendf(nil, "\t%s\n", command))
		}

		// Not on the last line
		if i < len(m.Targets)-1 {
			_, _ = w.Write([]byte("\n"))
		}
	}

	return nil
}
