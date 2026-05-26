---
name: release-flow
description: Maintainer release workflow for things3-cli. Use when preparing, validating, tagging, publishing, or verifying a public things3-cli release or Homebrew formula update.
---

# release-flow

Use this skill when shipping `things3-cli` to external users and agents.

Start here:
- Read `docs/RELEASING.md`.
- Read `references/release-flow.md` for the validation sequence and command details.

Core rule:
- Do not publish a release from assumptions. Verify tests, artifacts, changelog notes, formula checksums, GitHub release assets, and the installed binary path before calling the release done.

Tagging and the GitHub release are automated: merging a PR that adds a new
`## [X.Y.Z]` section to `CHANGELOG.md` makes `.github/workflows/release.yml` tag
and publish `vX.Y.Z`. Never tag by hand. The Homebrew formulas are NOT in CI —
the agent regenerates them from the RELEASED checksums after the release is live
(builds are not byte-reproducible, so locally built checksums will not match what
users download).

Quick path:
1. Confirm the release version (SemVer scope of the merged changes).
2. Add a `## [X.Y.Z] - YYYY-MM-DD` section to `CHANGELOG.md` with the bullets.
   Do NOT bump the formula in this PR — its checksums aren't known yet.
3. Run `make test`; build (`./scripts/build-release.sh vX.Y.Z`) and smoke-test a
   built binary's `--version`.
4. Sanity-check release notes with `./scripts/release-notes.sh vX.Y.Z`.
5. Open a release PR and merge it — CI tags and publishes the release.
6. After CI publishes, verify the GitHub release, then regenerate both formulas
   from the released checksums:
   `gh release download vX.Y.Z -p checksums.txt -O dist/checksums.txt --clobber`
   then `./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap`.
7. Commit + push `~/Developer/homebrew-tap`; open a follow-up PR for the in-repo
   `Formula/things3-cli.rb`.
8. Verify the Homebrew install path (`brew reinstall ossianhempel/tap/things3-cli`).

When CLI behavior changes, also update `skills/things/SKILL.md` and the mirrored `../agent-scripts/archived-skills/things/SKILL.md` if that mirror exists.
