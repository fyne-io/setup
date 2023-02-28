package pkg

import (
	"golang.org/x/sys/execabs"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func getShell() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "zsh"
	}
	out, err := execabs.Command("dscl", ".", "-read", home, "UserShell").Output()
	if err != nil {
		return "zsh"
	}

	items := strings.Split(string(out), ":")
	if len(items) < 2 {
		return "zsh"
	}

	return strings.TrimSpace(items[1])
}

func runCommand(name string, args ...string) *exec.Cmd {
	if runtime.GOOS == "darwin" { // darwin apps don't run in the user shell environment
		return execabs.Command(getShell(), "-i", "-c", name+" "+strings.Join(args, " "))
	}

	return execabs.Command(name, args...)
}
func setupPath() {
	if runtime.GOOS == "darwin" { // darwin apps don't run in the user shell environment
		shellPath, err := execabs.Command("zsh", "-i", "-c", "echo $PATH").Output()
		if err == nil {
			_ = os.Setenv("PATH", string(shellPath))
		}
	}
}
