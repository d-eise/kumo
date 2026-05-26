package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	name = "kumo"
	doc  = "kumo is a static analysis tool that detects issues in Go code"
)

// Analyzer is the main analyzer for kumo.
var Analyzer = &analysis.Analyzer{
	Name:     name,
	Doc:      doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.GoStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch node := n.(type) {
		case *ast.GoStmt:
			checkGoStmt(pass, node)
		case *ast.FuncDecl:
			checkFuncDecl(pass, node)
		}
	})

	return nil, nil
}

func checkGoStmt(pass *analysis.Pass, node *ast.GoStmt) {
	// Report bare goroutine launches without any synchronization context.
	if _, ok := node.Call.Fun.(*ast.FuncLit); ok {
		pass.Reportf(node.Pos(), "goroutine launched with anonymous function literal; consider naming it for better stack traces")
	}
}

func checkFuncDecl(pass *analysis.Pass, node *ast.FuncDecl) {
	if node.Body == nil {
		return
	}
	// Check for functions that are excessively long.
	const maxLines = 100
	startLine := pass.Fset.Position(node.Body.Lbrace).Line
	endLine := pass.Fset.Position(node.Body.Rbrace).Line
	if endLine-startLine > maxLines {
		pass.Reportf(node.Pos(), "function %q is too long (%d lines); consider refactoring", node.Name.Name, endLine-startLine)
	}
}
