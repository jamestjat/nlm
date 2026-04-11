//go:build windows

package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfilePathsUseUserHomeFallback(t *testing.T) {
	home := t.TempDir()

	t.Setenv("LOCALAPPDATA", "")
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", "")
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")

	expectedChrome := filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data")
	if got := getProfilePath(); got != expectedChrome {
		t.Fatalf("getProfilePath() = %q, want %q", got, expectedChrome)
	}

	expectedCanary := filepath.Join(home, "AppData", "Local", "Google", "Chrome SxS", "User Data")
	if got := getCanaryProfilePath(); got != expectedCanary {
		t.Fatalf("getCanaryProfilePath() = %q, want %q", got, expectedCanary)
	}

	expectedBrave := filepath.Join(home, "AppData", "Local", "BraveSoftware", "Brave-Browser", "User Data")
	if got := getBraveProfilePath(); got != expectedBrave {
		t.Fatalf("getBraveProfilePath() = %q, want %q", got, expectedBrave)
	}
}

func TestProfilePathsReturnEmptyWithoutAbsoluteBase(t *testing.T) {
	t.Setenv("LOCALAPPDATA", "")
	t.Setenv("USERPROFILE", "relative-home")
	t.Setenv("HOME", "")
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")

	if got := getProfilePath(); got != "" {
		t.Fatalf("getProfilePath() = %q, want empty", got)
	}
	if got := getCanaryProfilePath(); got != "" {
		t.Fatalf("getCanaryProfilePath() = %q, want empty", got)
	}
	if got := getBraveProfilePath(); got != "" {
		t.Fatalf("getBraveProfilePath() = %q, want empty", got)
	}
}

func TestGetChromePathIgnoresRelativeCandidates(t *testing.T) {
	clearWindowsBrowserEnv(t)

	workDir := t.TempDir()
	chdir(t, workDir)

	relativePath := filepath.Join("Google", "Chrome", "Application", "chrome.exe")
	writeEmptyFile(t, filepath.Join(workDir, relativePath))

	got := getChromePath()
	if got == relativePath {
		t.Fatalf("getChromePath() returned relative path %q", got)
	}
	if got != "" && !filepath.IsAbs(got) {
		t.Fatalf("getChromePath() returned non-absolute path %q", got)
	}
}

func TestGetBravePathIgnoresRelativeCandidates(t *testing.T) {
	clearWindowsBrowserEnv(t)

	workDir := t.TempDir()
	chdir(t, workDir)

	relativePath := filepath.Join("BraveSoftware", "Brave-Browser", "Application", "brave.exe")
	writeEmptyFile(t, filepath.Join(workDir, relativePath))

	got := getBravePath()
	if got == relativePath {
		t.Fatalf("getBravePath() returned relative path %q", got)
	}
	if got != "" && !filepath.IsAbs(got) {
		t.Fatalf("getBravePath() returned non-absolute path %q", got)
	}
}

func TestGetBrowserPathForProfileIgnoresRelativeCanaryCandidate(t *testing.T) {
	clearWindowsBrowserEnv(t)

	workDir := t.TempDir()
	chdir(t, workDir)

	relativePath := filepath.Join("Google", "Chrome SxS", "Application", "chrome.exe")
	writeEmptyFile(t, filepath.Join(workDir, relativePath))

	got := getBrowserPathForProfile("Chrome Canary")
	if got == relativePath {
		t.Fatalf("getBrowserPathForProfile() returned relative path %q", got)
	}
	if got != "" && !filepath.IsAbs(got) {
		t.Fatalf("getBrowserPathForProfile() returned non-absolute path %q", got)
	}
}

func TestGetBrowserPathForProfileUsesAbsoluteCanaryInstall(t *testing.T) {
	clearWindowsBrowserEnv(t)

	localAppData := t.TempDir()
	t.Setenv("LOCALAPPDATA", localAppData)

	expected := filepath.Join(localAppData, "Google", "Chrome SxS", "Application", "chrome.exe")
	writeEmptyFile(t, expected)

	if got := getBrowserPathForProfile("Chrome Canary"); got != expected {
		t.Fatalf("getBrowserPathForProfile() = %q, want %q", got, expected)
	}
}

func clearWindowsBrowserEnv(t *testing.T) {
	t.Helper()
	t.Setenv("LOCALAPPDATA", "")
	t.Setenv("PROGRAMFILES", "")
	t.Setenv("PROGRAMFILES(X86)", "")
	t.Setenv("USERPROFILE", "relative-home")
	t.Setenv("HOME", "")
	t.Setenv("HOMEDRIVE", "")
	t.Setenv("HOMEPATH", "")
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir(%q): %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})
}

func writeEmptyFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}
