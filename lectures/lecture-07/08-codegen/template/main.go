// Вторая половина генератора: данные о полях → text/template → go/format.
// Поля здесь заданы руками; в gen их достаёт из дерева go/ast.
//
//	go run ./code/08-codegen/template
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"text/template"
)

type field struct {
	Name   string // имя поля в Go
	Prefix string // готовый кусок JSON перед значением
	Append string // функция jsonenc для значения
}

// Отступы в шаблоне не важны: format.Source всё выровняет, как gofmt.
const tplText = `func (v {{.Name}}) MarshalJSON() ([]byte, error) {
buf := make([]byte, 0, 128)
{{- range .Fields}}
buf = append(buf, {{printf "%q" .Prefix}}...)
buf = jsonenc.{{.Append}}(buf, v.{{.Name}})
{{- end}}
return append(buf, '}'), nil
}`

func main() {
	data := map[string]any{
		"Name": "User",
		"Fields": []field{
			{Name: "Login", Prefix: `{"login":`, Append: "AppendString"},
			{Name: "Admin", Prefix: `,"admin":`, Append: "AppendBool"},
		},
	}
	tpl := template.Must(template.New("marshal").Parse(tplText))
	var src bytes.Buffer
	if err := tpl.Execute(&src, data); err != nil {
		log.Fatal(err)
	}
	// format.Source не только форматирует, но и парсит: ошибка в шаблоне всплывёт здесь.
	out, err := format.Source(src.Bytes())
	if err != nil {
		log.Fatalf("сгенерирован некорректный код: %v\n%s", err, src.Bytes())
	}
	fmt.Println(string(out))
}
