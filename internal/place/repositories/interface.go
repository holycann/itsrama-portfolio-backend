package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type ProvinceRepository interface {
	repository.BaseRepository[models.Province]

	// Specialized methods for provinces
	FindProvinceByName(ctx context.Context, name string) (*models.Province, error)
}

type CityRepository interface {
	repository.BaseRepository[models.City]

	// Specialized methods for cities
	FindCitiesByProvince(ctx context.Context, provinceID string) ([]models.City, error)
	FindCityByName(ctx context.Context, name string) (*models.City, error)
}

type LocationRepository interface {
	repository.BaseRepository[models.Location]

	// Specialized methods for locations
	FindLocationsByCity(ctx context.Context, cityID string) ([]models.Location, error)
	FindLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.Location, error)
}
