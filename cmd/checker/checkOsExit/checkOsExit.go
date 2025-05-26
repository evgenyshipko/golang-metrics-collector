package checkOsExit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var NoDirectOsExitAnalyzer = &analysis.Analyzer{
	Name:     "noosxit",
	Doc:      "forbid direct os.Exit calls in main function of main package",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runNoOsExit,
}

func runNoOsExit(pass *analysis.Pass) (interface{}, error) {
	// Проверяем, что это пакет main
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)

		// Нас интересует только функция main
		if fn.Name.Name != "main" {
			return
		}

		// Ищем вызовы os.Exit в теле функции
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			ident, ok := sel.X.(*ast.Ident)
			if !ok {
				return true
			}

			if ident.Name == "os" && sel.Sel.Name == "Exit" {
				pass.Reportf(call.Pos(), "direct call to os.Exit in main function of main package is forbidden")
			}

			return true
		})
	})

	return nil, nil
}
