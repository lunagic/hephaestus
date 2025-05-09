package formats

import (
	"encoding/json"
	"os"

	"github.com/lunagic/hephaestus/internal/utils"
)

type PackageJSON map[string]any

func (p PackageJSON) ReadFromDisk() error {
	file, err := os.Open("package.json")
	if err != nil {
		return err
	}

	// Reset the state
	for key := range p {
		delete(p, key)
	}

	if err := json.NewDecoder(file).Decode(&p); err != nil {
		return err
	}

	return nil
}

func (p PackageJSON) WriteToDisk() error {
	file, err := os.Create("package.json")
	if err != nil {
		return err
	}

	return utils.WriteJSON(p, file)
}

func (p PackageJSON) GetScript(name string) string {
	scripts, _ := p["scripts"].(map[string]any)

	scriptString, ok := scripts[name].(string)
	if !ok {
		return ""
	}

	return scriptString
}

func (p PackageJSON) SetScript(name string, script string) {
	scripts, _ := p["scripts"].(map[string]any)
	scripts[name] = script

	_ = p.WriteToDisk()
	_ = p.ReadFromDisk()
}

func (p PackageJSON) HasPackage(packageName string) bool {

	dependencies, _ := p["dependencies"].(map[string]any)
	if _, found := dependencies[packageName]; found {
		return true
	}

	peerDependencies, _ := p["peerDependencies"].(map[string]any)
	if _, found := peerDependencies[packageName]; found {
		return true
	}

	devDependencies, _ := p["devDependencies"].(map[string]any)
	if _, found := devDependencies[packageName]; found {
		return true
	}

	return false
}

func (p PackageJSON) ConfirmPackage(packageName string, dev bool, saveExact bool) error {
	if !dev {
		dependencies, _ := p["dependencies"].(map[string]any)
		if _, found := dependencies[packageName]; found {
			return nil
		}
	}

	if dev {
		devDependencies, _ := p["devDependencies"].(map[string]any)
		if _, found := devDependencies[packageName]; found {
			return nil
		}
	}

	args := []string{
		"install",
	}

	if dev {
		args = append(args, "--save-dev")
	}

	if saveExact {
		args = append(args, "--save-exact")
	}

	args = append(args, packageName)

	if err := utils.ShellExec("npm", args...); err != nil {
		return err
	}

	return p.ReadFromDisk()
}
