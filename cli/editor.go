package cli

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const defaultEditor = "vim"

func getPreferredEditorFromEnvironment() string {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		return defaultEditor
	}

	return editor
}

// resolveEditorArguments for now only with VS-Code, others can be added if needed
func resolveEditorArguments(executable string, filename string) []string {
	args := []string{filename}

	if strings.Contains(executable, "code") {
		args = append([]string{"--wait"}, args...)
	}

	return args
}

func openFileInEditor(filename string, editor string) error {
	// Get the full executable path for the editor.
	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, resolveEditorArguments(executable, filename)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Edit opens a temporary file in a text editor, write data into it for editing and
// returns the written bytes on success or an error on failure. It handles deletion
// of the temporary file behind the scenes.
// nolint: gosec
func Edit(data []byte) ([]byte, error) {
	resolveEditor := getPreferredEditorFromEnvironment()

	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return []byte{}, err
	}

	filename := file.Name()

	_, err = file.Write(data)
	if err != nil {
		return []byte{}, err
	}

	// Defer removal of the temporary file in case any of the next steps fail.
	// nolint: errcheck
	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = openFileInEditor(filename, resolveEditor); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
