# Шаг 4: HTML-шаблон

## Что добавили

Добавили страницу `GET /`. Хендлер передаёт список задач в `html/template`, а шаблон из `web/templates/index.html` строит HTML-список.

Посмотрите: поле `page`, метод `index` и маршрут `GET /` в `internal/httpapi/handler.go`.

## Запуск

```bash
go run ./steps/04_html_template/cmd/server
```

## Проверка

Откройте в браузере `http://localhost:8080/` или получите HTML через `curl`:

```bash
curl -i http://localhost:8080/
```

Создайте задачу через API, затем обновите страницу:

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"text":"Увидеть задачу в браузере"}'
curl http://localhost:8080/
```
