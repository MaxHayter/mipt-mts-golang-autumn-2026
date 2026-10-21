# Шаг 5: загрузить файл

## Что добавили

Добавили `POST /upload` с форматом `multipart/form-data`. Хендлер принимает поле `file`, ограничивает форму 5 MiB и копирует файл в папку `uploads/` этого шага.

Посмотрите: константу `maxUploadSize` и метод `upload` в `internal/httpapi/handler.go`.

## Запуск

```bash
go run ./steps/05_upload/cmd/server
```

## Проверка

Загрузите файл из корня семинара:

```bash
curl -i -F 'file=@README.md' http://localhost:8080/upload
```

Ожидаемый ответ: `200 OK` и `uploaded: README.md`. Файл появится в `steps/05_upload/uploads/README.md`.

Проверьте отсутствие обязательного поля:

```bash
curl -i -X POST http://localhost:8080/upload
```

Ожидаемый статус — `400 Bad Request`.
