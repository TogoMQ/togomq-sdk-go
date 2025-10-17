# ✅ Automated Release System Implemented

## Summary

Successfully implemented **automatic release workflow** for TogoMQ SDK using conventional commits. Releases will now be created automatically when code is merged to the `main` branch.

## 📁 Files Created

1. **`.github/workflows/release.yml`** - Automated release workflow
   - Triggers on push to `main` branch
   - Runs tests and linter
   - Analyzes commit messages
   - Creates version tags
   - Generates changelog
   - Creates GitHub releases
   - Comments on PRs

2. **`RELEASE_GUIDE.md`** - Complete quick start guide
   - How the system works
   - Commit message format
   - Examples for all scenarios
   - Troubleshooting guide

3. **`AGENTS.md`** (Updated) - Added comprehensive release documentation
   - Conventional commit format
   - Version bump rules
   - Pre-release checklist
   - Go module usage
   - Troubleshooting

## 🚀 How It Works

### Automatic Flow

```
Code Change → Commit (conventional format) → PR → Merge to main
                                                        ↓
                                                   Workflow runs
                                                        ↓
                                    ┌──────────────────┴──────────────────┐
                                    ↓                                      ↓
                              Run tests & lint                      Analyze commits
                                    ↓                                      ↓
                              ✅ Pass? No → ❌ Fail (no release)          │
                                    ↓ Yes                                 │
                                    └──────────────────┬──────────────────┘
                                                       ↓
                                            Determine version bump
                                                       ↓
                                    ┌─────────────────┼─────────────────┐
                                    ↓                 ↓                 ↓
                              feat: (minor)     fix: (patch)     BREAKING (major)
                              v1.0.0→v1.1.0    v1.0.0→v1.0.1    v1.0.0→v2.0.0
                                    └─────────────────┬─────────────────┘
                                                      ↓
                                              Create Git tag
                                                      ↓
                                            Create GitHub release
                                                      ↓
                                            Comment on PRs
                                                      ↓
                                            ✅ Release complete!
```

## 📝 Commit Message Format

```bash
<type>(<scope>): <subject>

<body>

<footer>
```

### Quick Reference

| Type | Description | Version Bump | Example |
|------|-------------|--------------|---------|
| `feat:` | New feature | Minor (v1.0.0 → v1.1.0) | `feat: add retry mechanism` |
| `fix:` | Bug fix | Patch (v1.0.0 → v1.0.1) | `fix: resolve memory leak` |
| `BREAKING CHANGE:` | Breaking change | Major (v1.0.0 → v2.0.0) | Footer in commit body |
| `docs:` | Documentation | No release | `docs: update README` |
| `test:` | Tests | No release | `test: add unit tests` |
| `chore:` | Maintenance | No release | `chore: update deps` |

## 🎯 Example Workflows

### 1. Adding a New Feature

```bash
# Make changes
git checkout -b add-batch-send
# ... code ...
git commit -m "feat: add batch send optimization"
git push origin add-batch-send

# Create and merge PR
gh pr create --title "Add batch send"
# Merge via GitHub → ✅ Release v1.1.0 created!
```

### 2. Fixing a Bug

```bash
# Fix bug
git checkout -b fix-connection-leak
# ... fix ...
git commit -m "fix: resolve connection leak in client"
git push origin fix-connection-leak

# Merge PR → ✅ Release v1.0.1 created!
```

### 3. Breaking Change

```bash
# Make breaking change
git checkout -b redesign-api
# ... changes ...
git commit -m "feat: redesign subscription API

BREAKING CHANGE: Sub() signature changed. See migration guide."
git push origin redesign-api

# Merge PR → ✅ Release v2.0.0 created!
```

### 4. Documentation Update (No Release)

```bash
# Update docs
git checkout -b update-docs
# ... changes ...
git commit -m "docs: improve installation guide"
git push origin update-docs

# Merge PR → ⏭️ No release created (expected)
```

## ✅ Pre-Commit Checklist

Before merging to `main`:

