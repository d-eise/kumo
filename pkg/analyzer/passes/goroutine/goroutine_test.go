package goroutine_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/sivchari/kumo/pkg/analyzer/passes/goroutine"
)

func TestGoroutineAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goroutine.Analyzer, "goroutine")
}

func TestGoroutineAnalyzerName(t *testing.T) {
	t.Parallel()

	if goroutine.Analyzer.Name != "goroutine" {
		t.Errorf("expected analyzer name %q, got %q", "goroutine", goroutine.Analyzer.Name)
	}
}

func TestGoroutineAnalyzerDoc(t *testing.T) {
	t.Parallel()

	if goroutine.Analyzer.Doc == "" {
		t.Error("expected non-empty doc string for goroutine analyzer")
	}
}
