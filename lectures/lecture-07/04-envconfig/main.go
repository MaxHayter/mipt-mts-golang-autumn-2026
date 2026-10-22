// Загрузка конфига из переменных окружения по тегам — мини-версия envconfig.
// Одна функция на любую структуру: типы полей становятся известны только в рантайме.
//
//	go run ./code/04-envconfig
//	PORT=9090 DEBUG=true TIMEOUT=2s go run ./code/04-envconfig
//	PORT=abc go run ./code/04-envconfig
package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Config struct {
	Host    string        `env:"HOST" default:"localhost"`
	Port    int           `env:"PORT" default:"8080"`
	Debug   bool          `env:"DEBUG"`
	Timeout time.Duration `env:"TIMEOUT" default:"500ms"`
	Secret  string        // без тега — не трогаем
}

func Load(dst any) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("envconfig: нужен указатель на структуру, получен %T", dst)
	}
	v = v.Elem()

	for f, fv := range v.Fields() {
		key, ok := f.Tag.Lookup("env")
		if !ok {
			continue
		}
		raw, ok := os.LookupEnv(key)
		if !ok {
			raw, ok = f.Tag.Lookup("default")
		}
		if !ok {
			continue
		}
		if err := set(fv, raw); err != nil {
			return fmt.Errorf("envconfig: %s (%s=%q): %w", f.Name, key, raw, err)
		}
	}
	return nil
}

func set(fv reflect.Value, raw string) error {
	// Сначала конкретный тип, потом Kind: у time.Duration Kind — Int64,
	// и без этой проверки "2s" ушло бы в ParseInt.
	if fv.Type() == reflect.TypeFor[time.Duration]() {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return err
		}
		fv.SetInt(int64(d))
		return nil
	}

	switch fv.Kind() {
	case reflect.String:
		fv.SetString(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(raw, 10, fv.Type().Bits())
		if err != nil {
			return err
		}
		fv.SetInt(n)
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	default:
		return fmt.Errorf("тип %v не поддерживается", fv.Type())
	}
	return nil
}

func main() {
	var cfg Config
	if err := Load(&cfg); err != nil {
		fmt.Println("ошибка:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", cfg)

	fmt.Println(Load(cfg)) // забыли &
}
