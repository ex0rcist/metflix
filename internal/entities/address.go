package entities

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// https://github.com/spf13/pflag/blob/master/README.md#usage
// must comply with https://pkg.go.dev/github.com/spf13/pflag@v1.0.5#Value
type Address string

func (a Address) String() string {
	return string(a)
}

func (a *Address) Set(src string) error {
	chunks := strings.Split(src, ":")
	if len(chunks) != 2 {
		return fmt.Errorf("set address failed: %w", errors.New("bad address"))
	}

	port := chunks[1]

	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("set address failed: %w", errors.New("bad port"))
	}

	*a = Address(src)

	return nil
}

func (a Address) Type() string {
	return "string"
}
