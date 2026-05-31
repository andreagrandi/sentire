package tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestGoInstallSmoke installs the sentire main package via `go install` into a
// temporary GOBIN and then exercises `sentire version` and `sentire --help`.
// This mirrors the public install path advertised in the README
// (`go install github.com/andreagrandi/sentire/cmd/sentire@latest`) without
// pulling from VCS or requiring Sentry credentials.
func TestGoInstallSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping go install smoke test in short mode")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	moduleRoot := filepath.Dir(wd)

	gobin := t.TempDir()
	gocache := t.TempDir()

	install := exec.Command("go", "install", "github.com/andreagrandi/sentire/cmd/sentire")
	install.Dir = moduleRoot
	install.Env = append(os.Environ(),
		"GOBIN="+gobin,
		"GOCACHE="+gocache,
		"GOFLAGS=-mod=mod",
	)
	var installErr bytes.Buffer
	install.Stderr = &installErr
	if err := install.Run(); err != nil {
		t.Fatalf("go install failed: %v\nstderr: %s", err, installErr.String())
	}

	binary := filepath.Join(gobin, "sentire")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if _, err := os.Stat(binary); err != nil {
		t.Fatalf("expected sentire binary at %s: %v", binary, err)
	}

	assertSmokeBinary(t, binary)
}

// assertSmokeBinary runs `version` and `--help` on the given binary and checks
// that both succeed and produce the expected user-facing output. It is shared
// with the release-artifact smoke test so both code paths agree on what a
// working install looks like.
func assertSmokeBinary(t *testing.T, binary string) {
	t.Helper()

	versionOut, err := exec.Command(binary, "version").CombinedOutput()
	if err != nil {
		t.Fatalf("sentire version failed: %v\noutput: %s", err, versionOut)
	}
	if !strings.Contains(string(versionOut), "sentire version") {
		t.Errorf("sentire version output missing expected prefix:\n%s", versionOut)
	}

	helpOut, err := exec.Command(binary, "--help").CombinedOutput()
	if err != nil {
		t.Fatalf("sentire --help failed: %v\noutput: %s", err, helpOut)
	}
	help := string(helpOut)
	if !strings.Contains(help, "Usage:") {
		t.Errorf("sentire --help output missing 'Usage:':\n%s", help)
	}
	if !strings.Contains(help, "sentire") {
		t.Errorf("sentire --help output missing program name:\n%s", help)
	}
}

// TestGoreleaserSnapshotSmoke builds the release artifacts via
// `goreleaser release --snapshot`, then unpacks the archive for the current
// OS/arch and runs the standard smoke checks on the extracted binary.
//
// The test is opt-in (set SENTIRE_SMOKE_RELEASE=1) because a full snapshot
// build cross-compiles for every platform and adds ~30s+ to the suite. The
// Makefile target `make smoke-release` flips the gate.
func TestGoreleaserSnapshotSmoke(t *testing.T) {
	if os.Getenv("SENTIRE_SMOKE_RELEASE") == "" {
		t.Skip("set SENTIRE_SMOKE_RELEASE=1 to run release artifact smoke test (or use `make smoke-release`)")
	}
	if _, err := exec.LookPath("goreleaser"); err != nil {
		t.Skip("goreleaser not installed; skipping release artifact smoke test")
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	moduleRoot := filepath.Dir(wd)
	distDir := filepath.Join(moduleRoot, "dist")

	_ = os.RemoveAll(distDir)
	t.Cleanup(func() { _ = os.RemoveAll(distDir) })

	build := exec.Command(
		"goreleaser", "release",
		"--snapshot", "--clean",
		"--skip=publish,announce,validate,homebrew,before",
		"--timeout=10m",
	)
	build.Dir = moduleRoot
	var buildOut bytes.Buffer
	build.Stdout = &buildOut
	build.Stderr = &buildOut
	if err := build.Run(); err != nil {
		t.Fatalf("goreleaser snapshot failed: %v\noutput:\n%s", err, buildOut.String())
	}

	archive := findReleaseArchive(t, distDir)
	binary := extractSentireBinary(t, archive)
	assertSmokeBinary(t, binary)
}

// findReleaseArchive returns the archive path for the current OS/arch produced
// by `goreleaser release --snapshot`. The name template lives in
// .goreleaser.yaml and follows `sentire_<Os>_<Arch>.<ext>`.
func findReleaseArchive(t *testing.T, distDir string) string {
	t.Helper()

	osName := goreleaserOSName(runtime.GOOS)
	arch := goreleaserArchName(runtime.GOARCH)
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}

	pattern := filepath.Join(distDir, fmt.Sprintf("sentire_%s_%s.%s", osName, arch, ext))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("glob %s failed: %v", pattern, err)
	}
	if len(matches) == 0 {
		entries, _ := os.ReadDir(distDir)
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("no archive matched %s; dist contents: %v", pattern, names)
	}
	return matches[0]
}

// extractSentireBinary extracts the sentire binary from a goreleaser archive
// into a temp directory and returns its path.
func extractSentireBinary(t *testing.T, archive string) string {
	t.Helper()

	dir := t.TempDir()
	var extract *exec.Cmd
	if strings.HasSuffix(archive, ".zip") {
		extract = exec.Command("unzip", "-o", archive, "-d", dir)
	} else {
		extract = exec.Command("tar", "-xzf", archive, "-C", dir)
	}
	if out, err := extract.CombinedOutput(); err != nil {
		t.Fatalf("failed to extract %s: %v\noutput: %s", archive, err, out)
	}

	binary := filepath.Join(dir, "sentire")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	if _, err := os.Stat(binary); err != nil {
		t.Fatalf("expected sentire binary in extracted archive at %s: %v", binary, err)
	}
	return binary
}

func goreleaserOSName(goos string) string {
	switch goos {
	case "darwin":
		return "Darwin"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		if goos == "" {
			return goos
		}
		return strings.ToUpper(goos[:1]) + goos[1:]
	}
}

func goreleaserArchName(goarch string) string {
	switch goarch {
	case "amd64":
		return "x86_64"
	case "386":
		return "i386"
	default:
		return goarch
	}
}
