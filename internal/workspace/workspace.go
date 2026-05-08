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

	for _, ws := range cfg.Workspaces {
		fmt.Printf("Creating workspace: %s\n", ws.Name)
		if err := cmux.CreateWorkspace(ws.Name); err != nil {
			return fmt.Errorf("failed to create workspace %q: %w", ws.Name, err)
		}
		cmux.Wait()

		for i, tab := range ws.Tabs {
			if i == 0 {
				// First tab is created with the workspace; just rename if needed.
				continue
			}
			fmt.Printf("  Creating tab: %s\n", tab.Name)
			if err := cmux.CreateTab(ws.Name, tab.Name); err != nil {
				return fmt.Errorf("failed to create tab %q: %w", tab.Name, err)
			}
			cmux.Wait()
		}

		// Set colors and run commands.
		for _, tab := range ws.Tabs {
			if tab.Color != "" {
				_ = cmux.SetTabColor(ws.Name, tab.Name, tab.Color)
			}

			cmd := buildCommand(tab)
			if cmd != "" {
				dir := expandHome(tab.Dir)
				if dir != "" {
					_ = cmux.SendKeys(ws.Name, tab.Name, fmt.Sprintf("cd %s && %s\n", dir, cmd))
				} else {
					_ = cmux.SendKeys(ws.Name, tab.Name, cmd+"\n")
				}
			}
			cmux.Wait()
		}
	}

	// Focus the first workspace.
	if len(cfg.Workspaces) > 0 {
		_ = cmux.FocusWorkspace(cfg.Workspaces[0].Name)
	}

	fmt.Println("Done! All workspaces launched.")
	return nil
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
