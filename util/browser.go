package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Open(url string) error {
	var cmd *exec.Cmd
	switch {
	case runtime.GOOS == "linux":
		// should exist on most linux environments
		// TODO: maybe check
		cmd = exec.Command("xdg-open", url)
	case runtime.GOOS == "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case runtime.GOOS == "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	fmt.Printf("Trying to open: %s in web browser\n\n", url)
	return cmd.Start()
}
