# Шаг 1: получить задачи

## Что добавили

Первый HTTP-маршрут `GET /api/tasks`. Сервер берёт задачи из памяти и кодирует их в JSON.

Посмотрите: `internal/httpapi/handler.go` — здесь появляются `http.ServeMux`, `ResponseWriter` и `json.NewEncoder`.

## Запуск

Из корня `http-seminar`:

```bash
go run ./steps/01_get_tasks/cmd/server
```

## Проверка

Во втором терминале:

```bash
curl -i http://localhost:8080/api/tasks
```

Ожидаемый ответ: статус `200 OK` и JSON-массив с задачей `"Подготовить семинар"`.
