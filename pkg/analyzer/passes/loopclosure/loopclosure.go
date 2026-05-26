package loopclosure

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "loopclosure"
	doc  = "loopclosure checks for goroutines launched inside loops that capture loop variables"
)

// Analyzer is the loopclosure analyzer.
var Analyzer = &analysis.Analyzer{
	Name:     name,
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.RangeStmt)(nil),
		(*ast.ForStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch loop := n.(type) {
		case *ast.RangeStmt:
			checkLoopBody(pass, loop.Body, collectRangeVars(loop))
		case *ast.ForStmt:
			checkLoopBody(pass, loop.Body, collectForVars(loop))
		}
	})

	return nil, nil
}

func collectRangeVars(loop *ast.RangeStmt) map[string]bool {
	vars := make(map[string]bool)
	if ident, ok := loop.Key.(*ast.Ident); ok && ident.Name != "_" {
		vars[ident.Name] = true
	}
	if loop.Value != nil {
		if ident, ok := loop.Value.(*ast.Ident); ok && ident.Name != "_" {
			vars[ident.Name] = true
		}
	}
	return vars
}

func collectForVars(loop *ast.ForStmt) map[string]bool {
	vars := make(map[string]bool)
	if loop.Init != nil {
		if assign, ok := loop.Init.(*ast.AssignStmt); ok {
			for _, lhs := range assign.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok {
					vars[ident.Name] = true
				}
			}
		}
	}
	return vars
}

func checkLoopBody(pass *analysis.Pass, body *ast.BlockStmt, loopVars map[string]bool) {
	if body == nil || len(loopVars) == 0 {
		return
	}
	for _, stmt := range body.List {
		goStmt, ok := stmt.(*ast.GoStmt)
		if !ok {
			continue
		}
		ast.Inspect(goStmt.Call, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if loopVars[ident.Name] {
				pass.Reportf(goStmt.Pos(), "goroutine captures loop variable %q", ident.Name)
			}
			return true
		})
	}
}
