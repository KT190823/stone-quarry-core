package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"mo-da-backend/internal/database"
)

// session is an in-memory auth token bound to a username with an expiry.
type session struct {
	Username  string
	Name      string
	Role      string
	Dept      string
	Email     string
	ExpiresAt time.Time
}

var (
	sessionStore = struct {
		sync.RWMutex
		m map[string]session
	}{m: map[string]session{}}
)

const sessionTTL = 12 * time.Hour

// HashPassword returns a salted SHA-256 hex digest of the password.
func HashPassword(password string) string {
	h := sha256.Sum256([]byte("mo-da-salt::" + password))
	return hex.EncodeToString(h[:])
}

func matchesPassword(password, hash string) bool {
	if hash == "" {
		return false
	}
	return HashPassword(password) == hash
}

// GenerateToken returns a cryptographically random opaque token.
func GenerateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// Authenticate validates username + password against the users table.
func Authenticate(ctx context.Context, username, password string) (session, error) {
	if username == "" || password == "" {
		return session{}, errors.New("Vui lòng nhập đủ tài khoản và mật khẩu")
	}

	var s session
	var hash, status, role, dept, email, name string
	err := database.Pool.QueryRow(ctx,
		`SELECT username, name, role, dept, email, status, COALESCE(password_hash,'') FROM users WHERE username = $1`,
		username,
	).Scan(&s.Username, &name, &role, &dept, &email, &status, &hash)
	if err != nil {
		return session{}, errors.New("Tài khoản không tồn tại")
	}
	if status == "Khóa" || status == "Inactive" || status == "Bị khóa" {
		return session{}, errors.New("Tài khoản đã bị khóa")
	}
	if !matchesPassword(password, hash) {
		return session{}, errors.New("Sai mật khẩu")
	}

	s.Name = name
	s.Role = role
	s.Dept = dept
	s.Email = email
	s.ExpiresAt = time.Now().Add(sessionTTL)

	// Update last login.
	database.Pool.Exec(ctx, `UPDATE users SET last_login = to_char(NOW(), 'DD/MM/YYYY HH24:MI') WHERE username = $1`, username)

	return s, nil
}

// CreateSession stores a session and returns its token.
func CreateSession(s session) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}
	sessionStore.Lock()
	sessionStore.m[token] = s
	sessionStore.Unlock()
	return token, nil
}

// GetSession resolves a token to a session, or an error if missing/expired.
func GetSession(token string) (session, error) {
	if token == "" {
		return session{}, errors.New("Thiếu mã phiên đăng nhập")
	}
	sessionStore.RLock()
	s, ok := sessionStore.m[token]
	sessionStore.RUnlock()
	if !ok {
		return session{}, errors.New("Phiên đăng nhập không hợp lệ")
	}
	if time.Now().After(s.ExpiresAt) {
		DeleteSession(token)
		return session{}, errors.New("Phiên đăng nhập đã hết hạn")
	}
	return s, nil
}

// DeleteSession invalidates a token.
func DeleteSession(token string) {
	sessionStore.Lock()
	delete(sessionStore.m, token)
	sessionStore.Unlock()
}

// SessionToMap renders a session as a safe API payload.
func SessionToMap(s session) map[string]interface{} {
	return map[string]interface{}{
		"username": s.Username,
		"name":     s.Name,
		"role":     s.Role,
		"dept":     s.Dept,
		"email":    s.Email,
	}
}
