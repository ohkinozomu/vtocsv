package vtocsv

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/bmatcuk/doublestar/v3"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func updateRecords(ssapkg *ssa.Package, outputCSV [][]string) [][]string {
	members := ssapkg.Members
	var functions []string
	for _, v := range members {
		_, ok := v.(*ssa.Function)
		if ok {
			functions = append(functions, v.(*ssa.Function).Name())
		}
	}

	for _, f := range functions {
		fun := ssapkg.Func(f)
		for _, l := range fun.Locals {
			outputCSV = append(outputCSV, []string{l.Comment, l.Type().String()})
		}
	}
	return outputCSV
}

func findGoFiles(path string) ([]string, error) {
	goFles, err := doublestar.Glob(path + "/**/*.go")
	if err != nil {
		return nil, err
	}
	return goFles, nil
}

func Output(path string) ([][]string, error) {
	var outputCSV [][]string

	files, err := findGoFiles(path)
	if err != nil {
		return outputCSV, err
	}

	outputCSV = append(outputCSV, []string{"comment", "type"})

	for _, f := range files {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, f, nil, parser.AllErrors)
		if err != nil {
			return outputCSV, err
		}
		files := []*ast.File{f}

		pkg := types.NewPackage(f.Name.Name, "")
		ssapkg, _, err := ssautil.BuildPackage(
			&types.Config{Importer: importer.Default()},
			fset, pkg, files,
			ssa.GlobalDebug,
		)
		if err != nil {
			return outputCSV, err
		}
		outputCSV = updateRecords(ssapkg, outputCSV)
	}

	return outputCSV, nil
}
