/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"strconv"
	"strings"
)

// ParseTopLevelConstants parses the values of the top-level
// constants in the specified directory (excludes *_test.go).
// Returns error or two slices: one with names and values
// where names[I] is a name for values[I] value.
func ParseTopLevelConstants(path, prefix string) ([]string, []string, error) {
	fset := token.NewFileSet()
	// ParseDir returns a map of package name to package ASTs. An AST is a representation of the source code
	// that can be traversed to extract information. The map is keyed by the package name.
	pkgs, err := parser.ParseDir(fset, path, func(info fs.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		return nil, nil, err
	}

	var constNames []string
	var constValues []string

	// Iterate all of the top-level declarations in each package's files,
	// looking for constants. When we find one, add its
	// name into constNames slice and value to the constValues.
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			for _, d := range f.Decls {
				genDecl, ok := d.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, s := range genDecl.Specs {
					valueSpec, ok := s.(*ast.ValueSpec)
					if !ok || len(valueSpec.Names) != 1 || valueSpec.Names[0].Obj.Kind != ast.Con || !strings.HasPrefix(valueSpec.Names[0].String(), prefix) {
						continue
					}
					for _, val := range valueSpec.Values {
						lit, ok := val.(*ast.BasicLit)
						if !ok || lit.Kind != token.STRING {
							continue
						}
						v, err := strconv.Unquote(lit.Value)
						if err != nil {
							return nil, nil, fmt.Errorf("unquote literal string %s: %v", lit.Value, err)
						}
						constNames = append(constNames, valueSpec.Names[0].String())
						constValues = append(constValues, v)
					}
				}
			}
		}
	}
	return constNames, constValues, nil
}
