package utils

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func HeadersToStr(headers http.Header) string {
	stringsSlice := []string{}

	for name, values := range headers {
		for _, value := range values {
			stringsSlice = append(stringsSlice, fmt.Sprintf("%s:%s", name, value))
		}
	}

	sort.Slice(stringsSlice, func(i, j int) bool {
		return stringsSlice[i] < stringsSlice[j]
	})

	return strings.Join(stringsSlice, ", ")
}
