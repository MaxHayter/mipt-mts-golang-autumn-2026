# Шаг 7: HTTPS с TLS

## Что добавили

Маршруты не меняются. Изменился запуск сервера: сначала создаём TCP-listener через `net.Listen`, затем при флаге `-tls` вызываем `server.ServeTLS` с сертификатом и ключом.

Посмотрите: `cmd/server/main.go`. Это важный вывод: TLS защищает транспорт, а хендлеры остаются прежними.

## Подготовка сертификата

Из корня `http-seminar` создайте локальный самоподписанный сертификат:

```bash
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout key.pem -out cert.pem -days 1 \
  -subj '/CN=localhost' \
  -addext 'subjectAltName=DNS:localhost'
```

## Запуск и проверка

```bash
go run ./steps/07_tls/cmd/server -tls
```

Во втором терминале:

```bash
curl -k -i https://localhost:8080/api/tasks
```

`-k` отключает проверку доверия к сертификату и допустим только для этой локальной демонстрации. В браузере будет предупреждение о недоверенном сертификате — это ожидаемо.
