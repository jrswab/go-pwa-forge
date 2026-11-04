package config

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strings"
	"testing"
)

func TestParseCORSOrigins_Empty(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestLoad_CORS(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		env         map[string]string
		want        []string
		wantErr     bool
		errContains string
	}{
		{name: "unset"},
		{name: "empty", env: map[string]string{"CORS_ALLOWED_ORIGINS": ""}},
		{name: "whitespace only", env: map[string]string{"CORS_ALLOWED_ORIGINS": " \t\n "}},
		{name: "single origin", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://app.example.com"}, want: []string{"https://app.example.com"}},
		{name: "multiple origins", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://hda.site.com,https://app.site.com"}, want: []string{"https://hda.site.com", "https://app.site.com"}},
		{name: "trim entries", env: map[string]string{"CORS_ALLOWED_ORIGINS": " https://a.example.com , https://b.example.com:8443 "}, want: []string{"https://a.example.com", "https://b.example.com:8443"}},
		{name: "localhost", env: map[string]string{"CORS_ALLOWED_ORIGINS": "http://localhost"}, want: []string{"http://localhost"}},
		{name: "explicit port", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://app.example.com:443"}, want: []string{"https://app.example.com:443"}},
		{name: "IPv4", env: map[string]string{"CORS_ALLOWED_ORIGINS": "http://127.0.0.1"}, want: []string{"http://127.0.0.1"}},
		{name: "IPv4 with port", env: map[string]string{"CORS_ALLOWED_ORIGINS": "http://127.0.0.1:8080"}, want: []string{"http://127.0.0.1:8080"}},
		{name: "IPv6", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://[::1]"}, want: []string{"https://[::1]"}},
		{name: "IPv6 with port", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://[::1]:8443"}, want: []string{"https://[::1]:8443"}},
		{name: "IPv6 non-loopback", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://[2001:db8::1]:443"}, want: []string{"https://[2001:db8::1]:443"}},
		{name: "trailing comma", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://a.example.com,"}, wantErr: true, errContains: "empty"},
		{name: "leading comma", env: map[string]string{"CORS_ALLOWED_ORIGINS": ",https://a.example.com"}, wantErr: true, errContains: "empty"},
		{name: "repeated comma", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://a.example.com,,https://b.example.com"}, wantErr: true, errContains: "empty"},
		{name: "only comma", env: map[string]string{"CORS_ALLOWED_ORIGINS": ","}, wantErr: true, errContains: "empty"},
		{name: "comma with whitespace", env: map[string]string{"CORS_ALLOWED_ORIGINS": " , "}, wantErr: true, errContains: "empty"},
		{name: "blank entry", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://a.example.com, ,https://b.example.com"}, wantErr: true, errContains: "empty"},
		{name: "wildcard", env: map[string]string{"CORS_ALLOWED_ORIGINS": "*"}, wantErr: true},
		{name: "wildcard host", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://*.example.com"}, wantErr: true},
		{name: "embedded wildcard", env: map[string]string{"CORS_ALLOWED_ORIGINS": "http://foo*.bar.com"}, wantErr: true},
		{name: "not a URL", env: map[string]string{"CORS_ALLOWED_ORIGINS": "not-a-url"}, wantErr: true},
		{name: "unsupported scheme", env: map[string]string{"CORS_ALLOWED_ORIGINS": "ftp://x.example.com"}, wantErr: true},
		{name: "path", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://x.example.com/path"}, wantErr: true},
		{name: "trailing slash", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://x.example.com/"}, wantErr: true},
		{name: "userinfo", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://user:pass@x.example.com"}, wantErr: true},
		{name: "query", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://x.example.com?q=1"}, wantErr: true},
		{name: "fragment", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://x.example.com#frag"}, wantErr: true},
		{name: "missing host", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://:443"}, wantErr: true},
		{name: "oversized port", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://example.com:99999"}, wantErr: true},
		{name: "multiple ports", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://example.com:443:44"}, wantErr: true},
		{name: "empty port", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://example.com:"}, wantErr: true},
		{name: "zero port", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://example.com:0"}, wantErr: true},
		{name: "port above maximum", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://example.com:65536"}, wantErr: true},
		{name: "invalid DNS label", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://-bad.example.com"}, wantErr: true},
		{name: "unbracketed IPv6", env: map[string]string{"CORS_ALLOWED_ORIGINS": "https://::1"}, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := tc.env

			var cfg Config
			cfg, err := Load(func(k string) string { return m[k] })

			if tc.wantErr {
				if err == nil {
					t.Fatal("Load: want error for malformed non-empty CORS config")
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("Load: want error containing %q, got %v", tc.errContains, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load: unexpected error %v", err)
			}
			if !slices.Equal(cfg.CORSAllowedOrigins, tc.want) {
				t.Fatalf("Load: want origins %v, got %v", tc.want, cfg.CORSAllowedOrigins)
			}
		})
	}
}

func TestLoad_ReadsOnlyCORSVar(t *testing.T) {
	t.Parallel()

	env := map[string]string{"CORS_ALLOWED_ORIGINS": "https://app.example.com"}
	var keys []string

	_, err := Load(func(key string) string {
		keys = append(keys, key)
		return env[key]
	})

	if err != nil {
		t.Fatalf("Load: unexpected error %v", err)
	}
	if !slices.Equal(keys, []string{"CORS_ALLOWED_ORIGINS"}) {
		t.Fatalf("Load: want only CORS_ALLOWED_ORIGINS lookup, got %v", keys)
	}
}

func TestLoad_NeverLogsOrReportsSecrets(t *testing.T) {
	// This test swaps the global logger.
	cases := []struct {
		name    string
		cors    string
		wantErr bool
	}{
		{name: "valid CORS", cors: "https://app.example.com"},
		{name: "malformed CORS", cors: "*", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := map[string]string{
				"CORS_ALLOWED_ORIGINS": tc.cors,
				"REDIS_PASS":           "s3cr3t-sentinel-REDIS_PASS",
				"DB_PASS":              "s3cr3t-sentinel-DB_PASS",
				"POSTGRES_PASSWORD":    "s3cr3t-sentinel-POSTGRES_PASSWORD",
				"REDIS_HOST":           "s3cr3t-sentinel-REDIS_HOST",
				"DB_USER":              "s3cr3t-sentinel-DB_USER",
			}
			var buf bytes.Buffer
			writer, flags := log.Writer(), log.Flags()
			t.Cleanup(func() {
				log.SetOutput(writer)
				log.SetFlags(flags)
			})
			log.SetOutput(&buf)
			log.SetFlags(0)

			cfg, err := Load(func(key string) string { return env[key] })

			if buf.Len() != 0 {
				t.Errorf("Load: want no log output, got %q", buf.String())
			}
			if (err != nil) != tc.wantErr {
				t.Errorf("Load: want error %t, got %v", tc.wantErr, err)
			}
			formatted := fmt.Sprintf("%+v", cfg)
			for key, sentinel := range env {
				if key == "CORS_ALLOWED_ORIGINS" {
					continue
				}
				if err != nil && strings.Contains(err.Error(), sentinel) {
					t.Errorf("Load: error contains %s sentinel: %v", key, err)
				}
				if strings.Contains(formatted, sentinel) {
					t.Errorf("Load: formatted config contains %s sentinel: %s", key, formatted)
				}
			}
		})
	}
}
