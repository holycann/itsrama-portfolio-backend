package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type ProvinceService interface {
	CreateProvince(ctx context.Context, province *models.Province) error
	GetProvinceByID(ctx context.Context, id string) (*models.Province, error)
	GetProvinceByName(ctx context.Context, name string) (*models.Province, error)
	ListProvinces(ctx context.Context, opts repository.ListOptions) ([]models.Province, error)
	UpdateProvince(ctx context.Context, province *models.Province) error
	DeleteProvince(ctx context.Context, id string) error
	CountProvinces(ctx context.Context, filters []repository.FilterOption) (int, error)
	SearchProvinces(ctx context.Context, query string, opts repository.ListOptions) ([]models.Province, error)
}

type CityService interface {
	CreateCity(ctx context.Context, city *models.City) error
	GetCityByID(ctx context.Context, id string) (*models.ResponseCity, error)
	GetCityByName(ctx context.Context, name string) (*models.ResponseCity, error)
	GetCitiesByProvince(ctx context.Context, provinceID string) ([]models.ResponseCity, error)
	ListCities(ctx context.Context, opts repository.ListOptions) ([]models.ResponseCity, error)
	UpdateCity(ctx context.Context, city *models.City) error
	DeleteCity(ctx context.Context, id string) error
	CountCities(ctx context.Context, filters []repository.FilterOption) (int, error)
	SearchCities(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseCity, error)
}

type LocationService interface {
	CreateLocation(ctx context.Context, location *models.Location) error
	GetLocationByID(ctx context.Context, id string) (*models.Location, error)
	GetLocationByName(ctx context.Context, name string) (*models.Location, error)
	GetLocationsByCity(ctx context.Context, cityID string) ([]models.Location, error)
	// GetLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.Location, error)
	ListLocations(ctx context.Context, opts repository.ListOptions) ([]models.Location, error)
	UpdateLocation(ctx context.Context, location *models.Location) error
	DeleteLocation(ctx context.Context, id string) error
	CountLocations(ctx context.Context, filters []repository.FilterOption) (int, error)
	SearchLocations(ctx context.Context, query string, opts repository.ListOptions) ([]models.Location, error)
}
