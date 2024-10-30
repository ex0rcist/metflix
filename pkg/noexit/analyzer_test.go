package noexit_test

import (
	"testing"

	"github.com/ex0rcist/metflix/pkg/noexit"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), noexit.Analyzer, "./...")
}
