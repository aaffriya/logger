package ctxmeta

import (
	"context"
)

const (
	TraceIDKey = "trace_id"
	UserIDKey  = "user_id"
	ActionKey  = "action"
)

type ContextData struct {
	TraceID string
	UserID  string
	Action  string
}

// FromContext reads these specific keys from ctxmeta store inside context
func FromContext(ctx context.Context) ContextData {
	if ctx == nil {
		return ContextData{}
	}

	data := ContextData{}
	if traceID, ok := GetPair(ctx, TraceIDKey)[TraceIDKey]; ok {
		data.TraceID = traceID
	}
	if userID, ok := GetPair(ctx, UserIDKey)[UserIDKey]; ok {
		data.UserID = userID
	}
	if action, ok := GetPair(ctx, ActionKey)[ActionKey]; ok {
		data.Action = action
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
