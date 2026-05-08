package main

import (
	"fmt"
	"os"

	"github.com/hyoyoungkim-mnetplus/cmux-kiro/internal/config"
	"github.com/hyoyoungkim-mnetplus/cmux-kiro/internal/workspace"
	"github.com/spf13/cobra"
)

var version = "dev"

var cfgFile string

func main() {
	root := &cobra.Command{
		Use:     "cmux-kiro",
		Short:   "Declarative workspace manager for Kiro CLI + cmux",
		Version: version,
	}

	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: ~/.config/cmux-kiro/workspace.yaml)")

	root.AddCommand(initCmd())
	root.AddCommand(launchCmd())
	root.AddCommand(destroyCmd())
	root.AddCommand(reloadCmd())
	root.AddCommand(validateCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	return config.DefaultConfigPath()
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Generate a default workspace config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := getConfigPath()

			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("config already exists at %s — use --config to specify a different path", path)
			}

			cfg := config.DefaultConfig()
			if err := config.Save(path, cfg); err != nil {
				return err
			}

			fmt.Printf("Config created at %s\n", path)
			fmt.Println("Edit it to define your workspaces, then run: cmux-kiro launch")
			return nil
		},
	}
}

func launchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "launch",
		Short: "Launch all workspaces defined in the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(getConfigPath())
			if err != nil {
				return err
			}
			return workspace.Launch(cfg)
		},
	}
}

func destroyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Close all workspaces defined in the config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(getConfigPath())
			if err != nil {
				return err
			}
			return workspace.Destroy(cfg)
		},
	}
}

func reloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reload",
		Short: "Destroy and re-launch all workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(getConfigPath())
			if err != nil {
				return err
			}
			_ = workspace.Destroy(cfg)
			return workspace.Launch(cfg)
		},
	}
}

func validateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the config file without launching",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(getConfigPath())
			if err != nil {
				return err
			}

			total := 0
			for _, ws := range cfg.Workspaces {
				total += len(ws.Tabs)
			}

			fmt.Printf("Config OK: %d workspace(s), %d tab(s) total\n", len(cfg.Workspaces), total)
			return nil
		},
	}
}
