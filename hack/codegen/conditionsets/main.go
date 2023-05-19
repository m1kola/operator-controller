package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"sort"

	"golang.org/x/tools/go/ast/astutil"

	internalast "github.com/operator-framework/operator-controller/internal/ast"
)

func main() {
	apiPathFlag := flag.String("apiPath", "", "Path to the API package on the filesystem")
	apiImportAliasFlag := flag.String("apiImportAlias", "", "API import alias. Will be used for adding import into the target file")
	apiImportFlag := flag.String("apiImport", "", "API import. Will be used for adding import into the target file")
	targetFilePathFlag := flag.String("targetFilePath", "", "Path to the target file")
	targetPackageName := flag.String("targetPackageName", "", "Package name for the target file")
	fromPrefixToSliceFlag := make(prefixMapFlag)
	flag.Var(&fromPrefixToSliceFlag, "fromPrefixToSlice", "Map of top level constants from from API package to variable name in the target file")
	flag.Parse()

	err := run(*apiPathFlag, *apiImportAliasFlag, *apiImportFlag, *targetPackageName, *targetFilePathFlag, fromPrefixToSliceFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Code generation completed successfully.")
}

func run(apiPath, apiImportAlias, apiImport, targetPackageName, targetFilePath string, fromPrefixToSlice map[string]string) error {
	fset := token.NewFileSet()

	// Generate modifications to the code
	file, err := doCodeGen(fset, apiPath, apiImportAlias, apiImport, targetPackageName, fromPrefixToSlice)
	if err != nil {
		return fmt.Errorf("failed code generation: %w", err)
	}

	outputFile, err := os.Create(targetFilePath)
	if err != nil {
		return fmt.Errorf("failed to create or update output file: %w", err)
	}
	defer outputFile.Close()

	// Write changes to the file
	err = format.Node(outputFile, fset, file)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func doCodeGen(fset *token.FileSet, apiPath, apiImportAlias, apiImport, targetPackageName string, fromPrefixToSlice map[string]string) (*ast.File, error) {
	// Parse constants
	targetVarNameToAPIConstNames := map[string][]string{}
	for apiConstPrefix, targetVarName := range fromPrefixToSlice {
		apiConstNames, _, err := internalast.ParseTopLevelConstants(apiPath, apiConstPrefix)
		if err != nil {
			return nil, fmt.Errorf("failed parse: %w", err)
		}

		targetVarNameToAPIConstNames[targetVarName] = apiConstNames
	}

	const infinity = 1 << 30
	nextPos := token.Pos(fset.AddFile("generating.go", -1, infinity).Base())

	comments := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Slash: nextPos,
				Text:  "// Code generated. DO NOT EDIT.",
			},
		},
	}

	// Add an empty line after the comment
	offset := fset.File(nextPos).Offset(nextPos)
	fset.File(nextPos).AddLine(offset + 1)
	fset.File(nextPos).AddLine(offset + 2)

	nextPos += 2

	// Create the AST file
	newFile := &ast.File{
		Package: nextPos,
		Name: &ast.Ident{
			NamePos: nextPos,
			Name:    targetPackageName,
		},
		Doc:      comments,
		Comments: []*ast.CommentGroup{comments},
	}

	// Ensure that the desired import is present
	astutil.AddNamedImport(fset, newFile, apiImportAlias, apiImport)

	generateGoFile(newFile, apiImportAlias, targetVarNameToAPIConstNames)

	return newFile, nil
}

func generateGoFile(file *ast.File, apiImportAlias string, targetVarNameToAPIConstNames map[string][]string) {
	// Ensure deterministic order
	keys := make([]string, 0, len(targetVarNameToAPIConstNames))
	for key := range targetVarNameToAPIConstNames {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Create the package level variables
	for _, varName := range keys {
		// Create the variable declaration
		varDecl := &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{
						{Name: varName},
					},
					Values: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.ArrayType{
								Elt: &ast.Ident{Name: "string"},
							},
							Elts: createImportedIdentifiers(apiImportAlias, targetVarNameToAPIConstNames[varName]),
						},
					},
				},
			},
		}

		// Add the variable declaration to the file
		file.Decls = append(file.Decls, varDecl)
	}
}

func createImportedIdentifiers(importAlias string, values []string) []ast.Expr {
	var exprs []ast.Expr
	for _, value := range values {
		exprs = append(exprs, &ast.SelectorExpr{
			X:   &ast.Ident{Name: importAlias},
			Sel: &ast.Ident{Name: value},
		})
	}
	return exprs
}
