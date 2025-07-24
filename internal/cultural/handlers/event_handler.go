package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
)

// EventHandler menangani permintaan HTTP terkait lokasi
type EventHandler struct {
	eventService services.EventService
	logger       *logger.Logger
}

// NewEventHandler membuat instance baru dari event handler
func NewEventHandler(eventService services.EventService, logger *logger.Logger) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		logger:       logger,
	}
}

// CreateEvent godoc
// @Summary Membuat lokasi baru
// @Description Menambahkan lokasi baru ke dalam sistem
// @Tags events
// @Accept json
// @Produce json
// @Param event body models.Event true "Informasi Lokasi"
// @Success 201 {object} response.Response "Lokasi berhasil dibuat"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		h.logger.Error("Error binding event: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.eventService.CreateEvent(c.Request.Context(), &event); err != nil {
		h.logger.Error("Error creating event: %v", err)
		response.InternalServerError(c, "Failed to create event", err.Error())
		return
	}

	response.SuccessCreated(c, event, "Event created successfully")
}

// SearchEvents godoc
// @Summary Mencari lokasi
// @Description Mencari lokasi berdasarkan ID atau nama
// @Tags events
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
// @Router /events/search [get]
func (h *EventHandler) SearchEvents(c *gin.Context) {
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
		event, err := h.eventService.GetEventByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding event by ID: %v", err)
			response.NotFound(c, "Event not found", err.Error())
			return
		}
		response.SuccessOK(c, event, "Event found")
		return
	}

	// Jika nama diberikan, cari berdasarkan nama
	if name != "" {
		event, err := h.eventService.GetEventByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error finding event by name: %v", err)
			response.NotFound(c, "Event not found", err.Error())
			return
		}
		response.SuccessOK(c, event, "Event found")
		return
	}

	// Jika tidak ada parameter spesifik, kembalikan daftar lokasi
	events, err := h.eventService.GetEvents(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving events: %v", err)
		response.InternalServerError(c, "Failed to retrieve events", err.Error())
		return
	}

	// Hitung total lokasi untuk pagination
	total, err := h.eventService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting events: %v", err)
		response.InternalServerError(c, "Failed to count events", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, events, total, offset/limit+1, limit)
}

// UpdateEvent godoc
// @Summary Memperbarui lokasi
// @Description Memperbarui informasi lokasi yang sudah ada berdasarkan ID
// @Tags events
// @Accept json
// @Produce json
// @Param event body models.Event true "Informasi Lokasi yang Diperbarui"
// @Success 200 {object} response.Response "Lokasi berhasil diperbarui"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Lokasi tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /events/{id} [put]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		h.logger.Error("Error binding event: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.eventService.UpdateEvent(c.Request.Context(), &event); err != nil {
		h.logger.Error("Error updating event: %v", err)
		response.InternalServerError(c, "Failed to update event", err.Error())
		return
	}

	response.SuccessOK(c, event, "Event updated successfully")
}

// DeleteEvent godoc
// @Summary Menghapus lokasi
// @Description Menghapus lokasi berdasarkan ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "ID Lokasi"
// @Success 204 {object} response.Response "Lokasi berhasil dihapus"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Event ID is required", nil)
		return
	}

	if err := h.eventService.DeleteEvent(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting event: %v", err)
		response.InternalServerError(c, "Failed to delete event", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Event deleted successfully")
}

// ListEvent godoc
// @Summary Mendapatkan daftar lokasi
// @Description Mengambil daftar lokasi dengan pagination
// @Tags events
// @Accept json
// @Produce json
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response "Daftar kota berhasil diambil"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /events [get]
func (h *EventHandler) ListEvent(c *gin.Context) {
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
	cities, err := h.eventService.GetEvents(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving cities: %v", err)
		response.InternalServerError(c, "Failed to retrieve cities", err.Error())
		return
	}

	// Hitung total kota untuk pagination
	total, err := h.eventService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting cities: %v", err)
		response.InternalServerError(c, "Failed to count cities", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, cities, total, offset/limit+1, limit)
}
