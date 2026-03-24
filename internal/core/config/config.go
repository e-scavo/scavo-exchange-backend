package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env           string
	HTTPAddr      string
	PublicBaseURL string

	CORSAllowOrigins []string

	Version string
	Commit  string

	// JWT
	JWTSecret string
	JWTIssuer string
	JWTTTLHrs int

	// PostgreSQL
	PostgresURL string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Readiness requirements
	ReadinessRequirePostgres bool
	ReadinessRequireRedis    bool
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

	c.JWTSecret = getenv("SCAVO_JWT_SECRET", "dev_dev_dev_dev_dev_dev_dev_dev")
	c.JWTIssuer = getenv("SCAVO_JWT_ISSUER", "scavo-exchange-backend")

	ttlH := getenv("SCAVO_JWT_TTL_HOURS", "24")
	n, _ := strconv.Atoi(ttlH)
	if n <= 0 {
		n = 24
	}
	c.JWTTTLHrs = n

	c.PostgresURL = getenv("SCAVO_POSTGRES_URL", "")

	c.RedisAddr = getenv("SCAVO_REDIS_ADDR", "")
	c.RedisPassword = getenv("SCAVO_REDIS_PASSWORD", "")

	redisDB, _ := strconv.Atoi(getenv("SCAVO_REDIS_DB", "0"))
	if redisDB < 0 {
		redisDB = 0
	}
	c.RedisDB = redisDB

	c.ReadinessRequirePostgres = getenvBool("SCAVO_READINESS_REQUIRE_POSTGRES", false)
	c.ReadinessRequireRedis = getenvBool("SCAVO_READINESS_REQUIRE_REDIS", false)

	if !strings.HasPrefix(c.HTTPAddr, ":") && !strings.Contains(c.HTTPAddr, ":") {
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

func getenvBool(k string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(k)))
	if v == "" {
		return def
	}

	switch v {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
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
