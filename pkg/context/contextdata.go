package ctxmeta

import (
	"context"
)

const (
	TraceIDKey     = "trace_id"
	UserIDKey      = "user_id"
	ActionKey      = "action"
	SessionIDKey   = "session_id"
	TokenKey       = "token"
	SessionDataKey = "session_data"
)

type ContextData struct {
	TraceID     string
	UserID      string
	SessionID   string
	Action      string
	Token       string
	SessionData map[string]any
}

// FromContext reads these specific keys from ctxmeta store inside context
func FromContext(ctx context.Context) ContextData {
	if ctx == nil {
		return ContextData{}
	}

	data := ContextData{}
	allData := GetPair(ctx, TraceIDKey, UserIDKey, SessionIDKey, ActionKey, TokenKey)

	if traceID, ok := allData[TraceIDKey]; ok {
		data.TraceID = traceID
	}
	if userID, ok := allData[UserIDKey]; ok {
		data.UserID = userID
	}
	if sessionID, ok := allData[SessionIDKey]; ok {
		data.SessionID = sessionID
	}
	if action, ok := allData[ActionKey]; ok {
		data.Action = action
	}
	if token, ok := allData[TokenKey]; ok {
		data.Token = token
	}

	// SessionData is stored separately as map[string]any
	if sessionData, ok := GetData(ctx, SessionDataKey); ok {
		if sessionMap, ok := sessionData.(map[string]any); ok {
			data.SessionData = sessionMap
		}
	}

	return data
}

// WithTraceID sets trace_id in ctxmeta context store
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return SetPair(ctx, TraceIDKey, traceID)
}

// WithUserID sets user_id in ctxmeta context store
func WithUserID(ctx context.Context, userID string) context.Context {
	return SetPair(ctx, UserIDKey, userID)
}

// WithAction sets action in ctxmeta context store
func WithAction(ctx context.Context, action string) context.Context {
	return SetPair(ctx, ActionKey, action)
}

// WithSessionID sets session_id in ctxmeta context store
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return SetPair(ctx, SessionIDKey, sessionID)
}

// GetSessionID retrieves session_id from ctxmeta context store
func GetSessionID(ctx context.Context) string {
	sessionID, _ := Get(ctx, SessionIDKey)
	return sessionID
}

// WithToken sets token in ctxmeta context store
func WithToken(ctx context.Context, token string) context.Context {
	return SetPair(ctx, TokenKey, token)
}

// GetToken retrieves token from ctxmeta context store
func GetToken(ctx context.Context) string {
	token, _ := Get(ctx, TokenKey)
	return token
}

// WithSessionData sets session_data in ctxmeta context store
func WithSessionData(ctx context.Context, sessionData map[string]any) context.Context {
	return SetData(ctx, SessionDataKey, sessionData)
}

// GetSessionData retrieves session_data from ctxmeta context store
func GetSessionData(ctx context.Context) map[string]any {
	sessionData, ok := GetData(ctx, SessionDataKey)
	if !ok {
		return nil
	}
	if sessionMap, ok := sessionData.(map[string]any); ok {
		return sessionMap
	}
	return nil
}
