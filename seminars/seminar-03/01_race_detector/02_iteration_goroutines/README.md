# Шаг 2: добавляем goroutine

## Запуск

```bash
go run ./01_race_detector/02_iteration_goroutines
go run -race ./01_race_detector/02_iteration_goroutines
```

Без `-race` число может быть меньше 1000. С `-race` ожидается `WARNING: DATA RACE` и ненулевой код завершения — это успешная демонстрация дефекта.

## Что сказать

- `go` не ждёт функцию: после запуска всех задач `main` продолжает работу.
- `time.Sleep` — это догадка о времени, а не гарантия завершения. На другом компьютере она может дать иной результат.
- Даже если увеличить Sleep и увидеть 1000, проблема остаётся: несколько goroutine читают и пишут `count` без синхронизации.
- Race detector показывает две конфликтующие операции и строки, где были созданы goroutine.

## Время

```bash
/usr/bin/time -p go run ./01_race_detector/02_iteration_goroutines
/usr/bin/time -p go run -race ./01_race_detector/02_iteration_goroutines
```

Не делайте вывод «быстрее значит правильнее»: результат и время здесь недетерминированы.
