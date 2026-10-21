# Шаг 2: создать задачу

## Что добавили

К `GET /api/tasks` добавили `POST /api/tasks`. Сервер читает JSON из тела запроса, проверяет `text`, сохраняет новую задачу в памяти и возвращает статус `201 Created`.

Посмотрите: `createTask` и `writeJSON` в `internal/httpapi/handler.go`.

## Запуск

```bash
go run ./steps/02_create_task/cmd/server
```

## Проверка

Создайте задачу:

```bash
curl -i -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"text":"Показать POST-запрос"}'
```

Затем получите весь список:

```bash
curl -i http://localhost:8080/api/tasks
```

Проверьте ошибку валидации:

```bash
curl -i -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"text":""}'
```

Последний запрос должен вернуть `400 Bad Request`.
