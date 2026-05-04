package apps

import (
	"context"
	"time"

	"github.com/Lakelimbo/kaeru/internal/containers"
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/jobs"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
)

// An App is basically a Compose stack (services). Usually that will be
// a bunch of services, with one of them being user-facing (example:
// Wordpress, where the wordpress service is user-facing/exposed, but it
// also needs MySQL, which is not user-facing).
type App struct {
	db     *dbx.DB
	docker *containers.DockerRepository
	jobs   *jobs.Job
}

func New(db *dbx.DB, docker *containers.DockerRepository, jobs *jobs.Job) *App {
	return &App{db, docker, jobs}
}

type AppRequest struct {
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
	Origin      string `json:"origin"`
	ComposeYAML string `json:"compose_yaml"`
}

type AppResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
	Origin      string `json:"origin,omitempty"`
}

type AppCreateResponse struct {
	JobID string `json:"job_queue_id"`
	AppRequest
}

func (a *App) NewApp(req AppRequest) (*AppCreateResponse, error) {
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.ProjectName, validation.Required, is.LowerCase),
		validation.Field(&req.Origin, is.URL),
	); err != nil {
		logger.Errorf("Invalid payload on app creation: %v", err)
		return nil, err
	}

	app := database.App{
		ID:          uuid.NewString(),
		Name:        req.Name,
		ProjectName: req.ProjectName,
		Origin:      req.Origin,
		ComposeYAML: req.ComposeYAML,
		Creation: database.Creation{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := a.db.Model(&app).Insert(); err != nil {
		logger.Errorf("Failed to create app on database: %v", err)
		return nil, err
	}

	jobID, err := a.jobs.NewJob("app_install", app.ID)
	if err != nil {
		return nil, err
	}

	a.jobs.Log(jobID, jobs.LevelInfo, "Pulling images...")
	// put the creation on a goroutine, so the job is run in the background
	// and doesn't hang.
	go func() {
		ctx := context.Background()

		a.jobs.UpdateJob(jobID, database.JobRunning, 10)
		if err := a.docker.Up(ctx, app); err != nil {
			a.jobs.Log(jobID, jobs.LevelError, err.Error())
			a.jobs.UpdateJob(jobID, database.JobError, -1)
			a.db.Delete("apps", dbx.HashExp{"id": app.ID})
			return
		}

		containers, err := a.docker.ListContainers(ctx, app.ProjectName)
		if err != nil {
			logger.Errorf("Unable to list services for project %s: %v", app.ProjectName, err)
		}

		for _, c := range containers.Items {
			a.NewService(app.ID, c)
		}

		a.jobs.Log(jobID, jobs.LevelInfo, "Containers started")
		a.jobs.UpdateJob(jobID, database.JobSuccess, 100)
	}()

	res := &AppCreateResponse{
		JobID: jobID,
		AppRequest: AppRequest{
			Name:        app.Name,
			ProjectName: app.ProjectName,
			Origin:      app.Origin,
		},
	}

	return res, nil
}

func (a *App) GetApp(id string) (*AppResponse, error) {
	var app database.App
	if err := a.db.Select("*").
		From("apps").
		Where(dbx.HashExp{"id": id}).
		One(&app); err != nil {
		logger.Errorf("Failed to get app: %v", err)
		return nil, err
	}

	res := &AppResponse{
		ID:          app.ID,
		Name:        app.Name,
		ProjectName: app.ProjectName,
		Origin:      app.Origin,
	}

	return res, nil
}

func (a *App) ListApps() ([]AppResponse, error) {
	var apps []database.App
	if err := a.db.Select("*").From("apps").All(&apps); err != nil {
		logger.Errorf("Failed to list apps: %v", err)
		return nil, err
	}

	var res []AppResponse
	for _, app := range apps {
		res = append(res, AppResponse{
			ID:          app.ID,
			Name:        app.Name,
			ProjectName: app.ProjectName,
			Origin:      app.Origin,
		})
	}

	return res, nil
}
