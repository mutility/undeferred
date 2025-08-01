// Package undefer implements the undeferred analyzer. It reports when a
// a named parameter is passed to a function that offers a different parameter
// with the same name.

package undefer

import (
	"go/ast"
	"slices"

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
	inspect.WithStack([]ast.Node{new(ast.DeferStmt)}, func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		if !push {
			return true
		}
		d := n.(*ast.DeferStmt)
		for _, a := range d.Call.Args {
			switch a := a.(type) {
			case *ast.Ident:
				if obj := pass.TypesInfo.Uses[a]; obj != nil {
					if slices.ContainsFunc(stack, func(n ast.Node) bool {
						if f := funcOf(n); f != nil && f.Results != nil {
							return slices.ContainsFunc(f.Results.List, func(r *ast.Field) bool {
								return slices.ContainsFunc(r.Names, func(id *ast.Ident) bool {
									return pass.TypesInfo.Defs[id] == obj
								})
							})
						}
						return false
					}) {
						pass.Reportf(a.Pos(), "defer captures current value of named result '%s'", a.Name)
					}
				}
			}
		}
		return true
	})
	return nil, nil
}

func funcOf(n ast.Node) *ast.FuncType {
	switch n := n.(type) {
	case *ast.FuncDecl:
		return n.Type
	case *ast.FuncType:
		return n
	case *ast.FuncLit:
		return n.Type
	}
	return nil
}
