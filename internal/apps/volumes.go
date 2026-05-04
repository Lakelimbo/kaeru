package apps

import (
	"context"
	"time"

	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
)

// (TO-DO)
//
// NewVolume creates an entry of a volume mostly for database purposes.
// Each service can have a volume, and volumes can also be explorable (exposed?)
// on the "x-kaeru" custom tag on the Compose file.
//
// Explorable volumes mean they are available to be accessed on the frontend's
// file explorer (eventually).
func (a *App) NewVolume(id string) error {
	var project database.App
	if err := a.db.Select("project_name").
		From("apps").
		Where(dbx.HashExp{"id": id}).
		One(&project); err != nil {
		logger.Errorf("Unable to find app to expose its volume on db: %v", err)
		return err
	}

	services, err := a.docker.ListContainers(context.Background(), project.ProjectName)
	if err != nil {
		return err
	}

	for _, s := range services.Items {
		for _, mnt := range s.Mounts {
			volume := database.AppVolume{
				ID:        uuid.NewString(),
				AppID:     project.ID,
				Name:      mnt.Name,
				Driver:    mnt.Driver,
				IsExposed: true,
				CreatedAt: time.Now(),
			}

			if err := a.db.Model(&volume).Insert(); err != nil {
				logger.Errorf("Failed to insert app volume into the database: %v", err)
				return err
			}
		}
	}

	return nil
}
