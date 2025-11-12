package ctxmeta

import (
	"context"
	"testing"
)

func TestSet(t *testing.T) {
	ctx := context.Background()
	ctx = Set(ctx, "test_key", "test_value")

	value, ok := Get(ctx, "test_key")
	if !ok {
		t.Error("Expected key to exist in context")
	}
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
}

func TestSetPair(t *testing.T) {
	ctx := context.Background()
	ctx = SetPair(ctx, "key1", "value1", "key2", "value2")

	value1, ok1 := Get(ctx, "key1")
	value2, ok2 := Get(ctx, "key2")

	if !ok1 || !ok2 {
		t.Error("Expected both keys to exist in context")
	}
	if value1 != "value1" || value2 != "value2" {
		t.Errorf("Expected 'value1' and 'value2', got '%s' and '%s'", value1, value2)
	}
}

func TestSetPairOddArguments(t *testing.T) {
	ctx := context.Background()
	// This should ignore the last argument
	ctx = SetPair(ctx, "key1", "value1", "key2", "value2", "orphan")

	value1, ok1 := Get(ctx, "key1")
	value2, ok2 := Get(ctx, "key2")

	if !ok1 || !ok2 {
		t.Error("Expected both keys to exist in context")
	}
	if value1 != "value1" || value2 != "value2" {
		t.Errorf("Expected 'value1' and 'value2', got '%s' and '%s'", value1, value2)
	}
}

func TestGetPair(t *testing.T) {
	ctx := context.Background()
	ctx = SetPair(ctx, "key1", "value1", "key2", "value2", "key3", "value3")

	result := GetPair(ctx, "key1", "key3", "nonexistent")

	if len(result) != 2 {
		t.Errorf("Expected 2 keys in result, got %d", len(result))
	}
	if result["key1"] != "value1" {
		t.Errorf("Expected 'value1' for key1, got '%s'", result["key1"])
	}
	if result["key3"] != "value3" {
		t.Errorf("Expected 'value3' for key3, got '%s'", result["key3"])
	}
	if _, exists := result["nonexistent"]; exists {
		t.Error("Expected nonexistent key to not be in result")
	}
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	ctx = SetPair(ctx, "key1", "value1", "key2", "value2")

	all := GetAll(ctx)

	if len(all) != 2 {
		t.Errorf("Expected 2 keys in GetAll result, got %d", len(all))
	}
	if all["key1"] != "value1" || all["key2"] != "value2" {
		t.Errorf("Expected correct values in GetAll result, got %v", all)
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = SetPair(ctx, "trace_id", "0af7651916cd43dd8448eb211c80319c", "span_id", "b7ad6b7169203331", "trace_flags", "01", "user_id", "user-456", "action", "TEST_ACTION")

	data := FromContext(ctx)

	if data.TraceID != "0af7651916cd43dd8448eb211c80319c" {
		t.Errorf("Expected TraceID '0af7651916cd43dd8448eb211c80319c', got '%s'", data.TraceID)
	}
	if data.SpanID != "b7ad6b7169203331" {
		t.Errorf("Expected SpanID 'b7ad6b7169203331', got '%s'", data.SpanID)
	}
	if data.TraceFlags != "01" {
		t.Errorf("Expected TraceFlags '01', got '%s'", data.TraceFlags)
	}
	if data.UserID != "user-456" {
		t.Errorf("Expected UserID 'user-456', got '%s'", data.UserID)
	}
	if data.Action != "TEST_ACTION" {
		t.Errorf("Expected Action 'TEST_ACTION', got '%s'", data.Action)
	}
}

func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	// Test with full W3C traceparent format
	ctx = WithTraceID(ctx, "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")

	data := FromContext(ctx)
	if data.TraceID != "0af7651916cd43dd8448eb211c80319c" {
		t.Errorf("Expected TraceID '0af7651916cd43dd8448eb211c80319c', got '%s'", data.TraceID)
	}
	if data.SpanID != "b7ad6b7169203331" {
		t.Errorf("Expected SpanID 'b7ad6b7169203331', got '%s'", data.SpanID)
	}
	if data.TraceFlags != "01" {
		t.Errorf("Expected TraceFlags '01', got '%s'", data.TraceFlags)
	}
}

