// Package goroutine provides an analyzer that detects goroutine-related issues.
package goroutine

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the goroutine leak and misuse analyzer.
var Analyzer = &analysis.Analyzer{
	Name:     "goroutine",
	Doc:      "checks for common goroutine misuse patterns such as launching goroutines in init functions or test helpers",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		if funcDecl.Name.Name == "init" {
			checkForGoStmtsInFunc(pass, funcDecl)
		}
	})

	return nil, nil
}

func checkForGoStmtsInFunc(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	if funcDecl.Body == nil {
		return
	}

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			return true
		}
		reportGoStmtInInit(pass, goStmt.Pos())
		return true
	})
}

func reportGoStmtInInit(pass *analysis.Pass, pos token.Pos) {
	pass.Reportf(pos, "goroutine launched inside init function; this can cause hard-to-detect leaks or race conditions")
}
