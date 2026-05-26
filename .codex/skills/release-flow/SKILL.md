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

Quick path:
1. Confirm the working tree and release version.
2. Move `CHANGELOG.md` entries from `Unreleased` into `vX.Y.Z` with today's date.
3. Run `make test`.
4. Build artifacts with `./scripts/build-release.sh vX.Y.Z`.
5. Update the formula with `./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap`.
6. Verify release notes with `./scripts/release-notes.sh vX.Y.Z`.
7. Commit the changelog/formula/docs changes.
8. Tag and publish with either the GitHub Actions release workflow or `./scripts/release.sh vX.Y.Z`.
9. Verify the GitHub release and Homebrew install path.

When CLI behavior changes, also update `skills/things/SKILL.md` and the mirrored `../agent-scripts/archived-skills/things/SKILL.md` if that mirror exists.
