---
description: Create a new release for this application, bumping the version and updating the changelog
mode: subagent
---

When asked to create a new release, you need to:
- Make sure `make test` passes without errors (this also runs the
  `go install` smoke test via `TestGoInstallSmoke`)
- If `goreleaser` is installed locally, also run `make smoke-release` to
  verify the release artifacts build and the unpacked binary works
- Use the provided version number or
  Bump the version number in internal/version/version.go:
  if you find `const Version = "0.1.0"` change to `const Version = "0.1.1"`
- Update the changelog writing a short summary of the changes since last release (with bullet points), follow existing format
- git commit the changes you just did
- git push the changes you just did
- do `git tag v<version>` (use the version you just bumped to in the `internal/version/version.go`)
- do `git push origin v<version>`
- After the GitHub `Release` workflow finishes, run the post-release
  verification checklist in `docs/releases.md` (GitHub artifact, `go install`,
  and Homebrew)