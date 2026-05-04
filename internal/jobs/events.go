package jobs

import "github.com/Lakelimbo/kaeru/internal/database"

type EventType string

const (
	TypeLog    EventType = "log"
	TypeStatus EventType = "status"
)

type CommonEvent struct {
	Type     EventType          `json:"string"`
	Status   database.JobStatus `json:"status"`
	Progress uint               `json:"progress"`
}

type LogLevel uint8

const (
	LevelInfo LogLevel = iota
	LevelError
)

type CommonLog struct {
	Type    EventType `json:"string"`
	Level   LogLevel  `json:"level"`
	Message string    `json:"message"`
}
