package state

import (
	"io"
	"os"
	"strings"

	"github.com/lunagic/hephaestus/internal/utils"
	"golang.org/x/mod/modfile"
)

func NewGo() (*Go, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return &Go{}, nil
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	goMod, err := modfile.Parse("go.mod", fileBytes, nil)
	if err != nil {
		return nil, err
	}

	return &Go{
		modfile: goMod,
	}, nil
}

type Go struct {
	modfile *modfile.File
}

func (g *Go) Enabled() bool {
	return g.modfile != nil
}

func (g *Go) InstallTool(path string) {
	parts := strings.Split(path, "@")
	for _, tool := range g.modfile.Tool {
		if tool.Path == parts[0] {
			return
		}
	}
	utils.ShellExec("go", "get", "-tool", path)
}

func (s *Go) AbleToBuild() bool {
	_, err := os.Open("main.go")

	return err == nil
}

func (s *Go) Path() string {
	if s.modfile == nil {
		return ""
	}

	return s.modfile.Module.Mod.Path
}

func (s *Go) Version() string {
	if s.modfile == nil {
		return ""
	}

	shortVersion := strings.Join(strings.Split(s.modfile.Go.Version, ".")[0:2], ".")

	return shortVersion
}

func (s *Go) HasPkg(pkg string) bool {
	for _, require := range s.modfile.Require {
		if require.Mod.Path == pkg {
			return true
		}
	}

	return false
}
