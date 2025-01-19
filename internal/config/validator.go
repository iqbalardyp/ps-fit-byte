package config

import (
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"

	"fit-byte/internal/activity/model"
	"time"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("activity_type", activityTypeValidator)
	validate.RegisterValidation("time_validator", timeValidator)
	validate.RegisterValidation("is_uri", uriValidator)
	return validate
}

var validActivityTypes = map[model.ActivityTypeEnum]struct{}{
	model.ActivityTypeEnumWalking:    {},
	model.ActivityTypeEnumYoga:       {},
	model.ActivityTypeEnumStretching: {},
	model.ActivityTypeEnumCycling:    {},
	model.ActivityTypeEnumSwimming:   {},
	model.ActivityTypeEnumDancing:    {},
	model.ActivityTypeEnumHiking:     {},
	model.ActivityTypeEnumRunning:    {},
	model.ActivityTypeEnumHIIT:       {},
	model.ActivityTypeEnumJumpRope:   {},
}

func activityTypeValidator(fl validator.FieldLevel) bool {
	activity, ok := fl.Field().Interface().(model.ActivityTypeEnum)
	if !ok {
		return false
	}
	_, exists := validActivityTypes[activity]
	return exists
}

func timeValidator(fl validator.FieldLevel) bool {
	t, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Example: Check if time is after 1970-01-01T00:00:00Z
	return !t.IsZero() && t.After(time.Unix(0, 0))
}

func uriValidator(fl validator.FieldLevel) bool {
	uri, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	// Check for empty string
	if strings.TrimSpace(uri) == "" {
		return false
	}

	// Parse the URI
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return false
	}

	// Check if scheme is present and valid
	if parsedURL.Scheme == "" {
		return false
	}

	// Validate host presence for network-based URIs
	if parsedURL.Scheme != "file" && parsedURL.Host == "" {
		return false
	}

	// Additional validation for specific schemes
	switch parsedURL.Scheme {
	case "http", "https":
		// For HTTP(S), ensure there's a valid host
		if !strings.Contains(parsedURL.Host, ".") && parsedURL.Host != "localhost" {
			return false
		}
	case "file":
		// For file scheme, ensure there's a path
		if parsedURL.Path == "" {
			return false
		}
	}

	return true
}
