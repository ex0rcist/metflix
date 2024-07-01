package validators

import (
	"errors"
	"regexp"
)

var nameRegexp = regexp.MustCompile(`^[A-Za-z\d]+$`)

func EnsureNamePresent(name string) error {
	if len(name) == 0 {
		return errors.New("missing name")
	}

	return nil
}

func ValidateName(name, kind string) error {
	if !nameRegexp.MatchString(name) {
		return errors.New("invalid name")
	}

	return nil
}

func ValidateKind(kind string) error {
	switch kind {
	case "counter", "gauge":
		return nil

	default:
		return errors.New("invalid metric type")
	}
}
