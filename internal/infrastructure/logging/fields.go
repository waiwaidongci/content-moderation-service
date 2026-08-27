// Package implementation for policy-driven content moderation and human review.
package logging

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func Fields(pairs ...any) map[string]any {
	out := map[string]any{}
	for i := 0; i+1 < len(pairs); i += 2 {
		out[fmt.Sprint(pairs[i])] = sanitize(pairs[i+1])
	}
	return out
}
func WithCaller(fields map[string]any) map[string]any {
	copy := map[string]any{}
	for k, v := range fields {
		copy[k] = sanitize(v)
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		copy["caller"] = fmt.Sprintf("%s:%d", file, line)
	}
	return copy
}
func WithDuration(fields map[string]any, start time.Time) map[string]any {
	copy := map[string]any{}
	for k, v := range fields {
		copy[k] = sanitize(v)
	}
	copy["duration_ms"] = time.Since(start).Milliseconds()
	return copy
}

// sanitize replaces values that json.Marshal cannot encode (channels,
// functions, and complex types whose elements contain them) with a string
// placeholder so the resulting field map is always JSON-serializable.
func sanitize(v any) any {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func:
		return fmt.Sprintf("%T(%p)", v, v)
	case reflect.Map:
		out := map[string]any{}
		for _, k := range rv.MapKeys() {
			out[fmt.Sprint(k.Interface())] = sanitize(rv.MapIndex(k).Interface())
		}
		return out
	case reflect.Slice, reflect.Array:
		out := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			out[i] = sanitize(rv.Index(i).Interface())
		}
		return out
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return nil
		}
		return sanitize(rv.Elem().Interface())
	default:
		return v
	}
}