```bash
# 1. Format code
gofmt -s -w .

# 2. Verify linter configuration (if .golangci.yml was modified)
golangci-lint config verify

# 3. Run linter
golangci-lint run

# 4. Run tests with race detection
go test -race ./...

# 5. Verify commit message format
git log -1 --pretty=%B
# Should follow: <type>: <description>

# 6. Verify no errors in all steps above
```

## 🎬 First Release

To get started, create the initial release:

**Option 1: Manual Tag (Recommended)**
```bash
git checkout main
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0
```

**Option 2: GitHub UI**
1. Go to Releases → "Draft a new release"
2. Tag: `v0.1.0`
3. Title: "Initial Release v0.1.0"
4. Description: List features
5. Publish

After v0.1.0 exists, all future releases are automatic!

## 📦 Installing Released Versions

Users can now install specific versions:

```bash
# Latest version
go get github.com/TogoMQ/togomq-sdk-go@latest

# Specific version
go get github.com/TogoMQ/togomq-sdk-go@v1.2.3

# List all versions
go list -m -versions github.com/TogoMQ/togomq-sdk-go
```

## 🔍 Monitoring Releases

### Check Workflow Status
1. Go to **Actions** tab
2. Look for "Release" workflow
3. View logs if issues occur

### View Releases
1. Go to **Releases** page
2. See all published versions
3. View generated changelogs

### PR Notifications
When released, PRs get a comment:
```
🚀 This PR has been released in v1.2.3
```

## ⚙️ Configuration

The workflow is configured in `.github/workflows/release.yml`:

- **Trigger**: Push to `main` branch
- **Permissions**: `contents: write`, `pull-requests: write`
- **Go Version**: 1.24
- **Default Bump**: `patch` (if no feat/fix in commits)
- **Tag Prefix**: `v`
- **Release Branch**: `main`

## 🐛 Troubleshooting

### No Release Created

**Check:**
1. Commit messages follow conventional format
2. At least one `feat:` or `fix:` commit exists
3. Tests passed (workflow blocks on failure)
4. View Actions tab for errors

### Wrong Version

**Fix:**
- Use `feat:` for minor bumps
- Use `fix:` for patch bumps
- Use `BREAKING CHANGE:` footer for major bumps

### Can't Install Version

**Wait 5-10 minutes** for Go module index to update, then:
```bash
GOPROXY=direct go get github.com/TogoMQ/togomq-sdk-go@v1.2.3
```

## 📚 Documentation

- **Quick Start**: [RELEASE_GUIDE.md](RELEASE_GUIDE.md)
- **Developer Guide**: [AGENTS.md](AGENTS.md) (Release Process section)
- **Conventional Commits**: https://www.conventionalcommits.org/
- **Semantic Versioning**: https://semver.org/

## 🎉 Benefits

✅ **Fully Automated** - No manual tagging or release creation  
✅ **Consistent Versioning** - Follows semver automatically  
✅ **Quality Gates** - Tests must pass before release  
✅ **Auto Changelog** - Generated from commit messages  
✅ **PR Tracking** - Comments on PRs when released  
✅ **Go Module Ready** - Tags work with `go get`  

## 📋 Next Steps

1. ✅ **Create first release** (`v0.1.0`)
2. ✅ **Train team** on conventional commits
3. ✅ **Update PR template** (optional) to remind about commit format
4. ✅ **Test workflow** with a test PR
5. ✅ **Monitor first few releases** to ensure working correctly

## 🚨 Important Notes

1. **Tests are mandatory** - Release is blocked if tests fail
2. **Commit format matters** - Wrong format = no release or wrong version
3. **Breaking changes need footer** - Don't forget `BREAKING CHANGE:` in commit body
4. **First release is manual** - Use `v0.1.0` tag to start
5. **Releases are immediate** - Happen on merge, not on schedule

---

## Status

✅ Release workflow implemented  
✅ Documentation complete  
✅ Ready for first release  

**Your automatic release system is ready to use! Just merge to `main` and get releases! 🚀**
