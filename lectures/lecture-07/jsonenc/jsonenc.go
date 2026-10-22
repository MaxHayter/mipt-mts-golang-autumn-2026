// Package jsonenc — общий помощник для 06-refljson и сгенерированного кода 08-codegen:
// экранирование строки по правилам JSON. Вынесен, чтобы бенчмарк сравнивал
// только способ обхода полей, а не разные реализации экранирования.
package jsonenc

import (
	"strconv"
	"unicode/utf8"
)

// AppendString дописывает s в buf как JSON-строку в кавычках.
// Упрощение: без HTML-экранирования <>& и U+2028/2029, которое делает encoding/json.
func AppendString(buf []byte, s string) []byte {
	buf = append(buf, '"')
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case c == '"' || c == '\\':
			buf = append(buf, '\\', c)
		case c == '\n':
			buf = append(buf, '\\', 'n')
		case c == '\r':
			buf = append(buf, '\\', 'r')
		case c == '\t':
			buf = append(buf, '\\', 't')
		case c < 0x20:
			buf = append(buf, `\u00`...)
			buf = append(buf, "0123456789abcdef"[c>>4], "0123456789abcdef"[c&0xf])
		case c < utf8.RuneSelf:
			buf = append(buf, c)
		default:
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				buf = append(buf, `�`...)
			} else {
				buf = append(buf, s[i:i+size]...)
			}
			i += size
			continue
		}
		i++
	}
	return append(buf, '"')
}

// AppendInt и AppendBool — чтобы сгенерированный код импортировал один пакет.
func AppendInt(buf []byte, n int64) []byte { return strconv.AppendInt(buf, n, 10) }

func AppendBool(buf []byte, b bool) []byte { return strconv.AppendBool(buf, b) }
