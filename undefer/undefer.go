// Package undefer implements the undeferred analyzer. It reports when a
// a named parameter is passed to a function that offers a different parameter
// with the same name.

package undefer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = `undefer reports parameters that were likely swapped

While this isn't necessarily a problem, and sometimes is very intentional,
the results of crossing your parameters can be disastrous.`

type undeferAnalyzer struct {
	*analysis.Analyzer
	ExactTypeOnly         bool
	IncludeGeneratedFiles bool
}

func Analyzer() *undeferAnalyzer {
	a := &undeferAnalyzer{
		Analyzer: &analysis.Analyzer{
			Name:     "undefer",
			Doc:      doc,
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		},
	}
	a.Flags.BoolVar(&a.ExactTypeOnly, "exact", false, "suppress undefer reports when types aren't an exact match")
	a.Flags.BoolVar(&a.IncludeGeneratedFiles, "gen", false, "include reports from generated files")

	a.Run = a.run

	return a
}

func (v *undeferAnalyzer) run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	results := make(map[types.Object]struct{})
	inspect.Preorder([]ast.Node{new(ast.FuncType), new(ast.DeferStmt)}, func(n ast.Node) {
		if f, ok := n.(*ast.FuncType); ok {
			clear(results)
			if f.Results == nil {
				return
			}
			for _, a := range f.Results.List {
				for _, n := range a.Names {
					if o := pass.TypesInfo.Defs[n]; o != nil {
						results[o] = struct{}{}
					}
				}
			}
			return
		}
		if len(results) == 0 {
			return
		}
		d := n.(*ast.DeferStmt)
		for _, a := range d.Call.Args {
			switch a := a.(type) {
			case *ast.Ident:
				if obj := pass.TypesInfo.Uses[a]; obj != nil {
					if _, ok := results[obj]; ok {
						pass.Reportf(a.Pos(), "defer captures current value of named result '%s'", a.Name)
					}
				}
			}
		}
	})
	return nil, nil
}
