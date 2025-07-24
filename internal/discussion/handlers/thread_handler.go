package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/internal/logger"
	"github.com/holycann/cultour-backend/internal/response"
)

// ThreadHandler menangani permintaan HTTP terkait thread
type ThreadHandler struct {
	threadService services.ThreadService
	logger        *logger.Logger
}

// NewThreadHandler membuat instance baru dari thread handler
func NewThreadHandler(threadService services.ThreadService, logger *logger.Logger) *ThreadHandler {
	return &ThreadHandler{
		threadService: threadService,
		logger:        logger,
	}
}

// CreateThread godoc
// @Summary Membuat thread baru
// @Description Menambahkan thread baru ke dalam sistem
// @Tags threads
// @Accept json
// @Produce json
// @Param thread body models.Thread true "Informasi Thread"
// @Success 201 {object} response.Response "Thread berhasil dibuat"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /threads [post]
func (h *ThreadHandler) CreateThread(c *gin.Context) {
	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		h.logger.Error("Error binding thread: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.threadService.CreateThread(c.Request.Context(), &thread); err != nil {
		h.logger.Error("Error creating thread: %v", err)
		response.InternalServerError(c, "Failed to create thread", err.Error())
		return
	}

	response.SuccessCreated(c, thread, "Thread created successfully")
}

// SearchThreads godoc
// @Summary Mencari thread
// @Description Mencari thread berdasarkan ID atau judul
// @Tags threads
// @Accept json
// @Produce json
// @Param id query string false "ID Thread"
// @Param title query string false "Judul Thread"
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response "Daftar thread yang ditemukan"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Thread tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /threads/search [get]
func (h *ThreadHandler) SearchThread(c *gin.Context) {
	// Ambil parameter query
	id := c.Query("id")
	// title := c.Query("title")
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
		thread, err := h.threadService.GetThreadByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding thread by ID: %v", err)
			response.NotFound(c, "Thread not found", err.Error())
			return
		}
		response.SuccessOK(c, thread, "Thread found")
		return
	}

	// Jika tidak ada parameter spesifik, kembalikan daftar thread
	threads, err := h.threadService.GetThreads(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving threads: %v", err)
		response.InternalServerError(c, "Failed to retrieve threads", err.Error())
		return
	}

	// Hitung total thread untuk pagination
	total := len(threads)

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, threads, total, offset/limit+1, limit)
}

// UpdateThread godoc
// @Summary Memperbarui thread
// @Description Memperbarui informasi thread yang sudah ada berdasarkan ID
// @Tags threads
// @Accept json
// @Produce json
// @Param thread body models.Thread true "Informasi Thread yang Diperbarui"
// @Success 200 {object} response.Response "Thread berhasil diperbarui"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Thread tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /threads/{id} [put]
func (h *ThreadHandler) UpdateThread(c *gin.Context) {
	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		h.logger.Error("Error binding thread: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.threadService.UpdateThread(c.Request.Context(), &thread); err != nil {
		h.logger.Error("Error updating thread: %v", err)
		response.InternalServerError(c, "Failed to update thread", err.Error())
		return
	}

	response.SuccessOK(c, thread, "Thread updated successfully")
}

// ListThreads godoc
// @Summary Mendapatkan daftar thread
// @Description Mengambil daftar thread dengan pagination
// @Tags threads
// @Accept json
// @Produce json
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response "Daftar thread berhasil diambil"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /threads [get]
func (h *ThreadHandler) ListThreads(c *gin.Context) {
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

	// Ambil daftar thread
	threads, err := h.threadService.GetThreads(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving threads: %v", err)
		response.InternalServerError(c, "Failed to retrieve threads", err.Error())
		return
	}

	// Hitung total thread untuk pagination
	total := len(threads)

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, threads, total, offset/limit+1, limit)
}

// DeleteThread godoc
// @Summary Menghapus thread
// @Description Menghapus thread berdasarkan ID
// @Tags threads
// @Accept json
// @Produce json
// @Param id path string true "ID Thread"
// @Success 204 {object} response.Response "Thread berhasil dihapus"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /threads/{id} [delete]
func (h *ThreadHandler) DeleteThread(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Thread ID is required", nil)
		return
	}

	if err := h.threadService.DeleteThread(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting thread: %v", err)
		response.InternalServerError(c, "Failed to delete thread", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Thread deleted successfully")
}
