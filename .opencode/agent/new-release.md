---
description: Create a new release for this application, bumping the version and updating the changelog
mode: subagent
---

When asked to create a new release, you need to:
- Make sure `make test` passes without errors
- Use the provided version number or
  Bump the version number in internal/version/version.go:
  if you find `const Version = "0.1.0"` change to `const Version = "0.1.1"`
- Update the changelog writing a short summary of the changes since last release (with bullet points), follow existing format
- git commit the changes you just did
- git push the changes you just did
- do `git tag v<version>` (use the version you just bumped to in the `internal/version/version.go`)
- do `git push origin v<version>`