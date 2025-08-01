package undefer_test

import (
	"testing"

	"github.com/mutility/undeferred/undefer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, undefer.Analyzer().Analyzer, "a")
}
