package cmux

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Run executes a cmux CLI command and returns the output.
func Run(args ...string) (string, error) {
	cmd := exec.Command("cmux", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cmux %s: %s", strings.Join(args, " "), string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// Ping checks if cmux is running.
func Ping() bool {
	_, err := Run("ping")
	return err == nil
}

// CreateWorkspace creates a new cmux workspace.
func CreateWorkspace(name string) error {
	_, err := Run("workspace", "new", name)
	return err
}

// CreateTab creates a new tab in the specified workspace.
func CreateTab(workspace, name string) error {
	_, err := Run("tab", "new", "--workspace", workspace, "--name", name)
	return err
}

// SetTabColor sets the color of a tab.
func SetTabColor(workspace, tab, color string) error {
	_, err := Run("tab", "color", "--workspace", workspace, "--name", tab, "--color", color)
	return err
}

// SendKeys sends keystrokes to a specific tab.
func SendKeys(workspace, tab, keys string) error {
	_, err := Run("send-keys", "--workspace", workspace, "--tab", tab, keys)
	return err
}

// ListWorkspaces returns the list of current workspaces.
func ListWorkspaces() ([]string, error) {
	out, err := Run("workspace", "list")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// FocusWorkspace switches focus to the specified workspace.
func FocusWorkspace(name string) error {
	_, err := Run("workspace", "focus", name)
	return err
}

// Wait pauses briefly between cmux commands.
func Wait() {
	time.Sleep(200 * time.Millisecond)
}
