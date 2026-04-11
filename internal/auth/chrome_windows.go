//go:build windows

package auth

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func detectChrome(debug bool) Browser {
	path := getChromePath()
	if path == "" {
		return Browser{Type: BrowserUnknown}
	}

	version := getChromeVersion(path)
	return Browser{
		Type:    BrowserChrome,
		Path:    path,
		Name:    "Google Chrome",
		Version: version,
	}
}

func getChromeVersion(path string) string {
	cmd := exec.Command(path, "--version")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(strings.TrimPrefix(string(out), "Google Chrome "))
}

func getProfilePath() string {
	return profilePathFor("Google", "Chrome")
}

func getChromePath() string {
	return browserExecutablePath("Google", "Chrome", "chrome.exe")
}

func getBrowserPathForProfile(browserName string) string {
	switch browserName {
	case "Brave":
		if path := getBravePath(); path != "" {
			return path
		}
	case "Chrome Canary":
		if path := getCanaryPath(); path != "" {
			return path
		}
	}
	return getChromePath()
}

func getBravePath() string {
	return browserExecutablePath("BraveSoftware", "Brave-Browser", "brave.exe")
}

func getCanaryPath() string {
	return browserExecutablePath("Google", "Chrome SxS", "chrome.exe")
}

func getCanaryProfilePath() string {
	return profilePathFor("Google", "Chrome SxS")
}

func getBraveProfilePath() string {
	return profilePathFor("BraveSoftware", "Brave-Browser")
}

func profilePathFor(vendor, product string) string {
	localAppData := getLocalAppDataPath()
	if localAppData == "" {
		return ""
	}
	return filepath.Join(localAppData, vendor, product, "User Data")
}

func browserExecutablePath(vendor, product, executable string) string {
	for _, root := range browserInstallRoots() {
		path := filepath.Join(root, vendor, product, "Application", executable)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func browserInstallRoots() []string {
	candidates := []string{
		os.Getenv("PROGRAMFILES"),
		os.Getenv("PROGRAMFILES(X86)"),
		getLocalAppDataPath(),
		`C:\Program Files`,
		`C:\Program Files (x86)`,
	}
	return uniqueAbsolutePaths(candidates)
}

func getLocalAppDataPath() string {
	if path := absoluteWindowsPath(os.Getenv("LOCALAPPDATA")); path != "" {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	home = absoluteWindowsPath(home)
	if home == "" {
		return ""
	}

	return filepath.Join(home, "AppData", "Local")
}

func uniqueAbsolutePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	result := make([]string, 0, len(paths))

	for _, path := range paths {
		path = absoluteWindowsPath(path)
		if path == "" {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		seen[path] = struct{}{}
		result = append(result, path)
	}

	return result
}

func absoluteWindowsPath(path string) string {
	if path == "" || !filepath.IsAbs(path) {
		return ""
	}
	return path
}
