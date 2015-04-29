package rest

import (
	"fmt"
	"crypto/sha256"
	"crypto/hmac"
	"net/http"
	"strconv"
	"strings"
	"time"
	"encoding/hex"
)

var tokenDuration time.Duration = 12 * time.Hour

func getCSRFKey() []byte {
	// TODO: fix
	return []byte("dummykey")
}

func makeCSRFToken(username string, expiration time.Time) string {
	mac_data := fmt.Sprintf("%s:%d", username, expiration.Unix())
	mac := hmac.New(sha256.New, getCSRFKey())
	mac.Write([]byte(mac_data))
	return fmt.Sprintf("%s:%d", hex.EncodeToString(mac.Sum(nil)), expiration.Unix())
}

func NewCSRFToken(r *http.Request) string {
	username, err := getUserEmail(r)
	if err != nil {
		username = ""
	}
	expiration := time.Now().UTC().Add(tokenDuration)
	return makeCSRFToken(username, expiration)
}

func getCSRFTokenFromRequest(r *http.Request) (string, error) {
	if val, ok := r.Header["X-XSRF-TOKEN"]; ok {
		if len(val) == 1 {
			return val[0], nil
		}
		return "", fmt.Errorf("Too many values for header X-XSRF-TOKEN")
	}
	if val, ok := r.Form["xsrf_token"]; ok {
		if len(val) == 1 {
			return val[0], nil
		}
		return "", fmt.Errorf("Too many values for field xsrf_token.")
	}
	return "", fmt.Errorf("No token found.")
}

func ValidateCSRFToken(r *http.Request) bool {
	username, err := getUserEmail(r)
	if err != nil {
		username = ""
	}

	received, err := getCSRFTokenFromRequest(r)
	if err != nil {
		// TODO: log failure
		return false
	}

	components := strings.SplitN(received, ":", 2)
	if len(components) != 2 {
		// Invalid token
		return false
	}

	expiration_ts, err := strconv.ParseInt(components[1], 10, 32)
	if err != nil {
		// Unparsable ts
		return false
	}

	expiration := time.Time.Unix(expiration_ts, 0)
	if expiration.Before(time.Now().UTC()) {
		// Expired token
		return false
	}

	expected := makeCSRFToken(username, expiration)
	return hmac.Equal([]byte(received), []byte(expected))
}

func csrfRequired(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			f(w ,r)
			return
		}
		if !ValidateCSRFToken(r) {
			http.Error(w, "Invalid CSRF Token", 403)
			return
		}
		f(w, r)
	}
}
