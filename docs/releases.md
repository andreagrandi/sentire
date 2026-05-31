# Releasing Sentire

This guide covers the steps to cut a new release and the verification you
should perform on each public install path before announcing the version.

The release pipeline itself is automated: pushing a tag matching `v*` to
GitHub triggers the `Release` workflow, which runs the test suite and then
invokes [GoReleaser](https://goreleaser.com) (see `.goreleaser.yaml`) to
build cross-platform archives, publish them to the GitHub release page, and
update the [Homebrew tap](https://github.com/andreagrandi/homebrew-tap).

For the version-bump and tagging workflow itself, see
`.opencode/agent/new-release.md`.

## Pre-release smoke tests

Before tagging, run the local smoke tests. Neither requires Sentry
credentials.

```bash
# Verifies the public `go install` path against the local module.
make smoke-install

# Builds a full goreleaser snapshot, extracts the current OS/arch archive,
# and runs `sentire version` and `sentire --help` against the unpacked
# binary. Requires goreleaser to be installed locally (`brew install
# goreleaser`).
make smoke-release
```

`make smoke-install` runs as part of normal development (it is included in
`make test` via `TestGoInstallSmoke`). `make smoke-release` is opt-in
because a full snapshot build cross-compiles for every supported platform
and adds noticeable time; it is also skipped automatically when goreleaser
is not on PATH, so it does not break CI environments that lack it.

## Post-release verification

After the `Release` workflow finishes, verify each advertised install path
end to end. Replace `<version>` with the tag you just pushed (without the
leading `v` for `go install` queries; with it for the GitHub release page).

### 1. GitHub release artifacts

1. Open the release page:
   <https://github.com/andreagrandi/sentire/releases/tag/v\<version\>>
2. Confirm archives for `Darwin_x86_64`, `Darwin_arm64`, `Linux_x86_64`,
   `Linux_arm64`, and `Windows_x86_64` are attached, along with
   `checksums.txt`.
3. Download the archive for your platform, extract it, and run:

   ```bash
   ./sentire version    # should print "sentire version <version>"
   ./sentire --help     # should print top-level usage
   ```

### 2. `go install`

In a clean shell (no `SENTRY_API_TOKEN` set), run:

```bash
GOBIN="$(mktemp -d)"
export GOBIN
go install github.com/andreagrandi/sentire/cmd/sentire@v<version>
"$GOBIN/sentire" version    # should report v<version>
"$GOBIN/sentire" --help
```

If `go install @latest` should also resolve to the new tag, run the same
check with `@latest`.

### 3. Homebrew

The release workflow pushes an updated formula to
`andreagrandi/homebrew-tap`. After the workflow finishes, verify the tap
end to end:

```bash
brew update
brew install andreagrandi/tap/sentire   # or `brew upgrade ...` if installed
sentire version    # should report the new version
sentire --help
```

If the formula has not yet been updated, the workflow log under
`Release → goreleaser` will show the push to the tap repository — check
both this repo's Actions tab and the tap repo's commit history.

## Troubleshooting

- **`go install` returns an older version.** Module proxies cache versions;
  `GOPROXY=direct go install github.com/andreagrandi/sentire/cmd/sentire@v<version>`
  forces a fresh fetch from GitHub.
- **`brew install` reports "no such formula".** The tap may not be tapped
  yet. Run `brew tap andreagrandi/tap` first, or use the fully qualified
  formula name `andreagrandi/tap/sentire`.
- **`sentire version` reports `unknown` or a stale version.** The release
  archives are built with ldflags that inject the version
  (`internal/version.Version`). If the value is wrong, the goreleaser
  config drifted — check `.goreleaser.yaml` and the most recent commit on
  the tag.
