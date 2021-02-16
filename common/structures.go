package common

// SessionData - Session Authorization Info
type SessionData struct {
	SessionID  string `json:"session_id"`
	ValidUntil string `json:"valid_until"`
}
