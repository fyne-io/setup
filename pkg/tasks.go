package pkg

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/tools"
	"golang.org/x/mod/semver"
)

type task struct {
	name, hint string
	test       func() (string, error)
}

var tasks = []*task{
	{"Go compiler", "Checking Go is installed and up to date", func() (string, error) {
		cmd := tools.CommandInShell("go", "version")
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
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = tools.CommandInShell("where", "gcc")
		} else {
			cmd = tools.CommandInShell("which", "gcc")
		}
		_, err := cmd.Output()
		if err == nil {
			return "gcc found", nil
		}

		if runtime.GOOS == "windows" {
			cmd = tools.CommandInShell("where", "clang")
		} else {
			cmd = tools.CommandInShell("which", "clang")
		}
		_, err = cmd.Output()
		if err == nil {
			return "clang found", nil
		}

		return "", errors.New("could not find gcc or clang compilers")
	}},
	{"Go bin PATH", "Verify that the PATH environment is set", func() (string, error) {
		home, _ := os.UserHomeDir()
		goPath := filepath.Join(home, "go", "bin")

		cmd := tools.CommandInShell("go", "env", "GOBIN")
		ret, err := cmd.Output()
		if err != nil {
			return "", err
		}
		bin := strings.TrimSpace(string(ret))
		if bin != "" {
			goPath = bin
		} else {
			cmd = tools.CommandInShell("go", "env", "GOPATH")
			ret, err = cmd.Output()
			if err != nil {
				return "", err
			}
			path := strings.TrimSpace(string(ret))
			if path != "" {
				goPath = filepath.Join(path, "bin")
			}
		}

		isMsys := false
		if runtime.GOOS == "windows" {
			cygpath, err := exec.LookPath("cygpath.exe")

			if err == nil {
				isMsys = true
				goPath = winToUnixPath(goPath, cygpath)
			}

			if isMsys { // use the shell environment
				cmd = exec.Command("env")
			} else {
				cmd = tools.CommandInShell("set")
			}
		} else {
			cmd = tools.CommandInShell("env")
		}
		ret, err = cmd.Output()
		allPath := string(ret)
		if err == nil {
			for _, line := range strings.Split(allPath, "\n") {
				line = strings.TrimSpace(line)
				if len(line) > 5 && (line[:5] == "PATH=" || line[:5] == "Path=") {
					allPath = line
					break
				}
			}
		} else if runtime.GOOS == "windows" {
			cmd = tools.CommandInShell("PowerShell", "-Command", "Get-ChildItem", "Env:Path", "|", "select", "Value", "|", "Format-List")
			ret, err = cmd.Output()
			if err != nil {
				return "", err
			}

			lines := strings.Split(string(ret), "\n")
			allPath = ""
			for _, line := range lines {
				allPath += strings.TrimSpace(line)
			}
		}

		if !strings.Contains(allPath, goPath) {
			return "", errors.New("PATH missing " + goPath)
		}
		return "confirmed", nil
	}},
	//{"Dependencies installed", "Various libraries required are present", func() (string, error) {
	//	return "", errors.New("no dependencies found")
	//}},
	{"Fyne helper", "Checking Fyne tool is installed for packaging", func() (string, error) {
		cmd := tools.CommandInShell("fyne", "--version")
		ret, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(ret)), nil
	}},
}

// winToUnixPath converts a home path from windows to unix using the given cygpath path.
func winToUnixPath(in, convert string) string {
	unixHome, _ := exec.Command(convert, "-u", "~").Output()
	winHome, _ := exec.Command(convert, "-w", "~").Output()
	winRoot, _ := exec.Command(convert, "-w", "/").Output()

	in = strings.ReplaceAll(in, strings.TrimSpace(string(winHome)), strings.TrimSpace(string(unixHome)))
	in = strings.ReplaceAll(in, "\\", "/")

	winRootUnix := strings.ReplaceAll(string(winRoot), "\\", "/")
	return strings.ReplaceAll(in, strings.TrimSpace(winRootUnix), "/")
}
