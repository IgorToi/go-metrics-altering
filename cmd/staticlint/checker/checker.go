// Package checker defines analyzer to check whether os.Exit is used.
package checker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "checker helps determine if os.exit exists",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if x, ok := node.(*ast.CallExpr); ok {
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if s.Sel.Name == "Exit" {
						pass.Reportf(s.Pos(), "using os.Exit is prohibbited")
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
