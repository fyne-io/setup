package pkg

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/semver"
	"golang.org/x/sys/execabs"
)

type task struct {
	name, hint string
	test       func() (string, error)
}

var tasks = []*task{
	{"Go compiler", "Checking Go is installed and up to date", func() (string, error) {
		path, err := execabs.LookPath("go")
		if err != nil {
			return "", err
		}

		cmd := execabs.Command(path, "version")
		ret, err := cmd.Output()
		if err != nil {
			return "", err
		}
		ver := parseGoVersion(strings.TrimSpace(string(ret)))
		before114 := semver.Compare("v"+ver, "v1.14.0") < 0
		if before114 {
			return "go" + ver, errors.New("go version is too old, must be 1.14 or newer")
		}
		return "go" + ver, nil
	}},
	{"C compiler", "Checking a C compiler is installed", func() (string, error) {
		_, err := execabs.LookPath("gcc")
		if err == nil {
			return "gcc found", nil
		}

		_, err = execabs.LookPath("clang")
		if err == nil {
			return "clang found", nil
		}

		return "", errors.New("could not find gcc or clang compilers")
	}},
	{"Go bin PATH", "Verify that the PATH environment is set", func() (string, error) {
		home, _ := os.UserHomeDir()
		goPath := filepath.Join(home, "go", "bin")

		cmd := runCommand("go", "env", "GOBIN")
		ret, err := cmd.Output()
		if err != nil {
			return "", err
		}
		bin := strings.TrimSpace(string(ret))
		if bin != "" {
			goPath = bin
		} else {
			cmd = runCommand("go", "env", "GOPATH")
			ret, err = cmd.Output()
			if err != nil {
				return "", err
			}
			path := strings.TrimSpace(string(ret))
			if path != "" {
				goPath = filepath.Join(path, "bin")
			}
		}

		allPath := os.Getenv("PATH")
		if !strings.Contains(allPath, goPath) {
			return "", errors.New("PATH missing " + goPath)
		}
		return "confirmed", nil
	}},
	//{"Dependencies installed", "Various libraries required are present", func() (string, error) {
	//	return "", errors.New("no dependencies found")
	//}},
	{"Fyne helper", "Checking Fyne tool is installed for packaging", func() (string, error) {
		path, err := execabs.LookPath("fyne")
		if err != nil {
			return "", err
		}

		cmd := runCommand(path, "--version")
		ret, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(ret)), nil
	}},
}
