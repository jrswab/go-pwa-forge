package config

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// @anchor:config:load environment → process config including CORS origins
// Config is the main configuration for the app
type Config struct {
	RedisHost string
	RedisPass string

	PostgresUser string
	PostgresPass string
	PostgresHost string
	PostgresPort string
	PostgresDB   string

	// CORSAllowedOrigins are exact origins permitted for app-route CORS.
	// Empty means same-origin / no CORS.
	CORSAllowedOrigins []string
}

// Load loads the config from the environment. Invalid CORS origin tokens
// return an error (strict reject; never a permissive wildcard fallback).
func (c *Config) Load() error {
	redisHost := os.Getenv("REDIS_HOST")
	c.RedisHost = redisHost
	if c.RedisHost == "" {
		log.Println("REDIS_HOST env var is blank", c)
	}

	redisPass := os.Getenv("REDIS_PASS")
	c.RedisPass = redisPass
	if c.RedisPass == "" {
		log.Println("REDIS_PASS env var is blank", c)
	}

	c.PostgresUser = os.Getenv("DB_USER")
	if c.PostgresUser == "" {
		log.Println("DB_USER env var is blank", c)
	}

	c.PostgresPass = os.Getenv("DB_PASS")
	if c.PostgresPass == "" {
		log.Println("DB_PASS env var is blank", c)
	}

	c.PostgresHost = os.Getenv("DB_HOST")
	if c.PostgresHost == "" {
		log.Println("DB_HOST env var is blank", c)
	}

	c.PostgresPort = os.Getenv("DB_PORT")
	if c.PostgresPort == "" {
		log.Println("DB_PORT env var is blank", c)
	}

	c.PostgresDB = os.Getenv("DB_NAME")
	if c.PostgresDB == "" {
		log.Println("DB_NAME env var is blank", c)
	}

	origins, err := parseCORSOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if err != nil {
		return err
	}
	c.CORSAllowedOrigins = origins
	return nil
}

// ValidateCORSOrigins checks caller-supplied origins with the same strict
// rules as environment parsing. Empty/nil is valid (no CORS). Any wildcard
// or malformed origin is an error — never a permissive policy.
func ValidateCORSOrigins(origins []string) error {
	for i, origin := range origins {
		if err := validateOrigin(origin); err != nil {
			return fmt.Errorf("cors origin[%d]: %w", i, err)
		}
	}
	return nil
}

// parseCORSOrigins splits a comma-separated origin list. Empty/whitespace-only
// whole input yields no CORS. When the whole input is non-empty, every
// comma-separated element must be a non-empty valid origin (leading, trailing,
// or repeated commas are rejected). Each origin must be an absolute http(s)
// origin (scheme + host [+ port], no path/query/fragment) with no '*'.
func parseCORSOrigins(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS: empty origin entry")
		}
		if err := validateOrigin(p); err != nil {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS: %w", err)
		}
		out = append(out, p)
	}
	return out, nil
}

func validateOrigin(origin string) error {
	if origin == "" {
		return fmt.Errorf("origin must not be empty")
	}
	// Reject '*' anywhere (bare wildcard, chi match-all, or host patterns).
	if strings.Contains(origin, "*") {
		return fmt.Errorf("wildcard origin %q is not allowed", origin)
	}
	u, err := url.Parse(origin)
	if err != nil {
		return fmt.Errorf("invalid origin %q: %w", origin, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid origin %q: scheme must be http or https", origin)
	}
	if u.Host == "" {
		return fmt.Errorf("invalid origin %q: missing host", origin)
	}
	// Exact origins only: reject path/query/fragment beyond scheme://host[:port]
	if u.Path != "" && u.Path != "/" {
		return fmt.Errorf("invalid origin %q: must not include a path", origin)
	}
	if u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" || u.User != nil {
		return fmt.Errorf("invalid origin %q: must be scheme://host[:port] only", origin)
	}
	// Canonical form is without trailing slash
	if strings.HasSuffix(origin, "/") {
		return fmt.Errorf("invalid origin %q: must not include a path", origin)
	}
	hostname := u.Hostname()
	if hostname == "" {
		return fmt.Errorf("invalid origin %q: missing host", origin)
	}
	port := u.Port()
	if port != "" {
		n, err := strconv.Atoi(port)
		// deliberate: reject 0 and leading zeros; origins use 1–65535 canonical decimal
		if err != nil || n < 1 || n > 65535 || port != strconv.Itoa(n) {
			return fmt.Errorf("invalid origin %q: invalid port", origin)
		}
	}
	// Require Host to be the exact canonical hostname[/IPv6 brackets][+port] form.
	// Catches empty ports ("host:"), multi-port ("host:443:44"), and bare IPv6.
	if u.Host != canonicalOriginHost(hostname, port) {
		return fmt.Errorf("invalid origin %q: invalid host", origin)
	}
	if ip := net.ParseIP(hostname); ip != nil {
		return nil
	}
	if !validDNSHostname(hostname) {
		return fmt.Errorf("invalid origin %q: invalid host", origin)
	}
	return nil
}

// canonicalOriginHost builds host or host:port, bracketing IPv6 literals.
func canonicalOriginHost(hostname, port string) string {
	host := hostname
	if strings.Contains(hostname, ":") {
		host = "[" + hostname + "]"
	}
	if port != "" {
		return host + ":" + port
	}
	return host
}

// validDNSHostname accepts ASCII DNS labels (RFC 1123): letters, digits, hyphens;
// no empty labels, no leading/trailing hyphen per label.
func validDNSHostname(host string) bool {
	if host == "" || len(host) > 253 {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if label == "" || len(label) > 63 {
			return false
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for i := 0; i < len(label); i++ {
			c := label[i]
			ok := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-'
			if !ok {
				return false
			}
		}
	}
	return true
}
