// Что видит генератор: go/parser превращает исходник в дерево go/ast.
// Тип поля в дереве — только синтаксис, имя как написано в коде.
//
//	go run ./code/08-codegen/ast
package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

//go:embed user.go
var src string

func main() {
	// Без ParseComments комментарии в дерево не попадут и метку //gen:json не найти.
	file, err := parser.ParseFile(token.NewFileSet(), "user.go", src, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		ts := gd.Specs[0].(*ast.TypeSpec)
		st, ok := ts.Type.(*ast.StructType)
		if !ok || gd.Doc == nil {
			continue
		}
		fmt.Printf("type %s, комментарий над ним: %s\n", ts.Name.Name, gd.Doc.List[0].Text)
		for _, f := range st.Fields.List {
			fmt.Printf("  %-8s %-18T %-10s тег %s\n",
				f.Names[0].Name, f.Type, types.ExprString(f.Type), f.Tag.Value)
		}
	}
}
