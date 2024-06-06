package utils

import "os/exec"

func Refresh() {
	// TODO: Find a better way to send keyboard events,
	//       without relying on external commands or need
	//       for root access
	cmd := exec.Command("xdotool", "key", "F5")
	cmd.Run()
}
