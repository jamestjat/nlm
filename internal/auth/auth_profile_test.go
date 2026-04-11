package auth

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestScanBrowserProfilesSupportsNetworkCookies(t *testing.T) {
	root := t.TempDir()
	profileDir := filepath.Join(root, "Default")
	writeProfileTestFile(t, filepath.Join(profileDir, "Network", "Cookies"), strings.Repeat("c", 2048))

	profiles, err := scanBrowserProfiles(root, "Chrome", "notebooklm.google.com")
	if err != nil {
		t.Fatalf("scanBrowserProfiles() error = %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("scanBrowserProfiles() returned %d profiles, want 1", len(profiles))
	}

	profile := profiles[0]
	if profile.Path != profileDir {
		t.Fatalf("profile.Path = %q, want %q", profile.Path, profileDir)
	}
	if !profile.HasTargetCookies {
		t.Fatalf("profile.HasTargetCookies = false, want true")
	}
	if !containsString(profile.Files, "Network/Cookies") {
		t.Fatalf("profile.Files = %v, want Network/Cookies marker", profile.Files)
	}
}

func TestCopyProfileDataFromPathCopiesNetworkCookiesAndLocalState(t *testing.T) {
	userDataDir := t.TempDir()
	sourceDir := filepath.Join(userDataDir, "Default")
	writeProfileTestFile(t, filepath.Join(sourceDir, "Network", "Cookies"), "cookie-db")
	writeProfileTestFile(t, filepath.Join(sourceDir, "Network", "Cookies-wal"), "cookie-wal")

	localState := `{"os_crypt":{"encrypted_key":"test-key"}}`
	writeProfileTestFile(t, filepath.Join(userDataDir, "Local State"), localState)

	ba := &BrowserAuth{tempDir: t.TempDir()}
	if err := ba.copyProfileDataFromPath(sourceDir); err != nil {
		t.Fatalf("copyProfileDataFromPath() error = %v", err)
	}

	copiedCookies, err := os.ReadFile(filepath.Join(ba.tempDir, "Default", "Network", "Cookies"))
	if err != nil {
		t.Fatalf("ReadFile(copied cookies): %v", err)
	}
	if string(copiedCookies) != "cookie-db" {
		t.Fatalf("copied cookies = %q, want %q", copiedCookies, "cookie-db")
	}

	copiedCookieWAL, err := os.ReadFile(filepath.Join(ba.tempDir, "Default", "Network", "Cookies-wal"))
	if err != nil {
		t.Fatalf("ReadFile(copied cookie wal): %v", err)
	}
	if string(copiedCookieWAL) != "cookie-wal" {
		t.Fatalf("copied cookie wal = %q, want %q", copiedCookieWAL, "cookie-wal")
	}

	copiedLocalState, err := os.ReadFile(filepath.Join(ba.tempDir, "Local State"))
	if err != nil {
		t.Fatalf("ReadFile(Local State): %v", err)
	}
	if string(copiedLocalState) != localState {
		t.Fatalf("copied Local State = %q, want %q", copiedLocalState, localState)
	}
}

func TestFindMostRecentProfileSupportsNetworkCookies(t *testing.T) {
	root := t.TempDir()
	olderProfile := filepath.Join(root, "Profile 1")
	newerProfile := filepath.Join(root, "Default")

	writeProfileTestFile(t, filepath.Join(olderProfile, "History"), "history")
	writeProfileTestFile(t, filepath.Join(newerProfile, "Network", "Cookies"), "cookies")

	setPathModTime(t, olderProfile, time.Now().Add(-2*time.Hour))
	setPathModTime(t, newerProfile, time.Now().Add(-time.Hour))

	if got := findMostRecentProfile(root); got != newerProfile {
		t.Fatalf("findMostRecentProfile() = %q, want %q", got, newerProfile)
	}
}

func writeProfileTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

func setPathModTime(t *testing.T, path string, modTime time.Time) {
	t.Helper()
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("Chtimes(%q): %v", path, err)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
