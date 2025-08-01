// Package main hosts the undeferred analyzer undefer.Analyzer()
package main

import (
	"github.com/mutility/undeferred/undefer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(undefer.Analyzer().Analyzer)
}
