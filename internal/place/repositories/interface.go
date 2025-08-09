package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

type ProvinceRepository interface {
	base.BaseRepository[models.Province, models.ProvinceDTO]
	FindProvinceByName(ctx context.Context, name string) (*models.ProvinceDTO, error)
	FindProvinceByCode(ctx context.Context, code string) (*models.ProvinceDTO, error)
	ListProvincesByRegion(ctx context.Context, region string) ([]models.ProvinceDTO, error)
}

type CityRepository interface {
	base.BaseRepository[models.City, models.CityDTO]
	FindCitiesByProvince(ctx context.Context, provinceID string) ([]models.CityDTO, error)
	FindCityByName(ctx context.Context, name string) (*models.CityDTO, error)
	FindCityByCode(ctx context.Context, code string) (*models.CityDTO, error)
	ListCitiesByPopulation(ctx context.Context, minPopulation, maxPopulation string) ([]models.CityDTO, error)
}

type LocationRepository interface {
	base.BaseRepository[models.Location, models.LocationDTO]
	FindLocationsByCity(ctx context.Context, cityID string) ([]models.LocationDTO, error)
	FindLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.LocationDTO, error)
	SearchLocationsByName(ctx context.Context, query string) ([]models.LocationDTO, error)
}
