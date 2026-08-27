package logging

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestFieldsAlwaysMarshalable(t *testing.T) {
	f := Fields("k", make(chan int))
	if _, err := json.Marshal(f); err != nil {
		t.Fatalf("Fields produced non-marshalable value: %v", err)
	}
}

func TestEventDoesNotPrintEmptyOnMarshalError(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{base: log.New(&buf, "", 0)}
	l.Event("info", "msg", map[string]any{"k": make(chan int)})
	if strings.TrimSpace(buf.String()) == "" {
		t.Fatal("expected a fallback log line when a field cannot be marshaled")
	}
}
