package workspace

import (
	"fmt"
	"os"
	"strings"

	"github.com/hyoyoungkim-mnetplus/cmux-kiro/internal/cmux"
	"github.com/hyoyoungkim-mnetplus/cmux-kiro/internal/config"
)

// Launch creates workspaces and tabs from the config, then runs commands.
func Launch(cfg *config.Config) error {
	if !cmux.Ping() {
		return fmt.Errorf("cmux is not running — please start cmux first")
	}

	var firstWsRef string

	for _, ws := range cfg.Workspaces {
		fmt.Printf("Creating workspace: %s\n", ws.Name)

		// First tab: use new-workspace --command to launch directly.
		firstTab := ws.Tabs[0]
		firstDir := expandHome(firstTab.Dir)
		firstCmd := buildCommand(firstTab)

		wsRef, err := cmux.NewWorkspace(firstDir, firstCmd)
		if err != nil {
			return fmt.Errorf("failed to create workspace %q: %w", ws.Name, err)
		}

		if firstWsRef == "" {
			firstWsRef = wsRef
		}

		// Rename workspace and set sidebar color.
		_ = cmux.RenameWorkspace(wsRef, ws.Name)
		if ws.Color != "" {
			label := ws.Label
			if label == "" {
				label = ws.Name
			}
			_ = cmux.SetStatus(wsRef, label, ws.Color, ws.Icon)
		}

		// Rename first tab.
		firstSurfRef := getFirstSurface(wsRef)
		if firstSurfRef != "" {
			_ = cmux.RenameTab(wsRef, firstSurfRef, firstTab.Name)
		}
		fmt.Printf("  Tab: %s\n", firstTab.Name)

		// Additional tabs: new-surface + respawn-pane.
		for j := 1; j < len(ws.Tabs); j++ {
			tab := ws.Tabs[j]
			surfRef, err := cmux.NewSurface(wsRef)
			if err != nil {
				fmt.Printf("  Warning: could not create tab %q: %v\n", tab.Name, err)
				continue
			}

			_ = cmux.RenameTab(wsRef, surfRef, tab.Name)

			// Build full command with cd if needed.
			dir := expandHome(tab.Dir)
			cmd := buildCommand(tab)
			fullCmd := ""
			if dir != "" && cmd != "" {
				fullCmd = fmt.Sprintf("cd %s && %s", dir, cmd)
			} else if dir != "" {
				fullCmd = fmt.Sprintf("cd %s", dir)
			} else if cmd != "" {
				fullCmd = cmd
			}

			if fullCmd != "" {
				_ = cmux.RespawnPane(wsRef, surfRef, fullCmd)
			}

			fmt.Printf("  Tab: %s\n", tab.Name)
		}
	}

	// Focus first workspace.
	if firstWsRef != "" {
		_ = cmux.SelectWorkspace(firstWsRef)
	}

	fmt.Println("Done!")
	return nil
}

// getFirstSurface retrieves the first surface ref of a workspace.
func getFirstSurface(wsRef string) string {
	out, err := cmux.Run("list-pane-surfaces", "--workspace", wsRef)
	if err != nil {
		return ""
	}
	return cmux.ParseRef(out, "surface")
}

// buildCommand constructs the shell command for a tab.
func buildCommand(tab config.Tab) string {
	if tab.Command != "" {
		if tab.Agent != "" {
			return fmt.Sprintf("%s --agent %s", tab.Command, tab.Agent)
		}
		return tab.Command
	}
	if tab.Agent != "" {
		return fmt.Sprintf("kiro-cli chat --agent %s", tab.Agent)
	}
	return ""
}

// expandHome replaces ~ with the actual home directory.
func expandHome(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return home + path[1:]
	}
	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
	return path
}
