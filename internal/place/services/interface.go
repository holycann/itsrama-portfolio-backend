package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

// ProvinceService defines operations for managing province-related data
type ProvinceService interface {
	// CRUD Operations
	CreateProvince(ctx context.Context, province *models.Province) (*models.Province, error)
	UpdateProvince(ctx context.Context, province *models.Province) (*models.Province, error)
	DeleteProvince(ctx context.Context, id string) error

	// Read Operations
	GetProvinceByID(ctx context.Context, id string) (*models.ProvinceDTO, error)
	GetProvinceByName(ctx context.Context, name string) (*models.ProvinceDTO, error)
	GetProvinceByCode(ctx context.Context, code string) (*models.ProvinceDTO, error)
	ListProvinces(ctx context.Context, opts base.ListOptions) ([]models.ProvinceDTO, error)
	ListProvincesByRegion(ctx context.Context, region string) ([]models.ProvinceDTO, error)

	// Query Operations
	CountProvinces(ctx context.Context, filters []base.FilterOption) (int, error)
	SearchProvinces(ctx context.Context, opts base.ListOptions) ([]models.ProvinceDTO, int, error)
}

// CityService defines operations for managing city-related data
type CityService interface {
	// CRUD Operations
	CreateCity(ctx context.Context, city *models.City) (*models.City, error)
	UpdateCity(ctx context.Context, city *models.City) (*models.City, error)
	DeleteCity(ctx context.Context, id string) error

	// Read Operations
	GetCityByID(ctx context.Context, id string) (*models.CityDTO, error)
	GetCityByName(ctx context.Context, name string) (*models.CityDTO, error)
	GetCityByCode(ctx context.Context, code string) (*models.CityDTO, error)
	GetCitiesByProvince(ctx context.Context, provinceID string) ([]models.CityDTO, error)
	ListCities(ctx context.Context, opts base.ListOptions) ([]models.CityDTO, int, error)
	ListCitiesByPopulation(ctx context.Context, minPopulation, maxPopulation string) ([]models.CityDTO, error)

	// Query Operations
	CountCities(ctx context.Context, filters []base.FilterOption) (int, error)
	SearchCities(ctx context.Context, opts base.ListOptions) ([]models.CityDTO, int, error)
}

// LocationService defines operations for managing location-related data
type LocationService interface {
	// CRUD Operations
	CreateLocation(ctx context.Context, location *models.LocationCreate) (*models.Location, error)
	UpdateLocation(ctx context.Context, location *models.LocationUpdate) (*models.Location, error)
	DeleteLocation(ctx context.Context, id string) error

	// Read Operations
	GetLocationByID(ctx context.Context, id string) (*models.LocationDTO, error)
	GetLocationByName(ctx context.Context, name string) (*models.LocationDTO, error)
	GetLocationsByCity(ctx context.Context, cityID string) ([]models.LocationDTO, error)
	ListLocations(ctx context.Context, opts base.ListOptions) ([]models.LocationDTO, error)

	// Advanced Search Operations
	GetLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.LocationDTO, error)

	// Query Operations
	CountLocations(ctx context.Context, filters []base.FilterOption) (int, error)
	SearchLocations(ctx context.Context, opts base.ListOptions) ([]models.LocationDTO, error)
}
