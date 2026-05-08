# Releasing

## Prerequisites (one-time setup)

1. Create a GitHub Personal Access Token at https://github.com/settings/tokens
   - Select **"Classic"** token
   - Scope: `repo` (full control)
2. Add it as a repository secret:
   - Go to https://github.com/hyoyoungkim-mnetplus/cmux-kiro/settings/secrets/actions
   - Click **"New repository secret"**
   - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
   - Value: your token

## How to Release

```bash
# 1. Make sure you're on main with everything pushed
git checkout main
git pull

# 2. Tag the release
git tag v0.1.0

# 3. Push the tag — this triggers the GitHub Actions workflow
git push origin v0.1.0
```

## What Happens Automatically

1. GitHub Actions runs goreleaser
2. Builds `cmux-kiro` for macOS (amd64 + arm64)
3. Creates a GitHub Release with binaries attached
4. Pushes the Homebrew formula to `homebrew-tap` repo

## Version Bumping

Follow semver:
- `v0.x.y` — pre-1.0 development
- Patch (`v0.1.1`): bug fixes
- Minor (`v0.2.0`): new features
- Major (`v1.0.0`): stable release
