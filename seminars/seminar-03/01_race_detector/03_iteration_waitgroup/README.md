# Шаг 3: WaitGroup ждёт, но не защищает данные

## Запуск

```bash
go run ./01_race_detector/03_iteration_waitgroup
go run -race ./01_race_detector/03_iteration_waitgroup
```

Программа теперь всегда дожидается workers, но команда с `-race` всё ещё должна вывести `WARNING: DATA RACE`.

## Что сказать

- `Add(1)` увеличивает счётчик до запуска goroutine; `Done()` уменьшает его при выходе; `Wait()` ждёт ноль.
- Теперь мы устранили одну проблему — раннее завершение `main`.
- Но `WaitGroup` не делает `count++` атомарным и не запрещает двум goroutine менять переменную одновременно.
- Поэтому ожидание и синхронизация доступа к данным — разные задачи.

## Время

```bash
/usr/bin/time -p go run ./01_race_detector/03_iteration_waitgroup
/usr/bin/time -p go run -race ./01_race_detector/03_iteration_waitgroup
```

`WaitGroup` добавлен ради корректного жизненного цикла, а не ради ускорения.
