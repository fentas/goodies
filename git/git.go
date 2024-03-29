package git

import (
	"os/exec"
)

func IsIgnored(path string) bool {
	// check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		// nothing to ignore if git is not available
		return true
	}

	cmd := exec.Command("git", "check-ignore", path)
	cmd.Stdout = nil
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
