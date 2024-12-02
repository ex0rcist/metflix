package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get runtime caller")
	}

	dir := filepath.Dir(filename)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (go.mod not found)")
		}

		dir = parent
	}
}
