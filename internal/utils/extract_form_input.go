package utils

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
)

// ExtractFileHeaders handles extracting file headers from a multipart form request
func ExtractFileHeaders(c *gin.Context, fieldName string, maxSizeInMB int64) ([]*multipart.FileHeader, error) {
	var fileHeaders []*multipart.FileHeader

	// First, try to get files from c.Request.MultipartForm
	if c.Request.MultipartForm != nil && len(c.Request.MultipartForm.File[fieldName]) > 0 {
		fileHeaders = c.Request.MultipartForm.File[fieldName]
	} else {
		// Then try to get files from c.MultipartForm()
		form, _ := c.MultipartForm()
		if form != nil && len(form.File[fieldName]) > 0 {
			fileHeaders = form.File[fieldName]
		} else {
			// fallback to single file
			f, err := c.FormFile(fieldName)
			if err == nil {
				fileHeaders = []*multipart.FileHeader{f}
			}
		}
	}

	// Validate file sizes
	for _, fileHeader := range fileHeaders {
		if fileHeader.Size > maxSizeInMB*1024*1024 {
			return nil, errors.New(
				errors.ErrValidation,
				fmt.Sprintf("%s files must be less than %dMB each", fieldName, maxSizeInMB),
				nil,
			)
		}
	}

	return fileHeaders, nil
}

// ExtractFormDataPayload handles extracting and unmarshaling a JSON payload from a form with validation
func ExtractFormDataPayload[T any](c *gin.Context, payload *T) error {
	// Extract JSON payload
	jsonPayload := c.PostForm("payload")
	if jsonPayload == "" {
		return errors.New(
			errors.ErrValidation,
			"Payload is required",
			nil,
		)
	}

	// Unmarshal JSON payload
	if err := json.Unmarshal([]byte(jsonPayload), payload); err != nil {
		return errors.New(
			errors.ErrValidation,
			"Invalid payload format: "+err.Error(),
			err,
		)
	}

	return nil
}
