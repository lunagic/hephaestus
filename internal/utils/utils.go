package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
)

func ShellExec(name string, arg ...string) error {
	log.Println("ShellExec:", name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func ShellOutput(name string, arg ...string) string {
	x := bytes.NewBufferString("")
	cmd := exec.Command(name, arg...)
	cmd.Stdout = x
	cmd.Stderr = x
	if err := cmd.Run(); err != nil {
		return err.Error()
	}

	return x.String()
}

func ToJSON(v any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	err := WriteJSON(v, buffer)
	return buffer.Bytes(), err
}

func WriteJSON(v any, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	return encoder.Encode(v)
}
