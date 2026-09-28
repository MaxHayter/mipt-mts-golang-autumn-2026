# Набор терминальных команд для семинара 02

## Вывод результатов escape analyse

```bash
go build -gcflags="-m" main.go
```

Отключить inline функций:
```bash
go build -gcflags="-l -m" main.go
```

Развернутый анализ компилятора:
```bash
go build -gcflags="-m -m" main.go
```

## Влияние escape analysis на производительность

С включенным механизмом inline
```bash
go test -bench=. -benchmem
```

Отключить inline механизм
```bash
go test -bench=. -gcflags="-l" -benchmem
```

## Работа с GC

Запуск трассировки для отслеживания работы GC в реальном времени
```bash
GODEBUG=gctrace=1 go run . 
```

Управление сборщиком мусора для одного зщапуска
```bash
GOGC=200 GODEBUG=gctrace=1 go run .
```
