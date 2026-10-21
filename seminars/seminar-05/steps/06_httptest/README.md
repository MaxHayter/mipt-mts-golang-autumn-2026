# Шаг 6: протестировать хендлер

## Что добавили

Добавили `handler_test.go`. В нём `httptest.NewRequest` создаёт запрос в памяти, `httptest.NewRecorder` записывает ответ, а `Routes().ServeHTTP` вызывает настоящий хендлер без порта и запущенного сервера.

Хендлер не меняется: тестируем поведение уже готового `POST /api/tasks`.

## Проверка

Из корня семинара выполните:

```bash
go test -v ./steps/06_httptest/internal/httpapi
```

Тест проверяет, что корректный запрос возвращает `201 Created` и JSON созданной задачи.

Чтобы сначала увидеть API вручную, можно запустить сервер:

```bash
go run ./steps/06_httptest/cmd/server
curl -i -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"text":"Проверить API вручную"}'
```
