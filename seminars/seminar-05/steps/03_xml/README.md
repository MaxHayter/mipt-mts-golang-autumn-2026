# Шаг 3: XML

## Что добавили

Добавили `GET /api/tasks.xml`. Данные те же, что в JSON API, но сервер сериализует их через `xml.NewEncoder` и указывает `Content-Type: application/xml`.

Посмотрите: `taskList` и `tasksXML` в `internal/httpapi/handler.go`.

## Запуск

```bash
go run ./steps/03_xml/cmd/server
```

## Проверка

Сравните два представления одного набора данных:

```bash
curl -i http://localhost:8080/api/tasks
curl -i http://localhost:8080/api/tasks.xml
```

Создайте задачу и снова запросите XML:

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H 'Content-Type: application/json' \
  -d '{"text":"Сравнить JSON и XML"}'
curl http://localhost:8080/api/tasks.xml
```
