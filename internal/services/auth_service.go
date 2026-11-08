package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
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
	ClientID  string
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
	var hash, status, role, dept, email, name, clientID string
	err := database.Pool.QueryRow(ctx,
		`SELECT username, name, role, dept, email, status, COALESCE(password_hash,''), COALESCE(client_id,'TTC') FROM users WHERE username = $1`,
		username,
	).Scan(&s.Username, &name, &role, &dept, &email, &status, &hash, &clientID)
	if err != nil {
		return session{}, errors.New("Tài khoản không tồn tại")
	}
	if IsLockedStatus(status) {
		return session{}, errors.New("Tài khoản đã bị khóa")
	}
	if !matchesPassword(password, hash) {
		return session{}, errors.New("Sai mật khẩu")
	}

	s.Name = name
	s.Role = role
	s.Dept = dept
	s.Email = email
	s.ClientID = normalizeClientID(clientID)
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

// IsLockedStatus reports whether an account status means the user may not log in.
func IsLockedStatus(status string) bool {
	statusNorm := strings.ToLower(strings.TrimSpace(status))
	switch statusNorm {
	case "khóa", "khoa", "tạm khóa", "tam khoa", "bị khóa", "bi khoa",
		"inactive", "ngừng hoạt động", "ngung hoat dong", "locked":
		return true
	}
	return false
}

// ChangePassword verifies the current password and stores a new one for the user.
func ChangePassword(ctx context.Context, username, currentPassword, newPassword string) error {
	if strings.TrimSpace(username) == "" {
		return errors.New("Thiếu tài khoản đăng nhập")
	}
	if len(newPassword) < 6 {
		return errors.New("Mật khẩu mới phải có ít nhất 6 ký tự")
	}
	if newPassword == currentPassword {
		return errors.New("Mật khẩu mới phải khác mật khẩu hiện tại")
	}

	var hash string
	err := database.Pool.QueryRow(ctx,
		`SELECT COALESCE(password_hash,'') FROM users WHERE username = $1`, username,
	).Scan(&hash)
	if err != nil {
		return errors.New("Tài khoản không tồn tại")
	}
	if !matchesPassword(currentPassword, hash) {
		return errors.New("Mật khẩu hiện tại không đúng")
	}

	newHash := HashPassword(newPassword)
	if _, err := database.Pool.Exec(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE username = $2`,
		newHash, username,
	); err != nil {
		return errors.New("Không cập nhật được mật khẩu")
	}
	return nil
}

// RevokeSessionsForUser invalidates every active session belonging to a username.
func RevokeSessionsForUser(username string) {
	if username == "" {
		return
	}
	sessionStore.Lock()
	for token, s := range sessionStore.m {
		if s.Username == username {
			delete(sessionStore.m, token)
		}
	}
	sessionStore.Unlock()
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
		"clientId": s.ClientID,
	}
}

// MenuForSession returns the permission-filtered navigation for a session.
func MenuForSession(ctx context.Context, s session) (MenuPayload, error) {
	return GetMenuForRole(ctx, s.Role, s.ClientID)
}
