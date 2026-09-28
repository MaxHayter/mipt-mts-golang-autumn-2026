## Переменные среды разработки

### Показать текущее окружение
```code:bash
go env
```

### Показать конкретную переменную
```code:bash
go env GOROOT
go env GOOS GOARCH
```

### Можно менять переменные для кроссплатформенной сборки
```code:bash
GOOS=linux GOARCH=amd64 go build
```

```code:bash
go tool dist list
```

```code:bash
// Windows: 
set GOOS=linux && set GOARCH=amd64 && go build
```

## Управление зависимостями

### Инициализация нового модуля
```code:bash
go mod init myapp
```

### Добавление зависимостей
#### Загрузить конкретную версию
```code:bash
go get github.com/google/uuid@v1.3.0
```

#### Загрузить последнюю версию
```code:bash
go get github.com/google/uuid@latest
```

#### Загрузить версию ветки
```code:bash
go get github.com/google/uuid@main
```

#### Обновить все зависимости
```code:bash
go get ./...
```

### Просмотр зависимостей
#### Список всех модулей
```code:bash
go list -m all
```

#### Информация о конкретном пакете
```code:bash
go list -m github.com/google/uuid
```

#### Список с версиями и обновлениями
```code:bash
go list -m -u all
```

#### Обслуживание зависимостей
##### Удалить неиспользуемые зависимости
```code:bash
go mod tidy
```

#### Проверить целостность
```code:bash
go mod verify
```

### Работа с кодом
#### Показать ошибки со статическим анализом
```code:bash
go vet ./...
```

#### Показать форматирование
```code:bash
go fmt ./...
```
или
```code:bash
gofmt -w main.go
```

#### Генерировать документацию
```code:bash
go doc CalculateSum
```

#### Вся документация пакета
```code:bash
go doc -all
```

### Юнит-тест
#### Запустить все тесты
```code:bash
go test ./...
```

#### Подробный вывод
```code:bash
go test -v ./...
```

#### Запустить конкретный тест
```code:bash
go test -run TestAdd
```

#### Запустить только бенчмарки
```code:bash
go test -bench=. -run=^$
```

#### Бенчмарк с конкретным количеством итераций
```code:bash
go test -bench=BenchmarkFibonacci -benchtime=10s
```

#### Покрытие кода
```code:bash
go test -cover ./...
```

#### Детальный отчёт покрытия (html)
```code:bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Мониторинг производительности бенчмарка
```code:bash
go test -bench=BenchmarkFibonacci -benchmem
```

### Профилирование

#### При запуске программы
```code:bash
go run -cpuprofile=cpu.prof ./main.go
go run -memprofile=mem.prof ./main.go
```

#### После компиляции
```code:bash
go build -o myapp
./myapp -cpuprofile=cpu.prof
./myapp -memprofile=mem.prof
```

#### Из бенчмарков
```code:bash
go test -bench=. -cpuprofile=cpu.prof
go test -bench=. -memprofile=mem.prof
```

#### Интерактивный анализ (команда top 10)
```code:bash
go tool pprof cpu.prof
```

#### Запустить веб-интерфейс (автоматически откроет браузер)
```code:bash
go tool pprof -http=:6060 cpu.prof
```

#### С указанием метрики (по умолчанию inuse_space)
```code:bash
go tool pprof -http=:6061 -alloc_space mem.prof
```

#### Или указать явно (без автооткрытия)
```code:bash
go tool pprof -http=localhost:6060 cpu.prof
```
#### Снятие профилей на лету
```code:bash
# CPU профиль (30 секунд)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory профиль (текущее состояние)
go tool pprof http://localhost:6060/debug/pprof/heap

# Горутины
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Блокировки (mutexes)
go tool pprof http://localhost:6060/debug/pprof/mutex

# Все доступные профили
go tool pprof http://localhost:6060/debug/pprof/

# В браузер (откроется localhost:6060)
go tool pprof -http=:6060 http://localhost:6060/debug/pprof/heap
```

### Трассировка

#### При запуске программы
```code:bash
go run -trace=trace.out ./main.go
```

# Из бенчмарков
```code:bash
go test -bench=. -trace=trace.out
```

# После компиляции
```code:bash
./myapp -trace=trace.out
```

# Открыть интерактивный просмотр
```code:bash
go tool trace trace.out
```


# Чего нет в коробке (5 мин)

- Линтер          golangci-lint   https://golangci-lint.run/
- Отладчик	 Delve	           https://github.com/go-delve/delve
- LSP-сервер  gopls	           https://github.com/golang/tools/tree/master/gopls
