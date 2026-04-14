package goadmin

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type sessionState struct {
	UserID uint   `json:"uid"`
	CSRF   string `json:"csrf"`
	Expiry int64  `json:"exp"`
}

func (a *App) readSession(r *http.Request) (*sessionState, error) {
	cookie, err := r.Cookie(a.cfg.SessionCookie)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid session cookie")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(signature, a.sign(payload)) {
		return nil, errors.New("invalid session signature")
	}

	var state sessionState
	if err := json.Unmarshal(payload, &state); err != nil {
		return nil, err
	}
	if time.Now().Unix() > state.Expiry {
		return nil, errors.New("session expired")
	}
	return &state, nil
}

func (a *App) writeSession(w http.ResponseWriter, userID uint) error {
	token, err := randomToken(24)
	if err != nil {
		return err
	}
	state := sessionState{
		UserID: userID,
		CSRF:   token,
		Expiry: time.Now().Add(24 * time.Hour).Unix(),
	}
	payload, err := json.Marshal(state)
	if err != nil {
		return err
	}
	value := base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString(a.sign(payload))
	http.SetCookie(w, &http.Cookie{
		Name:     a.cfg.SessionCookie,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func (a *App) clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     a.cfg.SessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
}

func (a *App) sign(payload []byte) []byte {
	mac := hmac.New(sha256.New, []byte(a.cfg.SessionSecret))
	mac.Write(payload)
	return mac.Sum(nil)
}

func randomToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
