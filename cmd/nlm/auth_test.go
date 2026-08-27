package main

import (
	"runtime"
	"testing"
)

func TestDetectAuthInfo(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", t.TempDir())
	} else {
		t.Setenv("HOME", t.TempDir())
	}

	tests := []struct {
		name        string
		input       string
		wantToken   string
		wantCookies string
		wantErr     bool
	}{
		{
			name:        "cookie header format",
			input:       `curl 'https://notebooklm.google.com/_/batchexecute' -H 'cookie: SID=abc; HSID=def' --data-raw 'f.req=x&at=token123&'`,
			wantToken:   "token123",
			wantCookies: "SID=abc; HSID=def",
		},
		{
			name:        "capitalized Cookie header",
			input:       `curl 'https://notebooklm.google.com/_/batchexecute' -H 'Cookie: SID=abc' --data-raw 'at=tok&'`,
			wantToken:   "tok",
			wantCookies: "SID=abc",
		},
		{
			name:        "curl -b flag format",
			input:       "curl --url 'https://notebook.google.com/_/batchexecute' \\\n  -b 'SID=abc; HSID=def' \\\n  --data-raw 'f.req=x&at=token456&'",
			wantToken:   "token456",
			wantCookies: "SID=abc; HSID=def",
		},
		{
			name:        "curl --cookie flag format",
			input:       `curl 'https://x' --cookie 'SID=abc' --data-raw 'at=tok789&'`,
			wantToken:   "tok789",
			wantCookies: "SID=abc",
		},
		{
			name:        "url-encoded token is decoded",
			input:       `curl 'https://notebook.google.com/_/batchexecute?bl=boq_labs-tailwind-frontend_20260825.16_p2' -b 'SID=abc' --data-raw 'f.req=x&at=AOr%3Aabc-def&'`,
			wantToken:   "AOr:abc-def",
			wantCookies: "SID=abc",
		},
		{
			name:    "no cookies",
			input:   `curl 'https://x' --data-raw 'at=tok&'`,
			wantErr: true,
		},
		{
			name:    "no at token",
			input:   `curl 'https://x' -b 'SID=abc' --data-raw 'f.req=x'`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, cookies, err := detectAuthInfo(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got token=%q cookies=%q", token, cookies)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if token != tt.wantToken {
				t.Errorf("token = %q, want %q", token, tt.wantToken)
			}
			if cookies != tt.wantCookies {
				t.Errorf("cookies = %q, want %q", cookies, tt.wantCookies)
			}
		})
	}
}
