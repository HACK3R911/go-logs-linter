# go-logs-linter

Инструмент для статического анализа лог-сообщений в Go-коде. Проверяет сообщения логов на соответствие стилевым правилам и выявляет потенциальные проблемы.

## Возможности

- **Проверка регистра** - сообщения должны начинаться со строчной буквы
- **Проверка языка** - только английские символы в сообщениях
- **Проверка специальных символов и эмодзи** - запрет на неалфавитные символы
- **Обнаружение чувствительных данных** - выявление паролей, токенов, ключей и т.д.
- **Пользовательские паттерны** - поддержка разрешенных/запрещенных регулярных выражений
- **Поддержка библиотек** - zap, slog
- **Интеграция** - плагин для golangci-lint

## Установка

```bash
go install github.com/HACK3R911/go-logs-linter/cmd/loglint@latest
```

Или соберите из исходников:

```bash
git clone https://github.com/HACK3R911/go-logs-linter.git
cd go-logs-linter
go build -o bin/loglint ./cmd/loglint
```

## Сборка

```bash
make build
```

## Запуск

### Без конфигурации (используются настройки по умолчанию)

```bash
loglint ./path/to/your/code
```

### С указанием конфигурации

```bash
loglint -config=config.yaml ./path/to/your/code
```

## Конфигурация

Создайте файл `config.yaml`:

```yaml
rules:
  # Разрешить сообщения с заглавной буквы
  allow_uppercase_start: false

  # Разрешенные паттерны (сообщение должно соответствовать хотя бы одному)
  allowed_patterns: []

  # Запрещенные паттерны
  disallowed_patterns: []

  # Разрешить неанглийские символы
  allow_non_english: false

  # Разрешить специальные символы
  allow_special_chars: false

  # Разрешить чувствительные данные
  allow_sensitive_data: false

  # Ключевые слова для определения чувствительных данных
  sensitive_keywords:
    - password
    - token
    - secret
    - key
    - api_key
    - apikey
    - private
    - credential
```

## Примеры использования

### Проверка файла

```bash
loglint ./main.go
```

### Проверка директории

```bash
loglint ./pkg/
```

### Проверка с настройками

```bash
loglint -config=./config.yaml ./src/
```

## Примеры сообщений

### ✅ Правильные сообщения

```go
logger.Info("starting server on port 8080")
logger.Error("failed to connect to database")
logger.Warn("something went wrong")
logger.Debug("api request completed")
```

### ❌ Неправильные сообщения

```go
// Начинается с заглавной буквы
logger.Info("Starting server on port 8080")

// Содержит неанглийские символы
logger.Info("Запуск сервера")

// Содержит специальные символы
logger.Info("server started!🚀")

// Содержит чувствительные данные
logger.Info("user password: " + password)
```

## Тесты проекта

```bash
make test
```

## Линтинг проекта

```bash
make lint
```

## CI/CD

Инструмент можно интегрировать в GitHub Actions:

```yaml
- name: Run loglint
  uses: golangci/golangci-lint-action@v6
  with:
    args: --enable loglint
```