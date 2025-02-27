# Metflix

<img src="gopher.png" style="width: 30%; float: left; margin: 0 15px 0px 0px; padding: 0 15px 0 0;" alt="Описание">

Учебный проект в рамках курса GO-advanced Яндекс.Практикума. Состоит из: 
- сервиса сбора метрик (агент)
- сервиса хранения метрик (сервер)
- утилита статического анализа кода (`multichecker`)

Для получения полного списка доступных команд выполните:
```bash
make help
```

<p style="clear: both">

## Миграции
Сервер автоматически применит новые миграции при запуске. 

Для работы с миграциями вручную можно установить утилиту [golang-migrate](https://github.com/golang-migrate/migrate):
```bash
go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# добавление новой миграции:
migrate create -ext sql -dir ./db/migrate -seq имя_миграции

# применить миграции:
migrate -database ${DATABASE_DSN} -path ./db/migrate up

# откатить миграции:
migrate -database ${DATABASE_DSN} -path ./db/migrate down -all
```


## Документация
Для генерации документации в формате OpenAPI (Swagger) необходимо установить `swag`:
```bash
go install github.com/swaggo/swag/cmd/swag@latest

make build # документация будет доступна на /swagger/index.html после запуска сервера
```

Также доступна документация в формате godoc: 
```bash
$ make godoc
# Project documentation is available at:
# http://127.0.0.1:3000/pkg/github.com/ex0rcist/metflix/pkg/
```

## Запуск сервера

Сервер отвечает за аггрегирование и хранение метрик. Для запуска выполните: 
```bash
./cmd/server/server
```

### Опции командной строки 
Имеют приоритет перед конфигурационным файлом. Для вывода списка доступных опций и их значений по умолчанию выполните команду:
```bash
./cmd/server/server --help

-a, --address string       address:port for HTTP API requests (default "0.0.0.0:8080")
-c, --config string        path to configuration file in JSON format
--crypto-key string    path to public key to encrypt agent -> server communications
-d, --database string      PostgreSQL database DSN
-r, --restore              whether to restore state on startup (default true)
-k, --secret string        a key to sign outgoing data
-f, --store-file string    path to file to store metrics
-i, --store-interval int   interval (s) for dumping metrics to the disk
-t, --trusted-subnet ipNet   trusted subnet in CIDR notation
```

### Переменные окружения сервера
Имеют приоритет перед опциями командной строки.

```bash
# Адрес и порт для http-api:
export ADDRESS=0.0.0.0:8080

# Адрес и порт для grpc-api:
export GRPC_ADDRESS=0.0.0.0:8080

# Доверенная подсеть в CIDR нотации (по умолчанию не задана).
export TRUSTED_SUBNET=

# Интервал времени в секундах для сохранения метрик на диск
# (значение 0 — делает запись синхронной):
export STORE_INTERVAL=300

# Имя файла, где хранятся значения метрик.
# Пустое значение — отключает функцию записи на диск:
export FILE_STORAGE_PATH="/tmp/devops-metrics-db.json"

# Восстанавливать ли сохраненные значения метрик из файла при старте сервера:
export RESTORE=true

# Секретный ключ для генерации подписи (по умолчанию не задан):
export KEY=

# Путь к приватному RSA ключу (в PEM формате) для расшифровки запросов агент -> сервер 
# (по умолчанию не задан):
export CRYPTO_KEY=

# DSN для подключения к базе данных (postgres-only):
export DATABASE_DSN=

# Адрес и порт, по которым доступен инструмент pprof:
export PROFILER_ADDRESS=0.0.0.0:8081

# Путь к конфигурационному файлу в JSON формате:
# Пример конфигурационного файла: ./config/server.example.json
export CONFIG=
```

## Запуск агента
Агент отвечает за сбор и отправку метрик на сервер. Для запуска выполните:
```bash
./cmd/agent/agent 
```

### Опции командной строки 
Имеют приоритет перед конфигурационным файлом. Посмотреть доступные оции: 
```bash
./cmd/agent/agent --help

-a, --address string        address:port for HTTP API requests (default "0.0.0.0:8080")
-c, --config string         path to configuration file in JSON format
    --crypto-key string         path to public key to encrypt agent -> server communications
-p, --poll-interval int     interval (s) for polling stats (default 2)
-l, --rate-limit int        number of max simultaneous requests to server (default -1)
-r, --report-interval int   interval (s) for polling stats (default 10)
-k, --secret string         a key to sign outgoing data
-t, --transport string      transport to use: http/grpc (default "http")
```

### Переменные окружения агента
Имеют приоритет перед опциями командной строки.

```bash
# Адрес и порт сервера, агрегирующего метрики:
export ADDRESS=0.0.0.0:8080

# Интервал опроса метрик (в секундах):
export POLL_INTERVAL=2

# Интервал отправки метрик (в секундах):
export REPORT_INTERVAL=10

# Секретный ключ для генерации подписи (по умолчанию не задан):
export KEY=

# Путь к публичному RSA ключу (в PEM формате) для шифрования запросов агент -> сервер
export CRYPTO_KEY=

# Путь к конфигурационному файлу в JSON формате (по умолчанию не задан):
# Пример конфигурационного файла: ./config/agent.example.json
export CONFIG=
```

## Запуск `multichecker`
```bash
./cmd/staticlint/staticlint <packages>
# или
./cmd/staticlint/staticlint ./...

# дополнительная настройка: 
./cmd/staticlint/staticlint -help
```