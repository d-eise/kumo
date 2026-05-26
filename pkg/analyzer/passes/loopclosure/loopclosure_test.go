package loopclosure_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/sivchari/kumo/pkg/analyzer/passes/loopclosure"
)

func TestLoopclosureAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, loopclosure.Analyzer, "loopclosure")
}

func TestLoopclosureAnalyzerName(t *testing.T) {
	t.Parallel()

	if loopclosure.Analyzer.Name != "loopclosure" {
		t.Errorf("expected analyzer name %q, got %q", "loopclosure", loopclosure.Analyzer.Name)
	}
}

func TestLoopclosureAnalyzerDoc(t *testing.T) {
	t.Parallel()

	if loopclosure.Analyzer.Doc == "" {
		t.Error("expected non-empty analyzer doc")
	}
}
