# things3-cli Release Flow

This repo is consumed by external users and agents. Treat releases as public API changes.

## Preconditions

- Work from `main`.
- Check `git status --short`; release scripts require a clean tree for publish.
- Confirm the next version from `CHANGELOG.md`, `git tag --sort=-version:refname`, and the scope of the change.
- Prefer SemVer:
  - patch for bug fixes and docs-only release support,
  - minor for new flags, commands, or output fields,
  - major for breaking command/output behavior.

## Required Validation

Run these before publishing:

```sh
make test
./scripts/build-release.sh vX.Y.Z
./scripts/release-notes.sh vX.Y.Z
```

Check the generated artifacts:

```sh
ls -la dist
cat dist/checksums.txt
tmpdir=$(mktemp -d)
tar -xzf dist/things-X.Y.Z-darwin-arm64.tar.gz -C "$tmpdir"
"$tmpdir/things" --version
rm -rf "$tmpdir"
```

## Changelog

Move `CHANGELOG.md` entries out of `Unreleased` into:

```md
## [X.Y.Z] - YYYY-MM-DD
```

The GitHub release body must be only the bullets for that version. `./scripts/release-notes.sh vX.Y.Z` is the source of truth.

## Homebrew Formula

Update the formula after artifacts exist:

```sh
./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap
```

Validate:

```sh
rg -n 'version "|download/v|sha256' Formula/things3-cli.rb ~/Developer/homebrew-tap/Formula/things3-cli.rb
```

After the GitHub release is live, verify install/upgrade from the tap when practical:

```sh
brew update
brew reinstall ossianhempel/tap/things3-cli
things --version
```

Do not claim Homebrew is ready until the formula references the released tag and tarball checksums.

## Publishing

Local publish:

```sh
./scripts/release.sh vX.Y.Z
```

CI publish:

```sh
git tag vX.Y.Z
git push origin vX.Y.Z
```

The release workflow also runs tests, builds artifacts, generates notes, and publishes the GitHub release.

## Post-Release Verification

Verify:

- GitHub release title is `things3-cli vX.Y.Z`.
- Release body contains only that version's changelog bullets.
- Assets include both darwin tarballs and `checksums.txt`.
- Formula version, URLs, and checksums point to the released version.
- Installed `things --version` reports `vX.Y.Z`.

## Agent Skill Sync

If command behavior changed, update:

- `skills/things/SKILL.md`
- `../agent-scripts/archived-skills/things/SKILL.md` when that mirror exists

If release mechanics changed, update:

- `.codex/skills/release-flow/SKILL.md`
- `.codex/skills/release-flow/references/release-flow.md`
- `docs/RELEASING.md`
