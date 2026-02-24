package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env           string // local|staging|prod
	HTTPAddr      string // :8080
	PublicBaseURL string // optional, e.g. https://api.scavo.exchange

	// CORS
	CORSAllowOrigins []string

	// Build info (optional)
	Version string
	Commit  string
}

func LoadFromEnv() (Config, error) {
	c := Config{}

	c.Env = getenv("SCAVO_ENV", "local")
	c.HTTPAddr = getenv("SCAVO_HTTP_ADDR", ":8080")
	c.PublicBaseURL = getenv("SCAVO_PUBLIC_BASE_URL", "")

	allowOrigins := getenv("SCAVO_CORS_ALLOW_ORIGINS", "*")
	c.CORSAllowOrigins = splitCSV(allowOrigins)

	c.Version = getenv("SCAVO_VERSION", "dev")
	c.Commit = getenv("SCAVO_COMMIT", "")

	if !strings.HasPrefix(c.HTTPAddr, ":") && !strings.Contains(c.HTTPAddr, ":") {
		// allow bare port: 8080
		if _, err := strconv.Atoi(c.HTTPAddr); err == nil {
			c.HTTPAddr = ":" + c.HTTPAddr
		} else {
			return c, fmt.Errorf("invalid SCAVO_HTTP_ADDR: %q", c.HTTPAddr)
		}
	}

	return c, nil
}

func getenv(k, def string) string {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	return v
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
