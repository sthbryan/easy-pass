package clipboard

import (
	"os/exec"
	"runtime"
	"strings"
)

type Clipboard interface {
	Copy(text string) error
}

type SystemClipboard struct{}

func NewSystemClipboard() *SystemClipboard {
	return &SystemClipboard{}
}

func (c *SystemClipboard) Copy(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("cmd", "/c", "echo", text, "|", "clip")
	default:
		return &ClipboardError{Message: "Unsupported OS"}
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

type ClipboardError struct {
	Message string
}

func (e *ClipboardError) Error() string {
	return e.Message
}