func TestWithTraceIDBackwardCompatibility(t *testing.T) {
	ctx := context.Background()
	// Test with just a trace ID (backward compatibility)
	ctx = WithTraceID(ctx, "trace-789")

	data := FromContext(ctx)
	if data.TraceID != "trace-789" {
		t.Errorf("Expected TraceID 'trace-789', got '%s'", data.TraceID)
	}
	if data.SpanID != "" {
		t.Errorf("Expected empty SpanID, got '%s'", data.SpanID)
	}
	if data.TraceFlags != "" {
		t.Errorf("Expected empty TraceFlags, got '%s'", data.TraceFlags)
	}
}

func TestWithUserID(t *testing.T) {
	ctx := context.Background()
	ctx = WithUserID(ctx, "user-999")

	data := FromContext(ctx)
	if data.UserID != "user-999" {
		t.Errorf("Expected UserID 'user-999', got '%s'", data.UserID)
	}
}

func TestWithAction(t *testing.T) {
	ctx := context.Background()
	ctx = WithAction(ctx, "LOGIN_ACTION")

	data := FromContext(ctx)
	if data.Action != "LOGIN_ACTION" {
		t.Errorf("Expected Action 'LOGIN_ACTION', got '%s'", data.Action)
	}
}

func TestFromContextEmpty(t *testing.T) {
	ctx := context.Background()

	data := FromContext(ctx)
	if data.TraceID != "" || data.UserID != "" || data.Action != "" {
		t.Errorf("Expected empty ContextData, got %+v", data)
	}
}

func TestFromContextNil(t *testing.T) {
	data := FromContext(context.TODO())
	if data.TraceID != "" || data.UserID != "" || data.Action != "" {
		t.Errorf("Expected empty ContextData for nil context, got %+v", data)
	}
}

func TestCopyMap(t *testing.T) {
	original := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	copied := copyMap(original)

	// Modify original
	original["key3"] = "value3"

	// Copied should not be affected
	if len(copied) != 2 {
		t.Errorf("Expected copied map to have 2 elements, got %d", len(copied))
	}
	if _, exists := copied["key3"]; exists {
		t.Error("Expected copied map to not contain key3")
	}
}

func TestChainedOperations(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithUserID(ctx, "user-456")
	ctx = WithAction(ctx, "CHAINED_ACTION")
	ctx = Set(ctx, "custom_key", "custom_value")

	data := FromContext(ctx)
	customValue, ok := Get(ctx, "custom_key")

	if data.TraceID != "trace-123" {
		t.Errorf("Expected TraceID 'trace-123', got '%s'", data.TraceID)
	}
	if data.UserID != "user-456" {
		t.Errorf("Expected UserID 'user-456', got '%s'", data.UserID)
	}
	if data.Action != "CHAINED_ACTION" {
		t.Errorf("Expected Action 'CHAINED_ACTION', got '%s'", data.Action)
	}
	if !ok || customValue != "custom_value" {
		t.Errorf("Expected custom_key 'custom_value', got '%s' (exists: %v)", customValue, ok)
	}
}

func TestGetTraceID(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")

	traceID := GetTraceID(ctx)
	if traceID != "0af7651916cd43dd8448eb211c80319c" {
		t.Errorf("Expected TraceID '0af7651916cd43dd8448eb211c80319c', got '%s'", traceID)
	}
}

func TestGetSpanID(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")

	spanID := GetSpanID(ctx)
	if spanID != "b7ad6b7169203331" {
		t.Errorf("Expected SpanID 'b7ad6b7169203331', got '%s'", spanID)
	}
}

func TestGetTraceFlags(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")

	traceFlags := GetTraceFlags(ctx)
	if traceFlags != "01" {
		t.Errorf("Expected TraceFlags '01', got '%s'", traceFlags)
	}
}

func TestWithSpanID(t *testing.T) {
	ctx := context.Background()
	ctx = WithSpanID(ctx, "b7ad6b7169203331")

	data := FromContext(ctx)
	if data.SpanID != "b7ad6b7169203331" {
		t.Errorf("Expected SpanID 'b7ad6b7169203331', got '%s'", data.SpanID)
	}
}

func TestWithTraceFlags(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceFlags(ctx, "01")

	data := FromContext(ctx)
	if data.TraceFlags != "01" {
		t.Errorf("Expected TraceFlags '01', got '%s'", data.TraceFlags)
	}
}
