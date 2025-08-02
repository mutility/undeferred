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
	referredShadows := map[types.Object]struct{}{}
	resultNames := map[string][]types.Object{}
	resultObjs := map[types.Object]struct{}{}
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspect.WithStack(
		[]ast.Node{new(ast.FuncDecl), new(ast.FuncLit), new(ast.DeferStmt), new(ast.Ident)},
		func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
			// update known result names/objs
			if f := funcOf(n); f != nil && f.Results != nil {
				for _, r := range f.Results.List {
					for _, id := range r.Names {
						o := pass.TypesInfo.Defs[id]
						named := resultNames[o.Name()]
						if push {
							resultNames[o.Name()] = append(named, o)
							resultObjs[o] = struct{}{}
						} else {
							resultNames[o.Name()] = named[:len(named)-1]
							delete(resultObjs, o)
						}
					}
				}
			}
			switch n := n.(type) {
			// check for use of named result directly in defer's call
			case *ast.DeferStmt:
				if !push {
					return true
				}
				for _, a := range n.Call.Args {
					switch a := a.(type) {
					case *ast.Ident:
						if _, ok := resultObjs[pass.TypesInfo.Uses[a]]; ok {
							pass.Reportf(a.Pos(), "defer captures current value of named result '%s'", a.Name)
						}
					}
				}
			// check for use of shadows matching a named result in a deferred func
			case *ast.Ident:
				deferStmt := closest[*ast.DeferStmt](stack)
				if deferStmt < 0 ||
					(deferStmt > closest[*ast.FuncDecl](stack) && deferStmt > closest[*ast.FuncLit](stack)) {
					return // don't consider ids outside a defer, or arguments to a defer
				}
				if o := pass.TypesInfo.Uses[n]; o != nil && push {
					if _, ok := resultObjs[o]; ok {
						return true // refers to the named return value; that's fine
					}
					for _, n := range stack[deferStmt+1:] {
						if scope := pass.TypesInfo.Scopes[n]; scope != nil && scope.Lookup(o.Name()) == o {
							return true // ignore references to ids defined within the defer func
						}
					}
					if scope := pass.TypesInfo.Scopes[funcOf(stack[deferStmt].(*ast.DeferStmt).Call.Fun)]; scope != nil {
						if scope.Lookup(o.Name()) == o {
							return true
						}
					}
					if objs := resultNames[o.Name()]; len(objs) > 0 {
						pass.Reportf(n.Pos(), "defer references shadow of named result '%s'", o.Name())
						if _, ok := referredShadows[o]; !ok {
							pass.Reportf(o.Pos(), "shadows named result '%s' referenced in later defer", o.Name())
							referredShadows[o] = struct{}{}
						}
					}
				}
				return true
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

func closest[T ast.Node](stack []ast.Node) int {
	for i := len(stack) - 1; i > 0; i-- {
		if _, ok := stack[i].(T); ok {
			return i
		}
	}
	return -1
}
