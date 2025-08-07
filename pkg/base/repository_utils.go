package base

import (
	"context"
	"reflect"
	"strings"

	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/response"
	"github.com/holycann/cultour-backend/pkg/validator"
)

// ValidateModel performs comprehensive validation for a model
func ValidateModel[T any](model T) error {
	if err := validator.ValidateStruct(model); err != nil {
		return errors.New(
			errors.ErrValidation,
			"Model validation failed",
			err,
			errors.WithContext("validation_errors", err.Error()),
		)
	}
	return nil
}

// BuildFilterFromStruct dynamically creates filter options from a struct
func BuildFilterFromStruct[T any](model T) []FilterOption {
	v := reflect.ValueOf(model)
	t := v.Type()
	var filters []FilterOption

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip zero/empty values
		if IsZero(field.Interface()) {
			continue
		}

		// Use field name as filter key
		filters = append(filters, FilterOption{
			Field:    strings.ToLower(fieldType.Name),
			Operator: OperatorEqual,
			Value:    field.Interface(),
		})
	}

	return filters
}

// SafeExecute provides a wrapper for repository operations with error handling
func SafeExecute[T any, R any](
	ctx context.Context,
	operation func() (R, error),
	errorMessage string,
) (R, error) {
	result, err := operation()
	if err != nil {
		return result, errors.New(
			errors.ErrDatabase,
			errorMessage,
			err,
			errors.WithContext("original_error", err.Error()),
		)
	}
	return result, nil
}

// PaginateResults applies pagination to a slice of results
func PaginateResults[T any](
	results []T,
	page, perPage int,
) ([]T, *response.Pagination) {
	total := len(results)

	// Calculate pagination
	start := (page - 1) * perPage
	end := start + perPage

	if start > total {
		return []T{}, &response.Pagination{
			Total:       total,
			Page:        page,
			PerPage:     perPage,
			TotalPages:  (total + perPage - 1) / perPage,
			HasNextPage: false,
		}
	}

	if end > total {
		end = total
	}

	paginatedResults := results[start:end]

	return paginatedResults, &response.Pagination{
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  (total + perPage - 1) / perPage,
		HasNextPage: end < total,
	}
}
