# Семинар: как находить проблемы конкурентности в Go

Этот практикум продолжает лекцию о goroutine, GMP, примитивах синхронизации, race condition и каналах. Он **не вводит** fan-in/fan-out, worker pool или `errgroup`: для студентов, ещё не создававших каналы, это будет следующий шаг.

В каждой теме есть две стадии:

- `01_initial` — минимальная программа с дефектом и наблюдаемым симптомом;
- `02_iteration` — исправление, которое устраняет конкретную причину.

## Маршрут семинара на 80 минут

1. **0–10 мин.** Повторите: goroutine конкурентны даже на одном ядре; планировщик может переключить выполнение между строками `count++`.
2. **10–35 мин.** `01_race_detector`: пройдите три версии. Сначала «голые» goroutine и ненадёжный `Sleep`; затем `WaitGroup`, который ждёт, но не устраняет гонку; в конце `Mutex` как синхронизация доступа к `count`.
3. **35–55 мин.** `02_goroutine_leak`: сравните число goroutine до и после запуска. Причина — нет условия выхода. Канал здесь показан только как «сигнал остановки», без паттернов.
4. **55–72 мин.** `03_memory_leak`: посмотрите `HeapAlloc`, затем heap profile. Причина — маленький subslice удерживает большой backing array.
5. **72–80 мин.** Итог: сначала воспроизвести, затем измерить, затем исправить и повторно проверить.

## Команды показа

```bash
cd concurrency-seminar

go run ./01_race_detector/01_initial
go run -race ./01_race_detector/02_iteration_goroutines
go run -race ./01_race_detector/03_iteration_waitgroup
go run -race ./01_race_detector/04_iteration_mutex

go run ./02_goroutine_leak/01_initial
go run ./02_goroutine_leak/02_iteration

go run ./03_memory_leak/01_initial
go tool pprof -top memory-initial.pprof
go run ./03_memory_leak/02_iteration
go tool pprof -top memory-fixed.pprof
```

Профили `memory-*.pprof` создаются командой примера в текущей рабочей папке. Их можно удалить после занятия.
