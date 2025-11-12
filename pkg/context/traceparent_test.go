package ctxmeta

import (
	"context"
	"testing"
)

func TestParseTraceparent_Valid(t *testing.T) {
	hdr := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
	tc, err := ParseTraceparent(hdr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tc.Version != "00" || tc.TraceID != "0af7651916cd43dd8448eb211c80319c" || tc.ParentID != "b7ad6b7169203331" || tc.TraceFlags != "01" {
		t.Fatalf("parsed values mismatch: %+v", tc)
	}
}

func TestParseTraceparent_Invalid(t *testing.T) {
	cases := []string{
		"bad-value",
		"00-zz-00-01",
		"00-0af7651916cd43dd8448eb211c80319c-badspan-01",
		"00-00000000000000000000000000000000-b7ad6b7169203331-01", // all zero trace id
		"00-0af7651916cd43dd8448eb211c80319c-0000000000000000-01", // all zero parent id
	}
	for _, hdr := range cases {
		if _, err := ParseTraceparent(hdr); err == nil {
			t.Errorf("expected error for %q", hdr)
		}
	}
}

func TestWithTraceparentStoresFields(t *testing.T) {
	ctx := context.Background()
	ctx, err := WithTraceparent(ctx, "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	data := FromContext(ctx)
	if data.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" || data.SpanID != "00f067aa0ba902b7" || data.TraceFlags != "00" {
		t.Fatalf("stored values mismatch: %+v", data)
	}
}

func TestGenerateTraceparentFromContext_UsesExistingTraceIDAndRegensParent(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "0af7651916cd43dd8448eb211c80319c")
	ctx = WithTraceFlags(ctx, "01")

	ctx, hdr1, err := GenerateTraceparentFromContext(ctx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	ctx, hdr2, err := GenerateTraceparentFromContext(ctx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	tc1, err := ParseTraceparent(hdr1)
	if err != nil {
		t.Fatalf("hdr1 parse error: %v", err)
	}
	tc2, err := ParseTraceparent(hdr2)
	if err != nil {
		t.Fatalf("hdr2 parse error: %v", err)
	}
	if tc1.TraceID != GetTraceID(ctx) || tc2.TraceID != GetTraceID(ctx) {
		t.Fatalf("trace id mismatch: tc1=%s, tc2=%s, expected=%s", tc1.TraceID, tc2.TraceID, GetTraceID(ctx))
	}
	if tc1.ParentID == tc2.ParentID {
		t.Fatalf("expected different parent IDs, got same: %s", tc1.ParentID)
	}
}

func TestGenerateTraceparentFromContext_GeneratesNewTraceIDWhenMissing(t *testing.T) {
	ctx := context.Background()
	ctx, hdr, err := GenerateTraceparentFromContext(ctx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	tc, err := ParseTraceparent(hdr)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if tc.TraceID == "" || tc.ParentID == "" {
		t.Fatalf("expected generated trace and parent ids, got: %+v", tc)
	}
	if GetTraceID(ctx) == "" || GetSpanID(ctx) == "" {
		t.Fatalf("expected ids saved in context")
	}
}

func TestGenerateTraceparent_Standalone(t *testing.T) {
	hdr1, err := GenerateTraceparent()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	hdr2, err := GenerateTraceparent()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if hdr1 == hdr2 {
		t.Fatalf("expected different headers for each call, got same: %s", hdr1)
	}
	tc1, err := ParseTraceparent(hdr1)
	if err != nil {
		t.Fatalf("parse error for hdr1: %v", err)
	}
	tc2, err := ParseTraceparent(hdr2)
	if err != nil {
		t.Fatalf("parse error for hdr2: %v", err)
	}
	if tc1.TraceID == "" || tc1.ParentID == "" {
		t.Fatalf("expected valid IDs in hdr1: %+v", tc1)
	}
	if tc2.TraceID == "" || tc2.ParentID == "" {
		t.Fatalf("expected valid IDs in hdr2: %+v", tc2)
	}
	if tc1.TraceFlags != "01" || tc2.TraceFlags != "01" {
		t.Fatalf("expected default flags '01', got: %s, %s", tc1.TraceFlags, tc2.TraceFlags)
	}
}

func TestGenerateTraceparentWithTraceID_UsesProvidedTraceID(t *testing.T) {
	traceID := "0af7651916cd43dd8448eb211c80319c"
	hdr1, err := GenerateTraceparentWithTraceID(traceID, "01")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	hdr2, err := GenerateTraceparentWithTraceID(traceID, "01")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	tc1, err := ParseTraceparent(hdr1)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	tc2, err := ParseTraceparent(hdr2)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if tc1.TraceID != traceID || tc2.TraceID != traceID {
		t.Fatalf("expected traceID %s, got: %s, %s", traceID, tc1.TraceID, tc2.TraceID)
	}
	if tc1.ParentID == tc2.ParentID {
		t.Fatalf("expected different parent IDs, got same: %s", tc1.ParentID)
	}
}
