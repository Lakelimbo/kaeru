package apps

import (
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/moby/moby/api/types/container"
)

// (TO-DO)
//
// NewService creates an app's service, mostly for database purposes.
// A service is technically also an application, but not the whole stack brought
// from the Compose file.
//
// For example, we may have MediaWiki, which consists of itself, a database like
// MariaDB, and also a cache like Redis; each one is a service, but under the umbrella
// of the app "MediaWiki".
func (a *App) NewService(appID string, req container.Summary) error {
	service := database.AppService{
		AppID:           appID,
		Name:            req.Names[0],
		Image:           req.Image,
		LastKnownStatus: req.Status,
		ContainerID:     req.ID,
	}

	if err := a.db.Model(&service).Insert(); err != nil {
		logger.Errorf("Failed to create service for app ID %s: %v", appID, err)
		return err
	}

	if err := a.NewVolume(appID); err != nil {
		return err
	}

	return nil
}
