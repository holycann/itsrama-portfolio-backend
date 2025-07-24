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

// LocalStoryHandler menangani permintaan HTTP terkait cerita lokal
type LocalStoryHandler struct {
	localStoryService services.LocalStoryService
	logger            *logger.Logger
}

// NewLocalStoryHandler membuat instance baru dari local story handler
func NewLocalStoryHandler(localStoryService services.LocalStoryService, logger *logger.Logger) *LocalStoryHandler {
	return &LocalStoryHandler{
		localStoryService: localStoryService,
		logger:            logger,
	}
}

// CreateLocalStory godoc
// @Summary Membuat cerita lokal baru
// @Description Menambahkan cerita lokal baru ke dalam sistem
// @Tags local_stories
// @Accept json
// @Produce json
// @Param local_story body models.LocalStory true "Informasi Cerita Lokal"
// @Success 201 {object} response.Response "Cerita lokal berhasil dibuat"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /local-stories [post]
func (h *LocalStoryHandler) CreateLocalStory(c *gin.Context) {
	var localStory models.LocalStory
	if err := c.ShouldBindJSON(&localStory); err != nil {
		h.logger.Error("Error binding local story: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.localStoryService.CreateLocalStory(c.Request.Context(), &localStory); err != nil {
		h.logger.Error("Error creating local story: %v", err)
		response.InternalServerError(c, "Failed to create local story", err.Error())
		return
	}

	response.SuccessCreated(c, localStory, "Local story created successfully")
}

// SearchLocalStories godoc
// @Summary Mencari cerita lokal
// @Description Mencari cerita lokal berdasarkan ID atau judul
// @Tags local_stories
// @Accept json
// @Produce json
// @Param id query string false "ID Cerita Lokal"
// @Param title query string false "Judul Cerita Lokal"
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response "Daftar cerita lokal yang ditemukan"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Cerita lokal tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /local-stories/search [get]
func (h *LocalStoryHandler) SearchLocalStories(c *gin.Context) {
	// Ambil parameter query
	id := c.Query("id")
	title := c.Query("title")
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
		localStory, err := h.localStoryService.GetLocalStoryByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding local story by ID: %v", err)
			response.NotFound(c, "Local story not found", err.Error())
			return
		}
		response.SuccessOK(c, localStory, "Local story found")
		return
	}

	// Jika judul diberikan, cari berdasarkan judul
	if title != "" {
		localStory, err := h.localStoryService.GetLocalStoryByTitle(c.Request.Context(), title)
		if err != nil {
			h.logger.Error("Error finding local story by title: %v", err)
			response.NotFound(c, "Local story not found", err.Error())
			return
		}
		response.SuccessOK(c, localStory, "Local story found")
		return
	}

	// Jika tidak ada parameter spesifik, kembalikan daftar cerita lokal
	localStories, err := h.localStoryService.GetLocalStories(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving local stories: %v", err)
		response.InternalServerError(c, "Failed to retrieve local stories", err.Error())
		return
	}

	// Hitung total cerita lokal untuk pagination
	total, err := h.localStoryService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting local stories: %v", err)
		response.InternalServerError(c, "Failed to count local stories", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, localStories, total, offset/limit+1, limit)
}

// UpdateLocalStory godoc
// @Summary Memperbarui cerita lokal
// @Description Memperbarui informasi cerita lokal yang sudah ada
// @Tags local_stories
// @Accept json
// @Produce json
// @Param local_story body models.LocalStory true "Informasi Cerita Lokal yang Diperbarui"
// @Success 200 {object} response.Response "Cerita lokal berhasil diperbarui"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /local-stories [put]
func (h *LocalStoryHandler) UpdateLocalStory(c *gin.Context) {
	var localStory models.LocalStory
	if err := c.ShouldBindJSON(&localStory); err != nil {
		h.logger.Error("Error binding local story: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.localStoryService.UpdateLocalStory(c.Request.Context(), &localStory); err != nil {
		h.logger.Error("Error updating local story: %v", err)
		response.InternalServerError(c, "Failed to update local story", err.Error())
		return
	}

	response.SuccessOK(c, localStory, "Local story updated successfully")
}

// DeleteLocalStory godoc
// @Summary Menghapus cerita lokal
// @Description Menghapus cerita lokal berdasarkan ID
// @Tags local_stories
// @Accept json
// @Produce json
// @Param id path string true "ID Cerita Lokal"
// @Success 204 {object} response.Response "Cerita lokal berhasil dihapus"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /local-stories/{id} [delete]
func (h *LocalStoryHandler) DeleteLocalStory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Local story ID is required", nil)
		return
	}

	if err := h.localStoryService.DeleteLocalStory(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting local story: %v", err)
		response.InternalServerError(c, "Failed to delete local story", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Local story deleted successfully")
}

// ListLocalStories godoc
// @Summary Mendapatkan daftar cerita lokal
// @Description Mengambil daftar cerita lokal dengan pagination
// @Tags local_stories
// @Accept json
// @Produce json
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response "Daftar cerita lokal berhasil diambil"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /local-stories [get]
func (h *LocalStoryHandler) ListLocalStories(c *gin.Context) {
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

	// Ambil daftar cerita lokal
	localStories, err := h.localStoryService.GetLocalStories(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving local stories: %v", err)
		response.InternalServerError(c, "Failed to retrieve local stories", err.Error())
		return
	}

	// Hitung total cerita lokal untuk pagination
	total, err := h.localStoryService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting local stories: %v", err)
		response.InternalServerError(c, "Failed to count local stories", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, localStories, total, offset/limit+1, limit)
}
