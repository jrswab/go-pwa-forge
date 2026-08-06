package config

import (
	"os"
	"strings"
	"testing"
)

func TestParseCORSOrigins_Empty(t *testing.T) {
	for _, raw := range []string{"", "   ", "\t"} {
		got, err := parseCORSOrigins(raw)
		if err != nil {
			t.Fatalf("parseCORSOrigins(%q): unexpected err %v", raw, err)
		}
		if got != nil {
			t.Fatalf("parseCORSOrigins(%q): want nil, got %v", raw, got)
		}
	}
}

func TestParseCORSOrigins_SingleAndMultiple(t *testing.T) {
	got, err := parseCORSOrigins("https://app.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "https://app.example.com" {
		t.Fatalf("got %v", got)
	}

	got, err = parseCORSOrigins(" https://a.example.com , https://b.example.com:8443 ")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "https://a.example.com" || got[1] != "https://b.example.com:8443" {
		t.Fatalf("got %v", got)
	}
}

func TestParseCORSOrigins_ValidExactOrigins(t *testing.T) {
	cases := []string{
		"http://localhost",
		"https://app.example.com",
		"https://app.example.com:443",
		"http://127.0.0.1",
		"http://127.0.0.1:8080",
		"https://[::1]",
		"https://[::1]:8443",
		"https://[2001:db8::1]:443",
	}
	for _, raw := range cases {
		got, err := parseCORSOrigins(raw)
		if err != nil {
			t.Fatalf("parseCORSOrigins(%q): unexpected err %v", raw, err)
		}
		if len(got) != 1 || got[0] != raw {
			t.Fatalf("parseCORSOrigins(%q): got %v", raw, got)
		}
	}
}

func TestParseCORSOrigins_EmptyElementsRejected(t *testing.T) {
	cases := []string{
		"https://a.example.com,",
		",https://a.example.com",
		"https://a.example.com,,https://b.example.com",
		",",
		" , ",
		"https://a.example.com, ,https://b.example.com",
	}
	for _, raw := range cases {
		if _, err := parseCORSOrigins(raw); err == nil {
			t.Fatalf("parseCORSOrigins(%q): want error for empty element", raw)
		}
	}
}

func TestParseCORSOrigins_MalformedRejected(t *testing.T) {
	cases := []string{
		"*",
		"https://*.example.com",
		"http://foo*.bar.com",
		"not-a-url",
		"ftp://x.example.com",
		"https://x.example.com/path",
		"https://x.example.com/",
		"https://user:pass@x.example.com",
		"https://x.example.com?q=1",
		"https://x.example.com#frag",
		"https://:443",
		"https://example.com:99999",
		"https://example.com:443:44",
		"https://example.com:",
		"https://example.com:0",
		"https://example.com:65536",
		"https://-bad.example.com",
		"https://::1",
	}
	for _, raw := range cases {
		if _, err := parseCORSOrigins(raw); err == nil {
			t.Fatalf("parseCORSOrigins(%q): want error", raw)
		}
	}
}

func TestValidateCORSOrigins(t *testing.T) {
	if err := ValidateCORSOrigins(nil); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCORSOrigins([]string{}); err != nil {
		t.Fatal(err)
	}
	for _, good := range [][]string{
		{"https://a.example.com"},
		{"http://localhost:3000", "https://app.example.com:443"},
		{"http://192.168.0.1", "https://[::1]", "https://[2001:db8::2]:8443"},
	} {
		if err := ValidateCORSOrigins(good); err != nil {
			t.Fatalf("ValidateCORSOrigins(%v): unexpected err %v", good, err)
		}
	}
	for _, bad := range [][]string{
		{"*"},
		{"https://ok.example.com", "*"},
		{"https://*.example.com"},
		{""},
		{"not-a-url"},
		{"https://:443"},
		{"https://example.com:99999"},
		{"https://example.com:443:44"},
		{"https://example.com:"},
		{"https://example.com:0"},
		{"https://-bad.example.com"},
		{"https://ok.example.com", "https://example.com:443:44"},
	} {
		if err := ValidateCORSOrigins(bad); err == nil {
			t.Fatalf("ValidateCORSOrigins(%v): want error", bad)
		}
	}
}

func TestLoad_CORSFromEnv(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://hda.site.com,https://app.site.com")
	// clear unrelated noise envs
	for _, k := range []string{"REDIS_HOST", "REDIS_PASS", "DB_USER", "DB_PASS", "DB_HOST", "DB_PORT", "DB_NAME"} {
		_ = os.Unsetenv(k)
	}

	cfg := new(Config)
	if err := cfg.Load(); err != nil {
		t.Fatal(err)
	}
	if len(cfg.CORSAllowedOrigins) != 2 {
		t.Fatalf("origins: %v", cfg.CORSAllowedOrigins)
	}

	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	cfg = new(Config)
	if err := cfg.Load(); err != nil {
		t.Fatal(err)
	}
	if cfg.CORSAllowedOrigins != nil {
		t.Fatalf("want nil origins, got %v", cfg.CORSAllowedOrigins)
	}

	t.Setenv("CORS_ALLOWED_ORIGINS", "*")
	cfg = new(Config)
	if err := cfg.Load(); err == nil {
		t.Fatal("want error for wildcard CORS config")
	}

	t.Setenv("CORS_ALLOWED_ORIGINS", "https://a.example.com,")
	cfg = new(Config)
	if err := cfg.Load(); err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("want empty-entry error, got %v", err)
	}
}
