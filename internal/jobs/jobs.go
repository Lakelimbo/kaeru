package jobs

import (
	"time"

	"github.com/Lakelimbo/kaeru/internal/containers"
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/Lakelimbo/kaeru/tools/misc"
	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
)

// Basically a way to have background jobs.
//
// A job consists of the DB entry and the Pub-Sub messaging for
// realtime updates.
type Job struct {
	db        *dbx.DB
	docker    *containers.DockerRepository
	Messaging *JobMessaging
}

type JobResponse struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Status     database.JobStatus `db:"status"`
	Progress   uint               `json:"progress"`
	Logs       string             `json:"logs"`
	AppID      string             `json:"app_id"`
	CreatedAt  time.Time          `json:"created_at"`
	FinishedAt time.Time          `json:"finished_at"`
}

// Spawn constructs a Job instance.
//
// Decided to call it "Spawn" because Messaging (pubsub) also lives in the same package, to
// avoid naming weirdness.
func Spawn(db *dbx.DB, docker *containers.DockerRepository, messaging *JobMessaging) *Job {
	return &Job{db, docker, messaging}
}

func (j *Job) NewJob(jobType string, appID string) (string, error) {
	id := uuid.NewString()

	job := database.Job{
		ID:        id,
		Type:      jobType,
		Status:    database.JobQueued,
		Progress:  0,
		CreatedAt: time.Now(),
	}
	if appID != "" {
		job.AppID = appID
	}

	if err := j.db.Model(&job).Insert(); err != nil {
		logger.Errorf("Failed to create background job: %v", err)
		return "", err
	}

	return id, nil
}

func (j *Job) GetJob(id string) (*JobResponse, error) {
	var job database.Job
	if err := j.db.Select("*").
		From("jobs").
		Where(dbx.HashExp{"id": id}).
		One(&job); err != nil {
		logger.Errorf("Background job not found: %v", err)
		return nil, err
	}

	res := JobResponse{
		ID:         job.ID,
		Type:       job.Type,
		Status:     job.Status,
		Progress:   job.Progress,
		Logs:       job.Logs,
		AppID:      job.AppID,
		CreatedAt:  job.CreatedAt,
		FinishedAt: job.FinishedAt,
	}

	return &res, nil
}

func (j *Job) ListJobs(status *database.JobStatus, appID string) ([]JobResponse, error) {
	filter := make(dbx.HashExp)
	if status != nil {
		filter["status"] = status
	}

	filter["app_id"] = appID

	var jobs []database.Job
	if err := j.db.Select("*").Where(filter).All(&jobs); err != nil {
		logger.Errorf("Failed to list jobs: %v", err)
		return nil, err
	}

	var res []JobResponse
	for _, j := range jobs {
		res = append(res, JobResponse{
			ID:         j.ID,
			Type:       j.Type,
			Status:     j.Status,
			Progress:   j.Progress,
			Logs:       j.Logs,
			AppID:      j.AppID,
			CreatedAt:  j.CreatedAt,
			FinishedAt: j.FinishedAt,
		})
	}

	return res, nil
}

func (j *Job) UpdateJob(id string, status database.JobStatus, progress int) error {
	var finish time.Time
	if status == database.JobError || status == database.JobSuccess {
		finish = time.Now().UTC()
	}

	update := dbx.Params{
		"status":      status,
		"finished_at": finish,
	}
	if progress != -1 {
		update["progress"] = uint(progress)
	}

	if _, err := j.db.Update("jobs", update, dbx.HashExp{"id": id}).Execute(); err != nil {
		logger.Errorf("Failed to update job: %v", err)
		return err
	}

	eventMsg, err := misc.MapToJSONString(CommonEvent{
		Type:     "status",
		Status:   status,
		Progress: uint(progress),
	})
	if err != nil {
		logger.Errorf("Unable to stream event message: %v", err)
		return err
	}

	j.Messaging.Pub(id, eventMsg)

	return nil
}

func (j *Job) DeleteJob(id string) error {
	if _, err := j.db.Delete("jobs", dbx.HashExp{"id": id}).Execute(); err != nil {
		logger.Errorf("Failed to delete job: %v", err)
		return err
	}

	return nil
}

func (j *Job) Log(id string, level LogLevel, msg string) error {
	msg, err := misc.MapToJSONString(CommonLog{
		Type:    "log",
		Level:   level,
		Message: msg,
	})
	if err != nil {
		logger.Errorf("Unable to publish event log: %v", err)
		return err
	}

	j.Messaging.Pub(id, msg)
	return nil
}
