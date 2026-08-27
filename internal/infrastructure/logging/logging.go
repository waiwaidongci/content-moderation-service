// Package implementation for policy-driven content moderation and human review.
package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Logger struct{ base *log.Logger }

func New() *Logger { return &Logger{base: log.Default()} }
func (l *Logger) Event(level, msg string, fields map[string]any) {
	if fields == nil {
		fields = map[string]any{}
	}
	fields["level"] = level
	fields["message"] = msg
	fields["ts"] = time.Now().UTC().Format(time.RFC3339Nano)
	b, err := json.Marshal(fields)
	if err != nil {
		b, _ = json.Marshal(map[string]any{
			"level":   level,
			"message": msg,
			"ts":      fields["ts"],
			"error":   fmt.Sprintf("log field serialization failed: %v", err),
		})
	}
	l.base.Print(string(b))
}
func (l *Logger) Info(msg string, fields map[string]any)  { l.Event("info", msg, fields) }
func (l *Logger) Error(msg string, fields map[string]any) { l.Event("error", msg, fields) }
