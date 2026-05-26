package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/sivchari/kumo/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "goroutine", "funclen")
}

func TestAnalyzerGoStmt(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	results := analysistest.Run(t, testdata, analyzer.Analyzer, "goroutine")
	if len(results) == 0 {
		t.Fatal("expected analysis results but got none")
	}
}

func TestAnalyzerFuncDecl(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	results := analysistest.Run(t, testdata, analyzer.Analyzer, "funclen")
	if len(results) == 0 {
		t.Fatal("expected analysis results but got none")
	}
}
