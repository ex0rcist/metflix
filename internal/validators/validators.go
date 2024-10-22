package validators

import (
	"regexp"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/pkg/metrics"
)

var nameRegexp = regexp.MustCompile(`^[A-Za-z\d]+$`)

// Ensure metric is valid
func ValidateMetric(name, kind string) error {
	if err := validateMetricName(name); err != nil {
		return err
	}

	if err := validateMetricKind(kind); err != nil {
		return err
	}

	return nil
}

func validateMetricName(name string) error {
	if len(name) == 0 {
		return entities.ErrMetricMissingName
	}

	if !nameRegexp.MatchString(name) {
		return entities.ErrMetricInvalidName
	}

	return nil
}

func validateMetricKind(kind string) error {
	switch kind {
	case metrics.KindCounter, metrics.KindGauge:
		return nil

	default:
		return entities.ErrMetricUnknown
	}
}
