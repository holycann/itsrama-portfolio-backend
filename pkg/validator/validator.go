package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
)

// ValidationRule defines a custom validation rule
type ValidationRule func(interface{}) error

// ValidationContext provides additional context for validation
type ValidationContext struct {
	Field     string
	Value     interface{}
	Rules     []ValidationRule
	CustomTag string
}

// ValidateStruct performs comprehensive validation for structs or slices of structs
func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)

	// Dereference pointer if needed
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var errors []string

	switch v.Kind() {
	case reflect.Struct:
		errors = append(errors, validateStructValue(v)...)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			// Dereference pointer inside slice if any
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			if item.Kind() != reflect.Struct {
				return fmt.Errorf("ValidateStruct: slice contains non-struct type")
			}
			errors = append(errors, validateStructValue(item)...)
		}
	default:
		return fmt.Errorf("ValidateStruct: unsupported type %s", v.Kind())
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
	}

	return nil
}

// ValidateModel performs comprehensive validation for a model
func ValidateModel[T any](model T) error {
	if err := ValidateStruct(model); err != nil {
		return errors.Wrap(err,
			errors.ErrValidation,
			"Model validation failed",
			errors.WithContext("validation_errors", err.Error()),
		)
	}
	return nil
}

// IsValueChanged checks if the new value is different from the existing one
func IsValueChanged[T any](existing, new *T, ignoredFields ...string) bool {
	if existing == nil || new == nil {
		return true
	}

	existingValue := reflect.ValueOf(existing).Elem()
	newValue := reflect.ValueOf(new).Elem()

	if existingValue.Type() != newValue.Type() {
		return true
	}

	ignoreMap := make(map[string]struct{})
	for _, f := range ignoredFields {
		ignoreMap[f] = struct{}{}
	}

	// Automatically ignore CreatedAt and UpdatedAt fields
	ignoreMap["CreatedAt"] = struct{}{}
	ignoreMap["UpdatedAt"] = struct{}{}

	for i := 0; i < existingValue.NumField(); i++ {
		fieldType := existingValue.Type().Field(i)
		if _, ignored := ignoreMap[fieldType.Name]; ignored {
			continue
		}

		field1 := existingValue.Field(i)
		field2 := newValue.Field(i)

		if !field1.CanInterface() {
			continue
		}

		if !reflect.DeepEqual(field1.Interface(), field2.Interface()) {
			return true
		}
	}

	return false
}

// parseValidationRules parses validation rules from a tag
func parseValidationRules(tag string) []ValidationRule {
	var rules []ValidationRule
	tagParts := strings.Split(tag, ",")

	for _, part := range tagParts {
		switch {
		case part == "required":
			rules = append(rules, requiredRule)
		case strings.HasPrefix(part, "min="):
			minValue := part[4:]
			rules = append(rules, minRule(minValue))
		case strings.HasPrefix(part, "max="):
			maxValue := part[4:]
			rules = append(rules, maxRule(maxValue))
		case part == "email":
			rules = append(rules, emailRule)
		case part == "uuid":
			rules = append(rules, uuidRule)
		case part == "password":
			rules = append(rules, passwordRule)
		}
	}

	return rules
}

// requiredRule checks if a value is not empty
func requiredRule(value interface{}) error {
	if isFieldEmpty(reflect.ValueOf(value)) {
		return fmt.Errorf("is required")
	}
	return nil
}

// minRule creates a rule to check minimum length/value
func minRule(minStr string) ValidationRule {
	return func(value interface{}) error {
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.String:
			min, _ := parseIntValue(minStr)
			if len(rv.String()) < min {
				return fmt.Errorf("must be at least %d characters long", min)
			}
		case reflect.Slice, reflect.Map:
			min, _ := parseIntValue(minStr)
			if rv.Len() < min {
				return fmt.Errorf("must have at least %d items", min)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			min, _ := parseIntValue(minStr)
			if rv.Int() < int64(min) {
				return fmt.Errorf("must be at least %d", min)
			}
		case reflect.Float32, reflect.Float64:
			min, _ := parseFloatValue(minStr)
			if rv.Float() < min {
				return fmt.Errorf("must be at least %f", min)
			}
		}
		return nil
	}
}

