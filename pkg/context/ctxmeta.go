package ctxmeta

import (
	"maps"
	"context"
	"log/slog"

)

type ctxKey string

const contextCarrierKey ctxKey = "ctxmeta-values"

// Set a single key-value pair (returns new context)
func Set(ctx context.Context, key string, value string) context.Context {
	carrier := copyMap(GetAll(ctx))
	carrier[key] = value
	return context.WithValue(ctx, contextCarrierKey, carrier)
}

// SetPair sets multiple key-value pairs (returns new context)
func SetPair(ctx context.Context, keysAndValues ...string) context.Context {
	if len(keysAndValues)%2 != 0 {
		slog.WarnContext(ctx, "[ctxmeta] SetPair called with odd number of arguments — ignoring last one.")
		keysAndValues = keysAndValues[:len(keysAndValues)-1]
	}

	carrier := copyMap(GetAll(ctx))
	for i := 0; i < len(keysAndValues); i += 2 {
		key := keysAndValues[i]
		value := keysAndValues[i+1]
		carrier[key] = value
	}

	return context.WithValue(ctx, contextCarrierKey, carrier)
}

// Get returns a single value
func Get(ctx context.Context, key string) (string, bool) {
	carrier := GetAll(ctx)
	val, ok := carrier[key]
	return val, ok
}

// GetPair fetches multiple keys from the context store
func GetPair(ctx context.Context, keys ...string) map[string]string {
	carrier := GetAll(ctx)
	result := make(map[string]string)
	for _, k := range keys {
		if v, ok := carrier[k]; ok {
			result[k] = v
		}
	}
	return result
}

// GetAll returns all stored key-value pairs
func GetAll(ctx context.Context) map[string]string {
	val := ctx.Value(contextCarrierKey)
	if carrier, ok := val.(map[string]string); ok {
		return carrier
	}
	return make(map[string]string)
}

// Helper: copyMap ensures we don't mutate shared maps
func copyMap(original map[string]string) map[string]string {
	newMap := make(map[string]string, len(original))
	maps.Copy(newMap, original)
	return newMap
}
