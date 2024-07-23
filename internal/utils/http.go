package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func HeadersToStr(headers http.Header) string {
	stringsSlice := []string{}

	for name, values := range headers {
		for _, value := range values {
			stringsSlice = append(stringsSlice, fmt.Sprintf("%s:%s", name, value))
		}
	}

	return strings.Join(stringsSlice, ", ")
}
