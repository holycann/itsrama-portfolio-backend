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

// MessageHandler menangani permintaan HTTP terkait pesan
type MessageHandler struct {
	messageService services.MessageService
	logger         *logger.Logger
}

// NewMessageHandler membuat instance baru dari message handler
func NewMessageHandler(messageService services.MessageService, logger *logger.Logger) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		logger:         logger,
	}
}

// CreateMessage godoc
// @Summary Membuat pesan baru
// @Description Menambahkan pesan baru ke dalam sistem
// @Tags messages
// @Accept json
// @Produce json
// @Param message body models.Message true "Informasi Pesan"
// @Success 201 {object} response.Response{data=models.ResponseMessage} "Pesan berhasil dibuat"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /messages [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.logger.Error("Error binding message: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.messageService.CreateMessage(c.Request.Context(), &message); err != nil {
		h.logger.Error("Error creating message: %v", err)
		response.InternalServerError(c, "Failed to create message", err.Error())
		return
	}

	response.SuccessCreated(c, message, "Message created successfully")
}

// ListMessages godoc
// @Summary Daftar pesan
// @Description Mendapatkan daftar pesan dengan pagination
// @Tags messages
// @Accept json
// @Produce json
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response{data=[]models.ResponseMessage} "Daftar pesan berhasil didapatkan"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /messages [get]
func (h *MessageHandler) ListMessages(c *gin.Context) {
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

	// Ambil daftar pesan
	messages, err := h.messageService.GetMessages(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving messages: %v", err)
		response.InternalServerError(c, "Failed to retrieve messages", err.Error())
		return
	}

	// Hitung total pesan untuk pagination
	total, err := h.messageService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting messages: %v", err)
		response.InternalServerError(c, "Failed to count messages", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, messages, total, offset/limit+1, limit)
}

// SearchMessages godoc
// @Summary Mencari pesan
// @Description Mencari pesan berdasarkan ID atau kriteria lainnya
// @Tags messages
// @Accept json
// @Produce json
// @Param id query string false "ID Pesan"
// @Param threadID query string false "ID Thread"
// @Param limit query int false "Jumlah data yang dikembalikan" default(10)
// @Param offset query int false "Offset untuk pagination" default(0)
// @Success 200 {object} response.Response{data=[]models.ResponseMessage} "Daftar pesan yang ditemukan"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Pesan tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /messages/search [get]
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	// Ambil parameter query
	id := c.Query("id")
	threadID := c.Query("threadID")
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
		message, err := h.messageService.GetMessageByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error finding message by ID: %v", err)
			response.NotFound(c, "Message not found", err.Error())
			return
		}
		response.SuccessOK(c, message, "Message found")
		return
	}

	// Jika threadID diberikan, cari berdasarkan threadID
	if threadID != "" {
		messages, err := h.messageService.GetMessagesByThreadID(c.Request.Context(), threadID, limit, offset)
		if err != nil {
			h.logger.Error("Error finding messages by thread ID: %v", err)
			response.NotFound(c, "Messages not found", err.Error())
			return
		}

		// Hitung total pesan untuk pagination
		total, err := h.messageService.CountByThreadID(c.Request.Context(), threadID)
		if err != nil {
			h.logger.Error("Error counting messages: %v", err)
			response.InternalServerError(c, "Failed to count messages", err.Error())
			return
		}

		response.WithPagination(c, messages, total, offset/limit+1, limit)
		return
	}

	// Jika tidak ada parameter spesifik, kembalikan daftar pesan
	messages, err := h.messageService.GetMessages(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Error retrieving messages: %v", err)
		response.InternalServerError(c, "Failed to retrieve messages", err.Error())
		return
	}

	// Hitung total pesan untuk pagination
	total, err := h.messageService.Count(c.Request.Context())
	if err != nil {
		h.logger.Error("Error counting messages: %v", err)
		response.InternalServerError(c, "Failed to count messages", err.Error())
		return
	}

	// Gunakan WithPagination untuk menambahkan metadata pagination
	response.WithPagination(c, messages, total, offset/limit+1, limit)
}

// UpdateMessage godoc
// @Summary Memperbarui pesan
// @Description Memperbarui informasi pesan yang sudah ada berdasarkan ID
// @Tags messages
// @Accept json
// @Produce json
// @Param message body models.RequestMessage true "Informasi Pesan yang Diperbarui"
// @Success 200 {object} response.Response{data=models.ResponseMessage} "Pesan berhasil diperbarui"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 404 {object} response.ErrorResponse "Pesan tidak ditemukan"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /messages/{id} [put]
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		h.logger.Error("Error binding message: %v", err)
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.messageService.UpdateMessage(c.Request.Context(), &message); err != nil {
		h.logger.Error("Error updating message: %v", err)
		response.InternalServerError(c, "Failed to update message", err.Error())
		return
	}

	response.SuccessOK(c, message, "Message updated successfully")
}

// DeleteMessage godoc
// @Summary Menghapus pesan
// @Description Menghapus pesan berdasarkan ID
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "ID Pesan"
// @Success 204 {object} response.Response "Pesan berhasil dihapus"
// @Failure 400 {object} response.ErrorResponse "Kesalahan validasi input"
// @Failure 500 {object} response.ErrorResponse "Kesalahan server internal"
// @Router /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Message ID is required", nil)
		return
	}

	if err := h.messageService.DeleteMessage(c.Request.Context(), id); err != nil {
		h.logger.Error("Error deleting message: %v", err)
		response.InternalServerError(c, "Failed to delete message", err.Error())
		return
	}

	response.Success(c, http.StatusNoContent, nil, "Message deleted successfully")
}
