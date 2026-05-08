# cmux-kiro

> Declarative workspace manager for [Kiro CLI](https://kiro.dev) + [cmux](https://github.com/manaflow-ai/cmux).

Define your agent workspaces in YAML. Launch them all with one command.  
Never lose your layout again.

## The Problem

You close cmux. Everything's gone — tabs, names, colors, agents. You set it all up again. Every. Single. Time.

## The Solution

```yaml
# ~/.config/cmux-kiro/workspace.yaml
workspaces:
  - name: devops
    color: "#4CAF50"
    tabs:
      - name: infra
        dir: ~/projects/terraform
        command: kiro-cli chat
      - name: query
        dir: ~/
        command: kiro-cli chat
        agent: query
      - name: jira
        dir: ~/
        command: kiro-cli chat
        agent: jira
```

```bash
cmux-kiro launch
```

Done. Your entire workspace — tabs, colors, directories, agents, dev tools — materialized in one shot.

Works with any command: `kiro-cli`, `k9s`, `lazygit`, `stern`, `htop`, or just a plain shell.

## Install

**Homebrew:**

```bash
brew tap hyoyoungkim-mnetplus/tap
brew install cmux-kiro
```

**From source:**

```bash
go install github.com/hyoyoungkim-mnetplus/cmux-kiro/cmd/cmux-kiro@latest
```

## Usage

```bash
# Generate a starter config
cmux-kiro init

# Edit your workspace definition
$EDITOR ~/.config/cmux-kiro/workspace.yaml

# Launch everything
cmux-kiro launch

# Validate config without launching
cmux-kiro validate
```

## Examples

**DevOps multi-tool setup:**

```yaml
workspaces:
  - name: devops
    color: "#4CAF50"
    tabs:
      - name: infra
        dir: ~/projects/terraform
        command: kiro-cli chat
      - name: k9s-dev
        command: k9s --context dev
      - name: k9s-prod
        command: k9s --context prod-legacy
      - name: logs
        command: stern -n default --context dev
      - name: shell
        dir: ~/projects
```

**Multi-agent workspace:**

```yaml
workspaces:
  - name: agents
    tabs:
      - name: code
        dir: ~/projects/my-app
        command: kiro-cli chat
        agent: developer
      - name: review
        dir: ~/projects/my-app
        agent: reviewer
      - name: jira
        agent: jira
```

## Config Reference

```yaml
workspaces:
  - name: string        # Workspace name (shown in cmux sidebar)
    color: string       # Hex color for the workspace (optional)
    tabs:
      - name: string    # Tab name
        dir: string     # Working directory (supports ~)
        command: string # Shell command to run (optional)
        agent: string   # Kiro CLI agent name (optional, appends --agent flag)
        color: string   # Tab color override (optional)
```

### Agent shorthand

If you set `agent` without `command`, it defaults to `kiro-cli chat --agent <name>`:

```yaml
tabs:
  - name: jira
    agent: jira
    # equivalent to: command: "kiro-cli chat --agent jira"
```

## Requirements

- [cmux](https://github.com/manaflow-ai/cmux) (must be running)
- [Kiro CLI](https://kiro.dev) (for agent features)

## How It Works

1. Reads your declarative YAML config
2. Creates cmux workspaces and tabs via the `cmux` CLI
3. Sets colors and sends commands to each tab
4. Focuses the first workspace

No daemons. No background processes. Just a single binary that talks to cmux.

## License

MIT
