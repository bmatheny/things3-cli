---
summary: 'things3-cli release checklist'
read_when:
  - preparing a release
  - writing release notes
---

# Releasing things3-cli

## Guardrails

- Title every GitHub release as `things3-cli <version>`.
- Release body = the CHANGELOG bullets for that version only (no extra prose).
- Do not call a release complete until tests, artifacts, release notes,
  Homebrew formula checksums, GitHub release assets, and the installed binary
  version have been verified.

## Checklist

- [ ] Confirm you are on `main`, the working tree is clean or only contains the
      intended release edits, and the next version matches SemVer scope.
- [ ] Update `CHANGELOG.md`: move items from **Unreleased** into a new version section with today’s date.
- [ ] Run tests: `make test`.
- [ ] Build release artifacts: `./scripts/build-release.sh vX.Y.Z`.
- [ ] Smoke-test an artifact binary and confirm `things --version` reports `vX.Y.Z`.
- [ ] Update Homebrew formula:
  `./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap`
  (commits to `Formula/` and copies into the tap repo).
- [ ] Verify `Formula/things3-cli.rb` and the tap formula point at the new
      release URLs and checksums.
- [ ] Generate release notes: `./scripts/release-notes.sh vX.Y.Z` (should output only the changelog bullets).
- [ ] Commit the changelog/formula/docs changes.
- [ ] Tag the release: `git tag vX.Y.Z` then `git push origin vX.Y.Z`.
- [ ] Create the GitHub release:
  - Option A (local): `./scripts/release.sh vX.Y.Z`
  - Option B (CI): push the tag and let `.github/workflows/release.yml` publish the release
- [ ] Verify the GitHub release title, notes, and assets match the guardrails.
- [ ] Verify Homebrew install or reinstall from `ossianhempel/tap/things3-cli`
      and confirm `things --version`.
