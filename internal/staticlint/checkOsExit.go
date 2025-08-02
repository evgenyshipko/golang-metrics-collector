package staticlint

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NoDirectOsExitAnalyzer запрещает использовать прямой вызов os.Exit в функции main пакета main.
var NoDirectOsExitAnalyzer = &analysis.Analyzer{
	Name:     "noosxit",
	Doc:      "forbid direct os.Exit calls in main function of main package",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runNoOsExit,
}

// runNoOsExit коллбэк с логикой анализатора NoDirectOsExitAnalyzer
func runNoOsExit(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)

		if fn.Name.Name != "main" {
			return
		}

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
