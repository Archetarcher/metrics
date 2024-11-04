// Package osexitchecker defines Analyzer that check os.Exit call in main package
package osexitchecker

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit call",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if s, ok := call.Fun.(*ast.SelectorExpr); ok {
				if s.Sel.Name == "Exit" {
					pass.Reportf(x.Pos(), "os exit call detected ")
				}
			}
		}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {

			if c, ok := node.(*ast.ExprStmt); ok {
				if strings.Contains(pass.Fset.Position(file.Package).Filename, "metrics/cmd/server/checker.go") ||
					strings.Contains(pass.Fset.Position(file.Package).Filename, "metrics/cmd/agent/checker.go") {
					expr(c)
				}
			}
			return true
		})
	}
	return nil, nil
}
