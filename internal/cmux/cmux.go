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

// ParseRef extracts a specific ref (e.g. "workspace:14") from cmux output like "OK workspace:14".
func ParseRef(output, prefix string) string {
	for _, part := range strings.Fields(output) {
		if strings.HasPrefix(part, prefix+":") {
			return part
		}
	}
	return ""
}

// Ping checks if cmux is running.
func Ping() bool {
	_, err := Run("ping")
	return err == nil
}

// NewWorkspace creates a new cmux workspace. Returns the workspace ref.
func NewWorkspace(cwd string) (string, error) {
	args := []string{"new-workspace"}
	if cwd != "" {
		args = append(args, "--cwd", cwd)
	}
	out, err := Run(args...)
	if err != nil {
		return "", err
	}
	ref := ParseRef(out, "workspace")
	if ref == "" {
		return "", fmt.Errorf("could not parse workspace ref from: %s", out)
	}
	return ref, nil
}

// RenameWorkspace renames a workspace.
func RenameWorkspace(wsRef string, title string) error {
	_, err := Run("rename-workspace", "--workspace", wsRef, title)
	return err
}

// NewSurface creates a new tab (surface) in a workspace. Returns the surface ref.
func NewSurface(wsRef string) (string, error) {
	out, err := Run("new-surface", "--workspace", wsRef)
	if err != nil {
		return "", err
	}
	ref := ParseRef(out, "surface")
	if ref == "" {
		return "", fmt.Errorf("could not parse surface ref from: %s", out)
	}
	return ref, nil
}

// RenameTab renames a tab/surface.
func RenameTab(wsRef string, surfRef string, title string) error {
	_, err := Run("rename-tab", "--workspace", wsRef, "--surface", surfRef, title)
	return err
}

// Send sends text to a surface.
func Send(wsRef string, surfRef string, text string) error {
	args := []string{"send", "--workspace", wsRef}
	if surfRef != "" {
		args = append(args, "--surface", surfRef)
	}
	args = append(args, text)
	_, err := Run(args...)
	return err
}

// SendCommand sends text and presses Return to execute it.
func SendCommand(wsRef string, surfRef string, command string) error {
	if err := Send(wsRef, surfRef, command); err != nil {
		return err
	}
	args := []string{"send-key", "--workspace", wsRef}
	if surfRef != "" {
		args = append(args, "--surface", surfRef)
	}
	args = append(args, "Return")
	_, err := Run(args...)
	return err
}

// SelectWorkspace focuses a workspace.
func SelectWorkspace(wsRef string) error {
	_, err := Run("select-workspace", "--workspace", wsRef)
	return err
}

// SetStatus sets a colored status indicator on a workspace sidebar.
func SetStatus(wsRef string, label string, color string, icon string) error {
	args := []string{"set-status", "project", label, "--color", color, "--workspace", wsRef}
	if icon != "" {
		args = append(args, "--icon", icon)
	}
	_, err := Run(args...)
	return err
}

// ListWorkspaces returns workspace list.
func ListWorkspaces() (string, error) {
	return Run("list-workspaces")
}

// Wait pauses briefly between cmux commands.
func Wait() {
	time.Sleep(300 * time.Millisecond)
}

// WaitForShell waits for a new shell to be ready by polling read-screen.
func WaitForShell(wsRef, surfRef string) {
	for i := 0; i < 10; i++ {
		time.Sleep(300 * time.Millisecond)
		out, err := Run("read-screen", "--workspace", wsRef, "--surface", surfRef, "--lines", "3")
		if err == nil && (strings.Contains(out, ">") || strings.Contains(out, "$") || strings.Contains(out, "%") || strings.Contains(out, "#")) {
			return
		}
	}
}
