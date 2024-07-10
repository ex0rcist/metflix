package validators

import (
	"regexp"

	"github.com/ex0rcist/metflix/internal/entities"
)

var nameRegexp = regexp.MustCompile(`^[A-Za-z\d]+$`)

func EnsureNamePresent(name string) error {
	if len(name) == 0 {
		return entities.ErrMetricMissingName
	}

	return nil
}

func ValidateName(name string) error {
	if !nameRegexp.MatchString(name) {
		return entities.ErrMetricInvalidName
	}

	return nil
}

func ValidateKind(kind string) error {
	switch kind {
	case "counter", "gauge":
		return nil

	default:
		return entities.ErrMetricUnknown
	}
}
