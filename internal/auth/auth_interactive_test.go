package auth

import (
	"errors"
	"testing"
	"time"
)

func TestAuthAttemptTimeout(t *testing.T) {
	t.Run("default timeout", func(t *testing.T) {
		ba := &BrowserAuth{}
		if got := ba.authAttemptTimeout(); got != defaultAuthAttemptTimeout {
			t.Fatalf("authAttemptTimeout() = %v, want %v", got, defaultAuthAttemptTimeout)
		}
	})

	t.Run("extends beyond keep-open window", func(t *testing.T) {
		ba := &BrowserAuth{keepOpenSeconds: 360}
		want := 7 * time.Minute
		if got := ba.authAttemptTimeout(); got != want {
			t.Fatalf("authAttemptTimeout() = %v, want %v", got, want)
		}
	})
}

func TestIsAuthPageURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{url: "https://accounts.google.com/v3/signin/identifier", want: true},
		{url: "https://notebooklm.google.com/login?continue=1", want: true},
		{url: "https://notebooklm.google.com/", want: false},
	}

	for _, tt := range tests {
		if got := isAuthPageURL(tt.url); got != tt.want {
			t.Fatalf("isAuthPageURL(%q) = %v, want %v", tt.url, got, tt.want)
		}
	}
}

func TestIsManualLoginError(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{err: errors.New("detected sign-in page - not authenticated"), want: true},
		{err: errors.New("redirected to authentication page - not logged in"), want: true},
		{err: errors.New("missing essential authentication cookies"), want: false},
	}

	for _, tt := range tests {
		if got := isManualLoginError(tt.err); got != tt.want {
			t.Fatalf("isManualLoginError(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}
