// SQLite doesn't really have enums, but considering the
// enums we use here are usually just integers or specific
// strings for DX purposes, we can just map them like below.

package database

type NotificationLevel uint8

const (
	LevelInfo NotificationLevel = iota
	LevelWarn
	LevelErr
	LevelSuccess
)

type AppServiceVolumeType string

const (
	VolumeVolume AppServiceVolumeType = "volume"
	VolumeBind   AppServiceVolumeType = "bind"
)

type JobStatus string

const (
	JobQueued  JobStatus = "queued"
	JobRunning JobStatus = "running"
	JobSuccess JobStatus = "success"
	JobError   JobStatus = "error"
)
