# Семинар: HTTP-сервис на стандартной библиотеке Go

## ТЗ

Нужно сделать небольшой сервис «Доска задач для учебной группы». Задачи храним в памяти.

- `GET /` показывает HTML-страницу со списком задач;
- `GET /api/tasks` отдаёт задачи в JSON;
- `POST /api/tasks` принимает JSON вида `{"text":"Подготовить слайды"}` и создаёт задачу;
- `GET /api/tasks.xml` отдаёт тот же список в XML — чтобы сравнить форматы;
- `POST /upload` принимает один файл из поля `file` в `multipart/form-data` и сохраняет его в `uploads/`;
- сервис запускается по HTTP, а с флагом `-tls` — по HTTPS с самоподписанным сертификатом;
- отдельный маленький клиент получает список задач;
- два теста проверяют HTTP-хендлер без запуска реального сервера.

Это ровно та функциональность, которую стоит писать вживую: один файл сервера, один шаблон, один клиент и один тестовый файл.

## Структура проекта

```text
api/openapi.yaml              # контракт API
cmd/server/main.go            # запуск HTTP/TLS-сервера
cmd/client/main.go            # маленький HTTP-клиент
internal/tasks/service.go     # задачи и хранение в памяти
internal/httpapi/handler.go   # маршруты, JSON, XML, multipart
internal/httpapi/handler_test.go
web/templates/index.html      # HTML-шаблон
```

`cmd` запускает программу, `internal/tasks` хранит состояние, а `internal/httpapi` переводит HTTP-запросы в операции с задачами. Так понятно, где искать код, и при этом для семинара нет лишних интерфейсов, репозиториев и слоёв.

## Версии для живого кодинга

В [steps](steps) каждая следующая папка — самостоятельная версия сервиса. Открывайте именно её `internal/httpapi/handler.go`, а не финальный хендлер в корне.

```text
01_get_tasks       GET /api/tasks и JSON
02_create_task     POST /api/tasks и JSON body
03_xml             XML-представление
04_html_template   HTML и html/template
05_upload          multipart-загрузка файлов
06_httptest        тест готового хендлера
07_tls             запуск готового хендлера через TLS
```

Запуск любой версии из корня семинара:

```bash
go run ./steps/01_get_tasks/cmd/server
go run ./steps/02_create_task/cmd/server
go run ./steps/03_xml/cmd/server
go run ./steps/04_html_template/cmd/server
go run ./steps/05_upload/cmd/server
go test ./steps/06_httptest/internal/httpapi
go run ./steps/07_tls/cmd/server -tls
```

Шаги 06 и 07 намеренно не изменяют маршруты: тест проверяет уже написанный HTTP-код, а TLS меняет способ запуска `http.Server`, а не сами хендлеры.

## Swagger / OpenAPI

Схема уже лежит в [api/openapi.yaml](api/openapi.yaml). Откройте её в Swagger Editor или любом OpenAPI Viewer.

Идите по ней в таком порядке:

1. `components/schemas/Task` — создаём Go-структуру `Task` и JSON-теги.
2. `GET /api/tasks` — создаём `ServeMux`, первый хендлер и `writeJSON`.
3. `POST /api/tasks` — добавляем структуру запроса, декодирование JSON, проверку пустого `text` и статус `201`.
4. `GET /api/tasks.xml` — добавляем XML-теги и второй формат ответа.
5. `POST /upload` — по схеме видно `multipart/form-data`, поле `file` и ограничение 5 MiB.
6. После реализации используйте блоки `example` в схеме как готовые запросы и ожидаемые ответы.

## Запуск

Откройте два терминала.

```bash
cd http-seminar
go run ./cmd/server
```

Во втором:

```bash
go run ./cmd/client
curl -i http://localhost:8080/api/tasks
curl -i -X POST http://localhost:8080/api/tasks \\
  -H 'Content-Type: application/json' \\
  -d '{"text":"Показать httptest"}'
curl http://localhost:8080/api/tasks.xml
curl -F 'file=@README.md' http://localhost:8080/upload
go test ./...
```

Откройте `http://localhost:8080/` в браузере, чтобы посмотреть HTML-шаблон.

## TLS: локальная демонстрация

Создайте сертификат один раз в папке семинара:

```bash
openssl req -x509 -newkey rsa:2048 -nodes \\
  -keyout key.pem -out cert.pem -days 1 \\
  -subj '/CN=localhost' \\
  -addext 'subjectAltName=DNS:localhost'
go run ./cmd/server -tls
```

Браузер предупредит о недоверенном сертификате — это ожидаемо. Для проверки из терминала используйте `-k` только в этом локальном примере:

```bash
curl -k https://localhost:8080/api/tasks
```
