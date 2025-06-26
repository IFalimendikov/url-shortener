package main

import (
	"go/ast"
	"os"
	"path/filepath"
	_ "url-shortener/docs"

	"encoding/json"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "addlint",
	Doc:  "reports integer additions",
	Run:  run,
}

const Config = `config.json`

type ConfigData struct {
	Staticcheck []string
}

func main() {

	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(workDir, Config))
	if err != nil {
		panic(err)
	}

	var cfgLint ConfigData
	if err := json.Unmarshal([]byte(data), &cfgLint); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		ExitCheckAnalyzer,
	}
	checks := make(map[string]bool)
	for _, v := range cfgLint.Staticcheck {
		checks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
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
