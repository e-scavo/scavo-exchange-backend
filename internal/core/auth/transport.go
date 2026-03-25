package auth

import (
	"net/http"
	"strings"
)

func ExtractBearerToken(headerValue string) string {
	headerValue = strings.TrimSpace(headerValue)
	if headerValue == "" {
		return ""
	}

	if !strings.HasPrefix(strings.ToLower(headerValue), "bearer ") {
		return ""
	}

	return strings.TrimSpace(headerValue[7:])
}

func ExtractTokenFromRequest(r *http.Request, allowQuery bool) string {
	if r == nil {
		return ""
	}

	token := ExtractBearerToken(r.Header.Get("Authorization"))
	if token != "" {
		return token
	}

	if allowQuery {
		token = strings.TrimSpace(r.URL.Query().Get("token"))
		if token != "" {
			return token
		}
	}

	return ""
}
