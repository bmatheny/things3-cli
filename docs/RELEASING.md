---
summary: 'things3-cli release checklist'
read_when:
  - preparing a release
  - writing release notes
---

# Releasing things3-cli

## How releases ship

Tagging and the GitHub release are **automated**. `.github/workflows/release.yml`
watches `main`: when a merged commit adds a new `## [X.Y.Z]` section to
`CHANGELOG.md`, CI tags `vX.Y.Z` (at the merge commit), runs tests, builds the
darwin tarballs, generates the release notes from that changelog section, and
publishes `things3-cli vX.Y.Z`. The workflow is idempotent — a later CHANGELOG
edit that only touches **Unreleased** finds the tag already present and skips.
A `workflow_dispatch` (Actions → Release → Run workflow, enter the version)
is the manual override.

So nobody tags by hand. The maintainer's only required action is to **merge the
release PR**. Both Homebrew formulas (in-repo `Formula/` and the tap) are updated
**after** the release publishes — by the agent, not in CI (see below).

> **Why the formula is updated after the release, not in the PR:** the release
> tarballs are built by CI, and the build is **not** byte-reproducible (a locally
> built tarball gets a different sha256 than CI's). So the formula must reference
> the **released** artifacts' checksums. Never put locally built checksums in the
> formula — they will not match what users download.

## Guardrails

- Title every GitHub release as `things3-cli <version>` (CI does this).
- Release body = the CHANGELOG bullets for that version only (no extra prose).
- Do not call a release complete until tests, artifacts, release notes,
  Homebrew formula checksums, GitHub release assets, and the installed binary
  version have been verified.

## Checklist (the release PR)

- [ ] Confirm the next version matches SemVer scope for the merged changes.
- [ ] Update `CHANGELOG.md`: add a new `## [X.Y.Z] - YYYY-MM-DD` section with
      the bullets for this release. (Do **not** bump the formula here — see the
      note above.)
- [ ] Run tests: `make test`.
- [ ] Build and smoke-test locally: `./scripts/build-release.sh vX.Y.Z` then
      confirm the built binary's `things --version` reports `vX.Y.Z`.
- [ ] Sanity-check release notes: `./scripts/release-notes.sh vX.Y.Z`
      (should output only the changelog bullets).
- [ ] Open the PR and merge it. CI tags and publishes the release automatically.

## Checklist (after CI publishes — the agent does this, not the maintainer)

- [ ] Confirm the GitHub release `things3-cli vX.Y.Z` is live with both darwin
      tarballs and `checksums.txt` attached, and the body is the changelog bullets.
- [ ] Generate both formulas from the **released** checksums (download the
      published checksums first so they match what users get):
  `gh release download vX.Y.Z -p checksums.txt -O dist/checksums.txt --clobber`
  then
  `./scripts/update-brew-formula.sh --version vX.Y.Z --tap-dir ~/Developer/homebrew-tap`
  (writes `Formula/things3-cli.rb` and copies it into the tap).
- [ ] Confirm `Formula/things3-cli.rb` and the tap formula sha256s equal the
      released `checksums.txt` values.
- [ ] Commit and push the tap repo; open a follow-up PR for the in-repo
      `Formula/things3-cli.rb` (its checksums change is too late for the release PR).
- [ ] Verify Homebrew install or reinstall from `ossianhempel/tap/things3-cli`
      and confirm `things --version`.
