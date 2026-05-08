package cmux

import (
	"fmt"
	"os/exec"
	"strings"
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

// ParseRef extracts a specific ref (e.g. "workspace:14") from cmux output.
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
func NewWorkspace(cwd string, command string) (string, error) {
	args := []string{"new-workspace"}
	if cwd != "" {
		args = append(args, "--cwd", cwd)
	}
	if command != "" {
		args = append(args, "--command", command)
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

// RespawnPane sends a command to a surface for execution.
func RespawnPane(wsRef string, surfRef string, command string) error {
	_, err := Run("respawn-pane", "--workspace", wsRef, "--surface", surfRef, "--command", command)
	return err
}

// SelectWorkspace focuses a workspace.
func SelectWorkspace(wsRef string) error {
	_, err := Run("select-workspace", "--workspace", wsRef)
	return err
}

// CloseWorkspace closes a workspace.
func CloseWorkspace(wsRef string) error {
	_, err := Run("close-workspace", "--workspace", wsRef)
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
