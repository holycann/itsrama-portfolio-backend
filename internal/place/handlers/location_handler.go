package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/internal/response"
)

// LocationHandler menangani permintaan HTTP terkait lokasi
type LocationHandler struct {
	locationService services.LocationService
	logger          *logger.Logger
}

// NewLocationHandler membuat instance baru dari location handler
func NewLocationHandler(locationService services.LocationService, logger *logger.Logger) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
		logger:          logger,
	}
}

// CreateLocation godoc
// @Summary Membuat lokasi baru
// @Description Menambahkan lokasi baru ke dalam sistem
// @Tags locations
// @Accept json
// @Produce json
// @Param location body models.Location true "Informasi Lokasi"
// @Success 201 {object} response.Response "Lokasi berhasil dibuat"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		h.logger.Error("Error binding location: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.locationService.CreateLocation(c.Request.Context(), &location); err != nil {
		h.logger.Error("Error creating location: %v", err)
		response.InternalServerError(c, "Failed to create location", err.Error())
		return
	}

	response.SuccessCreated(c, location, "Location created successfully")
}

// SearchLocations godoc
// @Summary Mencari lokasi
// @Description Mencari lokasi berdasarkan ID atau nama
// @Tags locations
// @Accept json
// @Produce json
// @Param id query string false "ID Lokasi"
// @Param name query string false "Nama Lokasi"
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response "Daftar lokasi yang ditemukan"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Lokasi tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /locations/search [get]
func (h *LocationHandler) SearchLocations(c *gin.Context) {
	// Ambil parameter query
	id := c.Query("id")
	name := c.Query("name")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Parsing limit dan offset
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter", err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.BadRequest(c, "Invalid offset parameter", err.Error())
		return
	}

	// Jika ID diberikan, cari berdasarkan ID
	if id != "" {
		location, err := h.locationService.GetLocationByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding location by ID: %v", err)
			response.NotFound(c, "Location not found", err.Error())
			return
		}
		response.SuccessOK(c, location, "Location found")
		return
	}

	// Jika nama diberikan, cari berdasarkan nama
	if name != "" {
		location, err := h.locationService.GetLocationByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error finding location by name: %v", err)
			response.NotFound(c, "Location not found", err.Error())
			return
		}
		response.SuccessOK(c, location, "Location found")
		return
	}

	// Jika tidak ada parameter spesifik, kembalikan daftar lokasi
	locations, err := h.locationService.GetLocations(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving locations: %v", err)
		response.InternalServerError(c, "Failed to retrieve locations", err.Error())
		return
	}

	// Hitung total lokasi untuk pagination
	total, err := h.locationService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting locations: %v", err)
		response.InternalServerError(c, "Failed to count locations", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, locations, total, offset/limit+1, limit)
}

// UpdateLocation godoc
// @Summary Memperbarui lokasi
// @Description Memperbarui informasi lokasi yang sudah ada berdasarkan ID
// @Tags locations
// @Accept json
// @Produce json
// @Param location body models.Location true "Informasi Lokasi yang Diperbarui"
// @Success 200 {object} response.Response "Lokasi berhasil diperbarui"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Lokasi tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /locations/{id} [put]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		h.logger.Error("Error binding location: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.locationService.UpdateLocation(c.Request.Context(), &location); err != nil {
		h.logger.Error("Error updating location: %v", err)
		response.InternalServerError(c, "Failed to update location", err.Error())
		return
	}

	response.SuccessOK(c, location, "Location updated successfully")
}

// DeleteLocation godoc
// @Summary Menghapus lokasi
// @Description Menghapus lokasi berdasarkan ID
// @Tags locations
// @Accept json
// @Produce json
// @Param id path string true "ID Lokasi"
// @Success 204 {object} response.Response "Lokasi berhasil dihapus"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Location ID is required", nil)
		return
	}

	if err := h.locationService.DeleteLocation(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting location: %v", err)
		response.InternalServerError(c, "Failed to delete location", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Location deleted successfully")
}

// ListLocation godoc
// @Summary Mendapatkan daftar lokasi
// @Description Mengambil daftar lokasi dengan pagination
// @Tags locations
// @Accept json
// @Produce json
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response "Daftar kota berhasil diambil"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /locations [get]
func (h *LocationHandler) ListLocation(c *gin.Context) {
	// Ambil parameter query untuk pagination
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Parsing limit dan offset
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		response.BadRequest(c, "Invalid limit parameter", err.Error())
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		response.BadRequest(c, "Invalid offset parameter", err.Error())
		return
	}

	// Ambil daftar kota
	cities, err := h.locationService.GetLocations(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving cities: %v", err)
		response.InternalServerError(c, "Failed to retrieve cities", err.Error())
		return
	}

	// Hitung total kota untuk pagination
	total, err := h.locationService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting cities: %v", err)
		response.InternalServerError(c, "Failed to count cities", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, cities, total, offset/limit+1, limit)
}
