package jobs_test

import (
	"testing"
	"time"

	"github.com/Lakelimbo/kaeru/internal/apps"
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/tools/tests"
)

var dummyCompose string = "x-kaeru:\n  project_name: nx\n\nservices:\n  nc:\n    image: nextcloud:apache\n    restart: always\n    ports:\n      - 80:80\n    volumes:\n      - nc_data:/var/www/html\n    networks:\n      - redisnet\n      - dbnet\n    environment:\n      - REDIS_HOST=redis\n      - MYSQL_HOST=db\n      - MYSQL_DATABASE=nextcloud\n      - MYSQL_USER=nextcloud\n      - MYSQL_PASSWORD=nextcloud\n  redis:\n    image: redis:alpine\n    restart: always\n    networks:\n      - redisnet\n    expose:\n      - 6379\n  db:\n    image: mariadb:10.5\n    command: --transaction-isolation=READ-COMMITTED --binlog-format=ROW\n    restart: always\n    volumes:\n      - db_data:/var/lib/mysql\n    networks:\n      - dbnet\n    environment:\n      - MYSQL_DATABASE=nextcloud\n      - MYSQL_USER=nextcloud\n      - MYSQL_ROOT_PASSWORD=nextcloud\n      - MYSQL_PASSWORD=nextcloud\n    expose:\n      - 3306\nvolumes:\n  db_data:\n  nc_data:\nnetworks:\n  dbnet:\n  redisnet:"

func TestJobCreation(t *testing.T) {
	t.Parallel()

	cfg := tests.TempDB(t.TempDir())
	app := tests.NewTestApp(cfg)

	appManager := apps.New(app.DB, app.Docker, app.Jobs)
	testStack, err := appManager.NewApp(apps.AppRequest{
		Name:        "testx1",
		ProjectName: "testx1",
		ComposeYAML: dummyCompose,
	})
	if err != nil {
		t.Fatalf("failed to create test app: %q", err)
	}

	job, err := app.Jobs.GetJob(testStack.JobID)
	if err != nil {
		t.Fatalf("expected job to be created, got nil")
	}

	if testStack.JobID != job.ID {
		t.Fatalf("expected job ID to be '%s', got '%s'", testStack.JobID, job.ID)
	}

	expectedCreation := time.Now().UTC().Truncate(time.Second)
	gotCreation := job.CreatedAt.UTC().Truncate(time.Second)
	if gotCreation != expectedCreation {
		t.Fatalf("expected job creation to be %s, got %s", expectedCreation, gotCreation)
	}
}

func TestJobUpdate(t *testing.T) {
	t.Parallel()

	cfg := tests.TempDB(t.TempDir())
	app := tests.NewTestApp(cfg)

	appManager := apps.New(app.DB, app.Docker, app.Jobs)
	testStack, err := appManager.NewApp(apps.AppRequest{
		Name:        "testx1",
		ProjectName: "testx1",
		ComposeYAML: dummyCompose,
	})
	if err != nil {
		t.Fatalf("failed to create test app: %q", err)
	}

	job, err := app.Jobs.GetJob(testStack.JobID)
	if err != nil {
		t.Fatalf("expected job to be created, got nil")
	}

	if err := app.Jobs.UpdateJob(job.ID, database.JobSuccess, 100); err != nil {
		t.Fatalf("failed to update job: %q", err)
	}

	updatedJob, _ := app.Jobs.GetJob(testStack.JobID)
	if updatedJob.Status != database.JobSuccess {
		t.Fatalf("expected status to be 'success', got '%s'", job.Status)
	}
	if updatedJob.FinishedAt.IsZero() {
		t.Fatal("expected finished time not not be zero")
	}
}
