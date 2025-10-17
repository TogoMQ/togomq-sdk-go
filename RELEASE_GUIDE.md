# Automated Release Setup - Quick Start

## ✅ What Was Implemented

Automatic releases using **conventional commits** - the simplest approach for automated versioning.

## 🚀 How It Works

1. **Merge to `main`** → Workflow triggers
2. **Analyze commits** → Determines version bump
3. **Run tests & lint** → Ensures quality
4. **Create release** → Tag, changelog, and GitHub release
5. **Comment on PRs** → Notifies contributors

## 📝 Commit Message Format

Use conventional commit messages to control versioning:

```bash
<type>(<scope>): <subject>

<body>

<footer>
```

### Version Bump Rules

| Commit Type | Example | Version Bump |
|-------------|---------|--------------|
| `feat:` or `feature:` | `feat: add retry mechanism` | v1.0.0 → v1.1.0 (minor) |
| `fix:` or `bugfix:` | `fix: correct race condition` | v1.0.0 → v1.0.1 (patch) |
| `BREAKING CHANGE:` | See below | v1.0.0 → v2.0.0 (major) |
| `docs:`, `test:`, `chore:` | No version bump | No release |

### Examples

#### Minor Version Bump (New Feature)
```bash
git commit -m "feat: add retry mechanism for failed publishes"
# Result: v1.0.0 → v1.1.0
```

#### Patch Version Bump (Bug Fix)
```bash
git commit -m "fix: correct race condition in subscription handler"
# Result: v1.0.0 → v1.0.1
```

#### Major Version Bump (Breaking Change)
```bash
git commit -m "feat: redesign client API

BREAKING CHANGE: removed WithQueue method, use WithTopic instead"
# Result: v1.0.0 → v2.0.0
```

#### No Release
```bash
git commit -m "docs: update installation instructions"
git commit -m "test: add integration tests"
git commit -m "chore: update dependencies"
# Result: No version bump, no release
```

## 🔄 Workflow

### Normal Development Flow

```bash
# 1. Create feature branch
git checkout -b add-retry-logic

# 2. Make changes
# ... code changes ...

# 3. Commit with conventional format
git commit -m "feat: add retry mechanism for failed publishes"

# 4. Push and create PR
git push origin add-retry-logic
gh pr create --title "Add retry mechanism" --body "Implements automatic retry"

# 5. Merge PR to main (via GitHub UI)
# ✅ Release v1.1.0 created automatically!
```

### Creating a Hotfix

```bash
# 1. Create fix branch
git checkout -b fix-memory-leak

# 2. Fix the issue
# ... fix code ...

# 3. Commit with fix type
git commit -m "fix: resolve memory leak in subscription goroutine"

# 4. Push and merge
git push origin fix-memory-leak
gh pr create --title "Fix memory leak"
# Merge via GitHub UI

# ✅ Release v1.0.1 created automatically!
```

### Breaking Change Release

```bash
# 1. Make breaking changes
git checkout -b redesign-api

# 2. Commit with BREAKING CHANGE footer
git commit -m "feat: redesign subscription API

BREAKING CHANGE: Sub() now returns separate channels for messages and errors
instead of a single channel. Update your code:

Before:
  msgChan, err := client.Sub(ctx, opts)

After:
  msgChan, errChan, err := client.Sub(ctx, opts)"

# 3. Push and merge
git push origin redesign-api
gh pr create --title "Redesign subscription API"
# Merge via GitHub UI

# ✅ Release v2.0.0 created automatically!
```

## 📋 Pre-Release Checklist

Before merging to `main`:

- [ ] All tests pass (`go test -race ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Code formatted (`gofmt -s -w .`)
- [ ] README updated if needed
- [ ] Examples tested
- [ ] Commit messages follow conventional format
- [ ] Breaking changes documented in commit message

## 🎯 First Release

To create the first release (v0.1.0):

**Option 1: Manual Tag (Recommended for first release)**
```bash
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0
```

**Option 2: Through GitHub UI**
1. Go to Releases → Draft a new release
2. Create new tag: `v0.1.0`
3. Title: "Initial Release v0.1.0"
4. Publish release

After the first tag exists, all subsequent releases will be automatic!

## 📦 Using Released Versions

After a release is created, users can install it:

```bash
# Install specific version
go get github.com/TogoMQ/togomq-sdk-go@v1.2.3

# Install latest version
go get github.com/TogoMQ/togomq-sdk-go@latest

# List all available versions
go list -m -versions github.com/TogoMQ/togomq-sdk-go
```

## 🔍 Monitoring Releases

### Check Release Status

1. Go to **Actions** tab in GitHub
2. Look for "Release" workflow
3. Check if it ran successfully

### View Releases

1. Go to **Releases** page in GitHub
2. See all published versions
3. View changelog for each release

### PR Comments

When a release is created, the workflow automatically comments on associated PRs:

```
🚀 This PR has been released in v1.2.3
```

## ⚠️ Important Notes

### Commit Message Best Practices

1. **Be clear**: `feat: add retry` ✅ vs `feat: stuff` ❌
2. **One feature per commit**: Break large changes into multiple commits
3. **Use scope when helpful**: `feat(client): add retry mechanism`
4. **Document breaking changes**: Always include BREAKING CHANGE footer

### Common Mistakes

❌ **Wrong:** `added new feature` (no type prefix)
✅ **Correct:** `feat: add new feature`

❌ **Wrong:** `feat: stuff` (unclear)
✅ **Correct:** `feat: add retry mechanism for failed publishes`

❌ **Wrong:** Breaking change without footer
```bash
git commit -m "feat: change client API"
```
✅ **Correct:** Breaking change with BREAKING CHANGE footer
```bash
git commit -m "feat: change client API

BREAKING CHANGE: removed old method, use new method instead"
```

## 🐛 Troubleshooting

### No Release Created After Merge

**Possible reasons:**
- Commit messages don't follow conventional format
- Only `docs:`, `test:`, or `chore:` commits (these don't trigger releases)
- Workflow failed (check Actions tab)
- Tests or linter failed

**Solution:**
1. Check Actions tab for errors
2. Verify commit message format
3. Ensure at least one `feat:` or `fix:` commit

### Wrong Version Number

**Problem:** Expected v1.1.0, got v1.0.1

**Reason:** Used `fix:` instead of `feat:`

**Solution:** 
- Use `feat:` for new features (minor bump)
- Use `fix:` for bug fixes (patch bump)
- Use `BREAKING CHANGE:` for major bumps

### Release Created But No Git Tag

**Problem:** Release appears in GitHub but `go get` doesn't find it

**Solution:**
1. Wait 5-10 minutes (Go module index delay)
2. Try `GOPROXY=direct go get github.com/TogoMQ/togomq-sdk-go@v1.2.3`
3. Verify tag exists: `git ls-remote --tags origin`

### Tests Failed, No Release Created

**This is correct behavior!** 

The workflow blocks releases if:
- Tests fail
- Linter fails

**Solution:**
1. Fix the failing tests or linter issues
2. Push fix to main
3. Workflow will run again and create release

## 📚 Additional Resources

- **Conventional Commits**: https://www.conventionalcommits.org/
- **Semantic Versioning**: https://semver.org/
- **GitHub Actions**: https://docs.github.com/en/actions

## 🎉 You're All Set!

Your automatic release system is now configured. Just:
1. Write code
2. Commit with conventional format
3. Merge to main
4. Get automatic releases! 🚀

---

**Questions?** Check the detailed documentation in [AGENTS.md](AGENTS.md) or open an issue.
