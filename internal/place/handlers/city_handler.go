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

// CityHandler menangani permintaan HTTP terkait lokasi
type CityHandler struct {
	cityService services.CityService
	logger      *logger.Logger
}

// NewCityHandler membuat instance baru dari city handler
func NewCityHandler(cityService services.CityService, logger *logger.Logger) *CityHandler {
	return &CityHandler{
		cityService: cityService,
		logger:      logger,
	}
}

// CreateCity godoc
// @Summary Membuat lokasi baru
// @Description Menambahkan lokasi baru ke dalam sistem
// @Tags cities
// @Accept json
// @Produce json
// @Param city body models.City true "Informasi Lokasi"
// @Success 201 {object} response.Response "Lokasi berhasil dibuat"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /cities [post]
func (h *CityHandler) CreateCity(c *gin.Context) {
	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		h.logger.Error("Error binding city: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.cityService.CreateCity(c.Request.Context(), &city); err != nil {
		h.logger.Error("Error creating city: %v", err)
		response.InternalServerError(c, "Failed to create city", err.Error())
		return
	}

	response.SuccessCreated(c, city, "City created successfully")
}

// SearchCitys godoc
// @Summary Mencari lokasi
// @Description Mencari lokasi berdasarkan ID atau nama
// @Tags cities
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
// @Router /cities/search [get]
func (h *CityHandler) SearchCity(c *gin.Context) {
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
		city, err := h.cityService.GetCityByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding city by ID: %v", err)
			response.NotFound(c, "City not found", err.Error())
			return
		}
		response.SuccessOK(c, city, "City found")
		return
	}

	// Jika nama diberikan, cari berdasarkan nama
	if name != "" {
		city, err := h.cityService.GetCityByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error finding city by name: %v", err)
			response.NotFound(c, "City not found", err.Error())
			return
		}
		response.SuccessOK(c, city, "City found")
		return
	}

	// Jika tidak ada parameter spesifik, kembalikan daftar lokasi
	cities, err := h.cityService.GetCities(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving cities: %v", err)
		response.InternalServerError(c, "Failed to retrieve cities", err.Error())
		return
	}

	// Hitung total lokasi untuk pagination
	total, err := h.cityService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting cities: %v", err)
		response.InternalServerError(c, "Failed to count cities", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, cities, total, offset/limit+1, limit)
}

// UpdateCity godoc
// @Summary Memperbarui kota
// @Description Memperbarui informasi kota yang sudah ada berdasarkan ID
// @Tags cities
// @Accept json
// @Produce json
// @Param city body models.City true "Informasi Kota yang Diperbarui"
// @Success 200 {object} response.Response "Kota berhasil diperbarui"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Kota tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /cities/{id} [put]
func (h *CityHandler) UpdateCity(c *gin.Context) {
	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		h.logger.Error("Error binding city: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.cityService.UpdateCity(c.Request.Context(), &city); err != nil {
		h.logger.Error("Error updating city: %v", err)
		response.InternalServerError(c, "Failed to update city", err.Error())
		return
	}

	response.SuccessOK(c, city, "City updated successfully")
}

// DeleteCity godoc
// @Summary Menghapus lokasi
// @Description Menghapus lokasi berdasarkan ID
// @Tags cities
// @Accept json
// @Produce json
// @Param id path string true "ID Lokasi"
// @Success 204 {object} response.Response "Lokasi berhasil dihapus"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /cities/{id} [delete]
func (h *CityHandler) DeleteCity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "City ID is required", nil)
		return
	}

	if err := h.cityService.DeleteCity(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting city: %v", err)
		response.InternalServerError(c, "Failed to delete city", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "City deleted successfully")
}

// ListCities godoc
// @Summary Mendapatkan daftar kota
// @Description Mengambil daftar kota dengan pagination
// @Tags cities
// @Accept json
// @Produce json
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response "Daftar kota berhasil diambil"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /cities [get]
func (h *CityHandler) ListCities(c *gin.Context) {
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
	cities, err := h.cityService.GetCities(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving cities: %v", err)
		response.InternalServerError(c, "Failed to retrieve cities", err.Error())
		return
	}

	// Hitung total kota untuk pagination
	total, err := h.cityService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting cities: %v", err)
		response.InternalServerError(c, "Failed to count cities", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, cities, total, offset/limit+1, limit)
}
