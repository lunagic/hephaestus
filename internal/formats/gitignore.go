package formats

import (
	"fmt"
	"io"
	"strings"
)

type GitIgnore struct {
	Sections []*GitIgnoreSection
}

type GitIgnoreSection struct {
	Title string
	Items []string
}

func (b GitIgnore) Build(w io.Writer) error {
	for i, section := range b.Sections {
		if len(section.Items) < 1 {
			continue
		}

		if i != 0 {
			_, _ = w.Write([]byte("\n"))
		}

		_, _ = w.Write([]byte(fmt.Sprintf("# %s\n", section.Title)))
		_, _ = w.Write([]byte(strings.Join(section.Items, "\n")))
		_, _ = w.Write([]byte("\n"))
	}

	return nil
}
