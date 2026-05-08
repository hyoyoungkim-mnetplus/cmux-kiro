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

		// Create workspace with first tab's directory.
		var firstDir string
		if len(ws.Tabs) > 0 {
			firstDir = expandHome(ws.Tabs[0].Dir)
		}
		wsRef, err := cmux.NewWorkspace(firstDir)
		if err != nil {
			return fmt.Errorf("failed to create workspace %q: %w", ws.Name, err)
		}
		cmux.Wait()

		if firstWsRef == "" {
			firstWsRef = wsRef
		}

		// Rename workspace.
		_ = cmux.RenameWorkspace(wsRef, ws.Name)
		if ws.Color != "" {
			label := ws.Label
			if label == "" {
				label = ws.Name
			}
			_ = cmux.SetStatus(wsRef, label, ws.Color, ws.Icon)
		}
		cmux.Wait()

		// Track surface refs. First surface comes with the workspace.
		var surfaceRefs []string

		// Get first surface ref from tree or create logic.
		// The first surface is auto-created; we need its ref.
		// Use new-surface output pattern: first surface isn't returned by new-workspace.
		// We'll list pane-surfaces to find it.
		firstSurfRef := getFirstSurface(wsRef)
		surfaceRefs = append(surfaceRefs, firstSurfRef)

		// Handle first tab.
		if len(ws.Tabs) > 0 {
			tab := ws.Tabs[0]
			fmt.Printf("  Tab: %s\n", tab.Name)
			if firstSurfRef != "" {
				_ = cmux.RenameTab(wsRef, firstSurfRef, tab.Name)
				cmd := buildCommand(tab)
				if cmd != "" {
					cmux.WaitForShell(wsRef, firstSurfRef)
					_ = cmux.SendCommand(wsRef, firstSurfRef, cmd)
				}
			}
			cmux.Wait()
		}

		// Create additional tabs.
		for j := 1; j < len(ws.Tabs); j++ {
			tab := ws.Tabs[j]
			fmt.Printf("  Tab: %s\n", tab.Name)

			surfRef, err := cmux.NewSurface(wsRef)
			if err != nil {
				fmt.Printf("  Warning: could not create tab: %v\n", err)
				continue
			}
			surfaceRefs = append(surfaceRefs, surfRef)
			cmux.WaitForShell(wsRef, surfRef)

			_ = cmux.RenameTab(wsRef, surfRef, tab.Name)

			// cd + command
			dir := expandHome(tab.Dir)
			cmd := buildCommand(tab)
			if dir != "" && cmd != "" {
				_ = cmux.SendCommand(wsRef, surfRef, fmt.Sprintf("cd %s && %s", dir, cmd))
			} else if dir != "" {
				_ = cmux.SendCommand(wsRef, surfRef, fmt.Sprintf("cd %s", dir))
			} else if cmd != "" {
				_ = cmux.SendCommand(wsRef, surfRef, cmd)
			}
			cmux.Wait()
		}
	}

	// Focus first workspace.
	if firstWsRef != "" {
		_ = cmux.SelectWorkspace(firstWsRef)
	}

	fmt.Println("Done! All workspaces launched.")
	return nil
}

// getFirstSurface retrieves the first surface ref of a workspace.
func getFirstSurface(wsRef string) string {
	out, err := cmux.Run("list-pane-surfaces", "--workspace", wsRef)
	if err != nil {
		return ""
	}
	// Output contains surface refs like "surface:21"
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
