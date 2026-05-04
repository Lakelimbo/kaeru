// Set the database models here, with the `db:` attributes for
// DBX.
//
// I found it better to place this here, because placing them in
// their own respective service (e.g: User in /users) ended up
// needing to rename them, and things got messy when a single service
// had multiple tables, or different methods needed tables from each
// other, causing circular-dependencies (not allowed in Go).
//
// As described in the DBX documentation, set the TableName() function
// as a method for a table so DBX can properly get the table name (if it
// wasn't clear already, on the DB the tables are in plural, because they're
// a collection, but the structs are representing an independent unit, so they're
// in singular).

package database

import (
	"time"
)

// Not a table, just a wrapper so we don't need to
// define CreatedAt and UpdatedAt on every table
// (though of course this doesn't cover situations like
// FinishedAt, or if there's only CreatedAt, etc.).
type Creation struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type User struct {
	ID            string `db:"id"`
	Username      string `db:"username"`
	FormattedName string `db:"formatted_name"`
	Email         string `db:"email"`
	Avatar        string `db:"avatar"`
	Password      string `db:"password"`
	Creation
}

func (User) TableName() string {
	return "users"
}

type UserAccess struct {
	UserID    string    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	IssuedAt  time.Time `db:"issued_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (UserAccess) TableName() string {
	return "user_access"
}

type Notification struct {
	ID             string            `db:"id"`
	UserID         string            `db:"user_id"`
	Category       string            `db:"category"`
	Title          string            `db:"title"`
	Description    string            `db:"description"`
	OnClickPath    string            `db:"onclick_path"`
	Image          string            `db:"image"`
	Level          NotificationLevel `db:"level"`
	DismissTimeout int               `db:"dismiss_timeout"`
	IsRead         bool              `db:"is_read"`
	Metadata       map[string]any    `db:"metadata"`
	CreatedAt      time.Time         `db:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}

type App struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	ProjectName string `db:"project_name"`
	Origin      string `db:"origin"`
	ComposeYAML string `db:"compose_yaml"`
	Creation
}

func (App) TableName() string {
	return "apps"
}

type AppVolume struct {
	ID        string    `db:"id"`
	AppID     string    `db:"app_id"`
	Name      string    `db:"name"`
	Driver    string    `db:"driver"`
	IsExposed bool      `db:"is_exposed"`
	CreatedAt time.Time `db:"created_at"`
}

func (AppVolume) TableName() string {
	return "app_volumes"
}

type AppService struct {
	AppID           string `db:"app_id"`
	Name            string `db:"name"`
	Image           string `db:"image"`
	LastKnownStatus string `db:"last_known_status"`
	ContainerID     string `db:"container_id"`
}

func (AppService) TableName() string {
	return "app_services"
}

type AppServiceVolume struct {
	ServiceID   string               `db:"service_id"`
	ServiceName string               `db:"service_name"`
	VolumeID    string               `db:"volume_id"`
	MountPath   string               `db:"mount_path"`
	Readonly    bool                 `db:"readonly"`
	Type        AppServiceVolumeType `db:"type"`
	BindSource  string               `db:"bind_source"`
}

func (AppServiceVolumeType) TableName() string {
	return "app_service_volumes"
}

type AppConfig struct {
	ID          string `db:"id"`
	AppID       string `db:"app_id"`
	Key         string `db:"key"`
	Value       string `db:"value"`
	IsSensitive bool   `db:"is_sensitive"`
	Creation
}

func (AppConfig) TableName() string {
	return "app_configs"
}

type Job struct {
	ID         string    `db:"id"`
	Type       string    `db:"type"`
	Status     JobStatus `db:"status"`
	Progress   uint      `db:"progress"`
	Logs       string    `db:"logs"`
	AppID      string    `db:"app_id"`
	CreatedAt  time.Time `db:"created_at"`
	FinishedAt time.Time `db:"finished_at"`
}

func (Job) TableName() string {
	return "jobs"
}
