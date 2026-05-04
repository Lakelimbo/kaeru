package logger_test

import (
	"bytes"
	"strings"
	"testing"

	g "charm.land/log/v2"
	"github.com/Lakelimbo/kaeru/internal/logger"
)

func TestLogLevels(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		input    string
		expected g.Level
	}{
		{"debug", g.DebugLevel},
		{"debug", g.DebugLevel},
		{"info", g.InfoLevel},
		{"warn", g.WarnLevel},
		{"warning", g.WarnLevel},
		{"err", g.ErrorLevel},
		{"error", g.ErrorLevel},
		{"fatal", g.FatalLevel},
		{"unknown", g.InfoLevel},
	}

	for _, s := range scenarios {
		t.Run(s.input, func(t *testing.T) {
			level := logger.ParseLogLevel(s.input)
			if level != s.expected {
				t.Fatalf("expected %v, got %v", s.expected, level)
			}
		})
	}
}

func TestLogSync(t *testing.T) {
	t.Parallel()

	l1 := logger.New()
	l2 := logger.New()

	if l1 == nil {
		t.Fatalf("logger should not be nil")
	}
	if l1 != l2 {
		t.Fatalf("l2 should've returned the same pointer as l1")
	}
}

func TestLoggerPrint(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	charm := g.New(&buf)
	l := &logger.Logger{charm}

	l.Printf("sonic %s", "the hedgehog")

	out := buf.String()
	if !strings.Contains(out, "sonic the hedgehog") {
		t.Fatalf("expected log to contain message, but got %v", out)
	}
}
