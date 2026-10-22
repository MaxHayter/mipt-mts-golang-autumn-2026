// Package refljson — JSON-маршалинг через reflect, «как encoding/json, но на одном экране».
// Поддерживает структуры с полями string, int*, bool и теги json:"name", json:"-",
// json:",omitempty". Два варианта: Marshal разбирает тип на каждом вызове,
// (*Encoder).Marshal — один раз на тип, дальше идёт по готовому плану.
package refljson

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"lecture-07/jsonenc"
)

func Marshal(v any) ([]byte, error) {
	rv, err := structValue(v)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 0, 128) // как в сгенерированном коде — сравниваем только обход полей
	buf = append(buf, '{')
	first := true
	// На каждом вызове заново: обход полей, разбор тегов, switch по Kind.
	t := rv.Type()
	for i := range t.NumField() {
		f, fv := t.Field(i), rv.Field(i)
		if !f.IsExported() {
			continue
		}
		name, omitEmpty := f.Name, false
		if tag, ok := f.Tag.Lookup("json"); ok {
			var opts string
			name, opts, _ = strings.Cut(tag, ",")
			if name == "-" {
				continue
			}
			if name == "" {
				name = f.Name
			}
			omitEmpty = opts == "omitempty"
		}
		if omitEmpty && fv.IsZero() {
			continue
		}

		if !first {
			buf = append(buf, ',')
		}
		first = false
		buf = jsonenc.AppendString(buf, name)
		buf = append(buf, ':')

		if buf, err = appendValue(buf, fv); err != nil {
			return nil, fmt.Errorf("refljson: поле %s: %w", f.Name, err)
		}
	}
	return append(buf, '}'), nil
}

// fieldPlan — всё, что Marshal выясняет про поле на каждом вызове, посчитанное заранее.
type fieldPlan struct {
	index     int
	key       string // готовый кусок "name":
	kind      reflect.Kind
	omitEmpty bool
}

// Encoder кеширует план полей по типу — так делает encoding/json.
// Нулевое значение готово к работе; копировать после первого использования нельзя.
type Encoder struct {
	cache sync.Map // reflect.Type → []fieldPlan
}

func (e *Encoder) Marshal(v any) ([]byte, error) {
	rv, err := structValue(v)
	if err != nil {
		return nil, err
	}
	plan, err := e.plan(rv.Type())
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 0, 128)
	buf = append(buf, '{')
	first := true
	// Ни тегов, ни IsExported: только значения по готовым индексам.
	for _, p := range plan {
		fv := rv.Field(p.index)
		if p.omitEmpty && fv.IsZero() {
			continue
		}
		if !first {
			buf = append(buf, ',')
		}
		first = false
		buf = append(buf, p.key...)
		buf, _ = appendValue(buf, fv) // Kind проверен при построении плана
	}
	return append(buf, '}'), nil
}

func (e *Encoder) plan(t reflect.Type) ([]fieldPlan, error) {
	if p, ok := e.cache.Load(t); ok {
		return p.([]fieldPlan), nil
	}
	var plan []fieldPlan
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		p := fieldPlan{index: i, kind: f.Type.Kind()}
		name := f.Name
		if tag, ok := f.Tag.Lookup("json"); ok {
			var opts string
			name, opts, _ = strings.Cut(tag, ",")
			if name == "-" {
				continue
			}
			if name == "" {
				name = f.Name
			}
			p.omitEmpty = opts == "omitempty"
		}
		if !supported(p.kind) {
			return nil, fmt.Errorf("refljson: поле %s: тип %v не поддерживается", f.Name, f.Type)
		}
		p.key = string(append(jsonenc.AppendString(nil, name), ':'))
		plan = append(plan, p)
	}
	// Две горутины могли построить план одновременно — в кеше останется первый.
	p, _ := e.cache.LoadOrStore(t, plan)
	return p.([]fieldPlan), nil
}

func structValue(v any) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("refljson: нужна структура, получен %T", v)
	}
	return rv, nil
}

func supported(k reflect.Kind) bool {
	switch k {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}

func appendValue(buf []byte, fv reflect.Value) ([]byte, error) {
	switch fv.Kind() {
	case reflect.String:
		return jsonenc.AppendString(buf, fv.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return jsonenc.AppendInt(buf, fv.Int()), nil
	case reflect.Bool:
		return jsonenc.AppendBool(buf, fv.Bool()), nil
	}
	return buf, fmt.Errorf("тип %v не поддерживается", fv.Type())
}
