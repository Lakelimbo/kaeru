package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"charm.land/lipgloss/v2"
	g "charm.land/log/v2"
)

// Custom pretty logger, based on Charm's Log.
type Logger struct {
	*g.Logger
}

var (
	logger     *Logger
	initLogger sync.Once
)

func New() *Logger {
	initLogger.Do(func() {
		styles := g.DefaultStyles()

		styles.Levels[g.FatalLevel] = lipgloss.NewStyle().SetString("FATAL").Foreground(lipgloss.Color("1"))
		styles.Levels[g.ErrorLevel] = lipgloss.NewStyle().SetString("ERROR").Foreground(lipgloss.Color("9"))
		styles.Levels[g.WarnLevel] = lipgloss.NewStyle().SetString("WARN ").Foreground(lipgloss.Color("3"))
		styles.Levels[g.InfoLevel] = lipgloss.NewStyle().SetString("INFO").Foreground(lipgloss.Color("4"))
		styles.Levels[g.DebugLevel] = lipgloss.NewStyle().SetString("DEBUG").Foreground(lipgloss.Color("5"))

		styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

		base := g.NewWithOptions(os.Stderr, g.Options{
			ReportCaller:    false,
			ReportTimestamp: true,
			TimeFormat:      fmt.Sprintf("%s |", time.TimeOnly),
		})
		base.SetStyles(styles)

		level := g.InfoLevel
		if envLevel := "debug"; envLevel != "" {
			level = ParseLogLevel(envLevel)
		}

		base.SetLevel(level)
		logger = &Logger{base}
	})

	return logger
}

func ParseLogLevel(level string) g.Level {
	switch strings.ToLower(level) {
	case "debug":
		return g.DebugLevel
	case "info":
		return g.InfoLevel
	case "warn", "warning":
		return g.WarnLevel
	case "err", "error":
		return g.ErrorLevel
	case "fatal":
		return g.FatalLevel
	default:
		return g.InfoLevel
	}
}

func (l *Logger) Printf(fmt string, v ...any) {
	l.Infof(fmt, v...)
}

func (l *Logger) Fatalf(fmt string, v ...any) {
	l.Fatal(fmt, v...)
}
