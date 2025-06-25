package staticlint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "addlint",
	Doc:  "reports integer additions",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
    for _, file := range pass.Files {
        if pass.Pkg.Name() != "main" {
            continue
        }
        for _, decl := range file.Decls {
            fn, ok := decl.(*ast.FuncDecl)
            if !ok || fn.Name.Name != "main" {
                continue
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
                pkgIdent, ok := sel.X.(*ast.Ident)
                if !ok {
                    return true
                }
                if pkgIdent.Name == "os" && sel.Sel.Name == "Exit" {
                    obj := pass.TypesInfo.Uses[pkgIdent]
                    if obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "os" {
                        pass.Reportf(call.Lparen, "direct call to os.Exit in main.main is forbidden")
                    }
                }
                return true
            })
        }
    }
    return nil, nil
}