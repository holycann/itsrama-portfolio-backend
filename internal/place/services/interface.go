package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
)

type CityService interface {
	CreateCity(ctx context.Context, location *models.City) error
	GetCities(ctx context.Context, limit, offset int) ([]*models.City, error)
	GetCityByID(ctx context.Context, id string) (*models.City, error)
	GetCityByName(ctx context.Context, name string) (*models.City, error)
	UpdateCity(ctx context.Context, user *models.City) error
	DeleteCity(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type LocationService interface {
	CreateLocation(ctx context.Context, location *models.Location) error
	GetLocations(ctx context.Context, limit, offset int) ([]*models.Location, error)
	GetLocationByID(ctx context.Context, id string) (*models.Location, error)
	GetLocationByName(ctx context.Context, name string) (*models.Location, error)
	UpdateLocation(ctx context.Context, user *models.Location) error
	DeleteLocation(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type ProvinceService interface {
	CreateProvince(ctx context.Context, province *models.Province) error
	GetProvinces(ctx context.Context, limit, offset int) ([]*models.Province, error)
	GetProvinceByID(ctx context.Context, id string) (*models.Province, error)
	GetProvinceByName(ctx context.Context, name string) (*models.Province, error)
	UpdateProvince(ctx context.Context, province *models.Province) error
	DeleteProvince(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}
