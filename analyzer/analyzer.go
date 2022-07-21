// Package staticlint This package is intended for static code analyzation.
package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

const (
	PackageName  = "main" // The name of the package being checked by the analyzer.
	FunctionName = "main" // The name of the function being checked by the analyzer.
	Expression   = "Exit" // The name of the expression is looking for by the analyzer.
)

var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for using os.Exit expression in the main file",
	Run:  run,
}

// run This function is run analyze static code.
func run(pass *analysis.Pass) (interface{}, error) {
	checkExitFunc := func(x *ast.CallExpr, packageName string, funName string) {
		sel, ok := x.Fun.(*ast.SelectorExpr)
		if ok {
			if packageName == PackageName && funName == FunctionName && sel.Sel.String() == Expression {
				pass.Reportf(x.Pos(), "you are using os.Exit() expression in the main file!")
			}
		}
	}
	lastFunName := ""
	for _, file := range pass.Files {
		// iterate traverse through all AST nodes.
		ast.Inspect(file, func(node ast.Node) bool {
			lastFunName = getLastFunName(node)
			switch x := node.(type) {
			case *ast.CallExpr:
				checkExitFunc(x, pass.Pkg.Name(), lastFunName)
			}
			return true
		})
	}
	return nil, nil
}

// getLastFunName returns last function name was used.
func getLastFunName(node ast.Node) string {
	lastFunName := ""
	switch x := node.(type) {
	case *ast.FuncDecl:
		lastFunName = x.Name.Name
	}

	return lastFunName
}