// maxRule creates a rule to check maximum length/value
func maxRule(maxStr string) ValidationRule {
	return func(value interface{}) error {
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.String:
			max, _ := parseIntValue(maxStr)
			if len(rv.String()) > max {
				return fmt.Errorf("must be no more than %d characters long", max)
			}
		case reflect.Slice, reflect.Map:
			max, _ := parseIntValue(maxStr)
			if rv.Len() > max {
				return fmt.Errorf("must have no more than %d items", max)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			max, _ := parseIntValue(maxStr)
			if rv.Int() > int64(max) {
				return fmt.Errorf("must be no more than %d", max)
			}
		case reflect.Float32, reflect.Float64:
			max, _ := parseFloatValue(maxStr)
			if rv.Float() > max {
				return fmt.Errorf("must be no more than %f", max)
			}
		}
		return nil
	}
}

// emailRule validates email format
func emailRule(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(str) {
		return fmt.Errorf("must be a valid email address")
	}
	return nil
}

// uuidRule validates UUID format
func uuidRule(value interface{}) error {
	uuidVal, ok := value.(uuid.UUID)
	if !ok {
		return fmt.Errorf("must be a valid UUID")
	}

	if isUUIDNil(uuidVal) {
		return fmt.Errorf("cannot be nil")
	}
	return nil
}

// passwordRule validates password complexity
func passwordRule(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	var (
		hasMinLength   = len(str) >= 8
		hasUppercase   = false
		hasLowercase   = false
		hasNumber      = false
		hasSpecialChar = false
	)

	for _, char := range str {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	var errors []string
	if !hasMinLength {
		errors = append(errors, "at least 8 characters long")
	}
	if !hasUppercase {
		errors = append(errors, "contain at least one uppercase letter")
	}
	if !hasLowercase {
		errors = append(errors, "contain at least one lowercase letter")
	}
	if !hasNumber {
		errors = append(errors, "contain at least one number")
	}
	if !hasSpecialChar {
		errors = append(errors, "contain at least one special character")
	}

	if len(errors) > 0 {
		return fmt.Errorf("must %s", strings.Join(errors, ", "))
	}
	return nil
}

// validateStructValue validates a single struct value and returns slice of error strings
func validateStructValue(v reflect.Value) []string {
	t := v.Type()
	var errs []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		ctx := ValidationContext{
			Field:     fieldType.Name,
			Value:     field.Interface(),
			CustomTag: tag,
			Rules:     parseValidationRules(tag),
		}

		if err := validateField(ctx); err != nil {
			errs = append(errs, err.Error())
		}
	}

	return errs
}

// validateField validates a single field with its rules
func validateField(ctx ValidationContext) error {
	var fieldErrors []string

	for _, rule := range ctx.Rules {
		if err := rule(ctx.Value); err != nil {
			fieldErrors = append(fieldErrors, err.Error())
		}
	}

	if len(fieldErrors) > 0 {
		return fmt.Errorf("%s %s", ctx.Field, strings.Join(fieldErrors, ", "))
	}
	return nil
}

// isUUIDNil checks if a UUID is nil
func isUUIDNil(id uuid.UUID) bool {
	return id.String() == "00000000-0000-0000-0000-000000000000"
}

// isFieldEmpty checks if a field is considered empty
func isFieldEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		// Special handling for time.Time
		if v.Type().String() == "time.Time" {
			return v.Interface().(time.Time).IsZero()
		}
		return false
	default:
		return false
	}
}

// parseIntValue safely parses an integer value
func parseIntValue(val string) (int, error) {
	var result int
	_, err := fmt.Sscanf(val, "%d", &result)
	return result, err
}

// parseFloatValue safely parses a float value
func parseFloatValue(val string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(val, "%f", &result)
	return result, err
}

// ValidateUUID checks if a UUID is valid and not nil
func ValidateUUID(id uuid.UUID, fieldName string) error {
	if isUUIDNil(id) {
		return fmt.Errorf("%s is required and cannot be nil", fieldName)
	}
	return nil
}

// ValidateString checks string fields
func ValidateString(value, fieldName string, minLength, maxLength int) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if minLength > 0 && len(value) < minLength {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, minLength)
	}
	if maxLength > 0 && len(value) > maxLength {
		return fmt.Errorf("%s must be no more than %d characters long", fieldName, maxLength)
	}
	return nil
}
