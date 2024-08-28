package entities

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type Secret string

func (s *Secret) Set(src string) error {
	if len([]byte(src)) < 32 {
		log.Warn().Msg("secret is too short")
	}

	*s = Secret(src)

	return nil
}

func (s Secret) Type() string {
	return "string"
}

func (s Secret) String() string {
	if len(s) <= 2 {
		return string(s)
	}

	masked := strings.Repeat("*", len(s)-2)
	return string(s[0]) + masked + string(s[len(s)-1])
}
