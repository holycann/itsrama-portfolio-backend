package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type ProvinceRepository interface {
	repository.BaseRepository[models.Province, models.Province]
	FindProvinceByName(ctx context.Context, name string) (*models.Province, error)
}

type CityRepository interface {
	repository.BaseRepository[models.City, models.ResponseCity]
	FindCitiesByProvince(ctx context.Context, provinceID string) ([]models.ResponseCity, error)
	FindCityByName(ctx context.Context, name string) (*models.ResponseCity, error)
}

type LocationRepository interface {
	repository.BaseRepository[models.Location, models.Location]
	FindLocationsByCity(ctx context.Context, cityID string) ([]models.Location, error)
	// FindLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.Location, error)
}
