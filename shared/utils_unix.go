// +build darwin linux

package shared

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// EnsurePython3Venv ensures python3 venv for hashicorp python plugin
// venvDir should be directory path of target venv
func EnsurePython3Venv(venvDir string, packages ...string) (python3 string, err error) {
	python3 = filepath.Join(venvDir, "bin", "python3")

	log.Info().
		Str("python3", python3).
		Strs("packages", packages).
		Msg("ensure python3 venv")

	// check if python3 venv is available
	if err := exec.Command(python3, "--version").Run(); err != nil {
		// python3 venv not available, create one
		// check if system python3 is available
		if err := execCommand("python3", "--version"); err != nil {
			return "", errors.Wrap(err, "python3 not found")
		}

		// check if .venv exists
		if _, err := os.Stat(venvDir); err == nil {
			// .venv exists, remove first
			if err := execCommand("rm", "-rf", venvDir); err != nil {
				return "", errors.Wrap(err, "remove existed venv failed")
			}
		}

		// create python3 .venv
		if err := execCommand("python3", "-m", "venv", venvDir); err != nil {
			return "", errors.Wrap(err, "create python3 venv failed")
		}
	}

	for _, pkg := range packages {
		err := InstallPythonPackage(python3, pkg)
		if err != nil {
			return python3, errors.Wrap(err, fmt.Sprintf("pip install %s failed", pkg))
		}
	}

	return python3, nil
}

func execCommand(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	log.Info().Str("cmd", cmd.String()).Msg("exec command")

	// print output with colors
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("exec command failed")
		return err
	}

	return nil
}