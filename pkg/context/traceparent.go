package ctxmeta

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// TraceContext represents W3C trace context information
type TraceContext struct {
	Version    string // 2 hex characters, currently "00"
	TraceID    string // 32 hex characters (16 bytes), must not be all zeros
	ParentID   string // 16 hex characters (8 bytes), must not be all zeros (aka span id)
	TraceFlags string // 2 hex characters (1 byte)
}

// ParseTraceparent parses a W3C traceparent header
// Format: version-trace-id-parent-id-trace-flags
// Example: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
func ParseTraceparent(traceparent string) (*TraceContext, error) {
	parts := strings.Split(traceparent, "-")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid traceparent format: expected 4 parts, got %d", len(parts))
	}

	tc := &TraceContext{
		Version:    parts[0],
		TraceID:    strings.ToLower(parts[1]),
		ParentID:   strings.ToLower(parts[2]),
		TraceFlags: strings.ToLower(parts[3]),
	}

	if err := validateTraceContext(tc); err != nil {
		return nil, err
	}
	return tc, nil
}

// String returns the traceparent in W3C format
func (tc *TraceContext) String() string {
	return fmt.Sprintf("%s-%s-%s-%s", tc.Version, tc.TraceID, tc.ParentID, tc.TraceFlags)
}

// IsSampled returns true if the trace is sampled (should be recorded)
func (tc *TraceContext) IsSampled() bool {
	return tc.TraceFlags == "01"
}

// validateTraceContext ensures fields match W3C constraints
func validateTraceContext(tc *TraceContext) error {
	// Version: must be "00" (only version currently supported)
	if tc.Version != "00" {
		return errors.New("invalid version in traceparent (only '00' is supported)")
	}
	if !isValidHexLen(tc.TraceID, 16) || isAllZeroHex(tc.TraceID) {
		return errors.New("invalid trace-id in traceparent")
	}
	if !isValidHexLen(tc.ParentID, 8) || isAllZeroHex(tc.ParentID) {
		return errors.New("invalid parent-id in traceparent")
	}
	if !isValidHexLen(tc.TraceFlags, 1) {
		return errors.New("invalid trace-flags in traceparent")
	}
	return nil
}

// isValidHexLen checks if s is hex string of exactly byteLen bytes
func isValidHexLen(s string, byteLen int) bool {
	if len(s) != byteLen*2 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func isAllZeroHex(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '0' {
			return false
		}
	}
	return true
}

// generateID creates a random ID with n bytes and returns a lowercase hex string
func generateID(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Avoid the (extremely unlikely) all-zero ID to comply with W3C
	allZero := true
	for _, v := range b {
		if v != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		// regenerate once
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
	}
	return strings.ToLower(hex.EncodeToString(b)), nil
}

// NewTraceContext constructs a new TraceContext using provided traceID (optional).
// If traceID is empty, a new one is generated. ParentID is always (re)generated.
// TraceFlags defaults to "01" (sampled) if empty.
func NewTraceContext(traceID, traceFlags string) (*TraceContext, error) {
	var err error
	if traceID == "" {
		traceID, err = generateID(16)
		if err != nil {
			return nil, err
		}
	} else {
		// validate provided traceID
		if !isValidHexLen(strings.ToLower(traceID), 16) || isAllZeroHex(strings.ToLower(traceID)) {
			return nil, errors.New("invalid provided traceID")
		}
	}
	parentID, err := generateID(8)
	if err != nil {
		return nil, err
	}
	if traceFlags == "" {
		traceFlags = "01"
	} else {
		if !isValidHexLen(strings.ToLower(traceFlags), 1) {
			return nil, errors.New("invalid provided traceFlags")
		}
		traceFlags = strings.ToLower(traceFlags)
	}

	return &TraceContext{
		Version:    "00",
		TraceID:    strings.ToLower(traceID),
		ParentID:   parentID,
		TraceFlags: traceFlags,
	}, nil
}

// GenerateTraceparent generates a new W3C traceparent string without context.
// Creates new traceID, parentID/spanID, and uses default flags "01" (sampled).
// Returns the full traceparent header string: "00-{traceID}-{spanID}-{flags}"
func GenerateTraceparent() (string, error) {
	tc, err := NewTraceContext("", "01")
	if err != nil {
		return "", err
	}
	return tc.String(), nil
}

// GenerateTraceparentWithTraceID generates a W3C traceparent string using the provided traceID.
// Creates a new parentID/spanID. TraceFlags defaults to "01" (sampled) if empty.
// Returns the full traceparent header string.
func GenerateTraceparentWithTraceID(traceID, traceFlags string) (string, error) {
	tc, err := NewTraceContext(traceID, traceFlags)
	if err != nil {
		return "", err
	}
	return tc.String(), nil
}

// WithTraceparent parses the header and stores trace_id, span_id(parent), and trace_flags in context
func WithTraceparent(ctx context.Context, header string) (context.Context, error) {
	tc, err := ParseTraceparent(header)
	if err != nil {
		return ctx, err
	}
	ctx = SetPair(ctx, TraceIDKey, tc.TraceID, SpanIDKey, tc.ParentID, TraceFlagsKey, tc.TraceFlags)
	return ctx, nil
}

// GetTraceparent builds a traceparent string from values in context.
// Returns the traceparent string and true if successful; returns false if trace_id or span_id is missing.
// If false is returned, the caller should generate a new traceparent using GenerateTraceparent or similar.
func GetTraceparent(ctx context.Context) (string, bool) {
	vals := GetPair(ctx, TraceIDKey, SpanIDKey, TraceFlagsKey)
	traceID, ok := vals[TraceIDKey]
	if !ok || traceID == "" {
		return "", false
	}
	parentID := vals[SpanIDKey]
	if parentID == "" {
		// span_id not present; cannot fabricate here. Caller should generate a new traceparent.
		return "", false
	}
	flags := vals[TraceFlagsKey]
	if flags == "" {
		flags = "01"
	}
	return fmt.Sprintf("00-%s-%s-%s", strings.ToLower(traceID), strings.ToLower(parentID), strings.ToLower(flags)), true
}

// GetOrGenerateTraceID returns a context that guarantees a trace_id stored and returns it.
func GetOrGenerateTraceID(ctx context.Context) (context.Context, string, error) {
	if tid, ok := Get(ctx, TraceIDKey); ok && tid != "" {
		return ctx, tid, nil
	}
	tid, err := generateID(16)
	if err != nil {
		return ctx, "", err
	}
	ctx = SetPair(ctx, TraceIDKey, tid)
	return ctx, tid, nil
}

// GenerateTraceparentFromContext ensures a trace_id exists (use or generate), always generates a new parent(span) id,
// uses existing trace_flags if present (else defaults to "01"), stores span_id and flags into context, and returns the header string.
func GenerateTraceparentFromContext(ctx context.Context) (context.Context, string, error) {
	var err error
	ctx, tid, err := GetOrGenerateTraceID(ctx)
	if err != nil {
		return ctx, "", err
	}
	pid, err := generateID(8)
	if err != nil {
		return ctx, "", err
	}
	flags := GetTraceFlags(ctx)
	if flags == "" {
		flags = "01"
	}
	ctx = SetPair(ctx, SpanIDKey, pid, TraceFlagsKey, strings.ToLower(flags))
	header := fmt.Sprintf("00-%s-%s-%s", strings.ToLower(tid), pid, strings.ToLower(flags))
	return ctx, header, nil
}
