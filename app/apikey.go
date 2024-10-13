package app

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"strings"

	"golang.org/x/exp/slog"
)

type ApiKeyConfig struct {
	APIKeyHeader string
	APIKeys      map[string]string
	APIKeyMaxLen int
}

func ApiKeyMiddleware(cfg ApiKeyConfig) (func(handler http.Handler) http.Handler, error) {
	apiKeyHeader := cfg.APIKeyHeader
	apiKeys := cfg.APIKeys
	// apiKeyMaxLen := cfg.APIKeyMaxLen

	decodedAPIKeys := make(map[string][]byte)
	for name, value := range apiKeys {
		decodedKey, err := hex.DecodeString(value)
		if err != nil {
			return nil, err
		}

		decodedAPIKeys[name] = decodedKey
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			apiKey, err := apiToken(r, apiKeyHeader)
			if err != nil {
				slog.Error("request failed API key authentication", "error", err)
				http.Error(w, "invalid API key", http.StatusUnauthorized)
				return
			}

			if _, ok := apiKeyIsValid(apiKey, decodedAPIKeys); !ok {
				hostIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					slog.Error("failed to parse remote address", "error", err)
					hostIP = r.RemoteAddr
				}
				slog.Error("no matching API key found", "remoteIP", hostIP)

				http.Error(w, "invalid api key", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

// apiKeyIsValid checks if the given API key is valid and returns the principal if it is.
func apiKeyIsValid(rawKey string, availableKeys map[string][]byte) (string, bool) {
	hash := sha256.Sum256([]byte(rawKey))
	key := hash[:]

	for name, value := range availableKeys {
		contentEqual := subtle.ConstantTimeCompare(value, key) == 1

		if contentEqual {
			return name, true
		}
	}

	return "", false
}

// bearerToken extracts the content from the header, striping the Bearer prefix
func bearerToken(r *http.Request, header string) (string, error) {
	if header == "" {
		// header = "X-API-KEY"
		header = "Authorization"
	}
	rawToken := r.Header.Get(header)
	pieces := strings.SplitN(rawToken, " ", 2)

	if len(pieces) < 2 {
		return "", errors.New("token with incorrect bearer format")
	}

	token := strings.TrimSpace(pieces[1])

	return token, nil
}

// apiToken extracts the content from the header, striping whitespaces
func apiToken(r *http.Request, header string) (string, error) {
	if header == "" {
		// header = "X-API-KEY"
		header = "Authorization"
	}
	rawToken := r.Header.Get(header)
	token := strings.TrimSpace(rawToken)

	return token, nil
}
