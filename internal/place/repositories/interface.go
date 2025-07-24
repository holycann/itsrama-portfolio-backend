package repositories

import (
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type CityRepository interface {
	repository.BaseRepository[models.City]
}

type LocationRepository interface {
	repository.BaseRepository[models.Location]
}
