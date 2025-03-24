# Ипотечный калькулятор

[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org/dl/)
![Repository Top Language](https://img.shields.io/github/languages/top/evt/callback)

### Сервис для расчета параметров ипотеки с кэшированием результатов.

## Возможности

- Расчет ключевых параметров ипотеки:
  - Процентная ставка (в зависимости от программы кредитования)
  - Сумма кредита
  - Ежемесячный аннуитетный платеж
  - Общая переплата за весь срок
  - Дата последнего платежа
- Поддержка трех программ кредитования:
  1. Корпоративная программа (8%)
  2. Военная ипотека (9%)
  3. Базовая программа (10%)
- Валидация входных данных:
  - Проверка минимального первоначального взноса (20%)
  - Проверка выбора только одной программы
- Кэширование результатов расчетов в памяти
- Логирование запросов через middleware

## API

### `POST /execute`

Производит расчет параметров ипотеки.

**Входные данные:**
```json
{
    "object_cost": 5000000,
    "initial_payment": 1000000,
    "months": 240,
    "program": {
        "salary": true
    }
}
```

**Успешный ответ (200 OK):**
```json
{
   "result": {
      "params": {
         "object_cost": 5000000,
         "initial_payment": 1000000,                
         "months": 240
      },
      "program": {
         "salary": true
      },
      "aggregates": {
         "rate": 8,
         "loan_sum": 4000000,
         "monthly_payment": 33458,
         "overpayment": 4029920,
         "last_payment_date": "2044-02-18"
      }
   }
}
```

**Возможные ошибки:**
- 400 Bad Request:
  - `{"error": "choose program"}` - не выбрана программа
  - `{"error": "choose only 1 program"}` - выбрано несколько программ
  - `{"error": "the initial payment should be more"}` - недостаточный первоначальный взнос

### `GET /cache`

Возвращает все сохраненные в кэше расчеты.

**Успешный ответ (200 OK):**
```json
[
   {
      "id": 0,
      "params": {
         "object_cost": 5000000,
         "initial_payment": 1000000,
         "months": 240
      },
      "program": {
         "salary": true
      },
      "aggregates": {
         "rate": 8,
         "loan_sum": 4000000,
         "monthly_payment": 33458,
         "overpayment": 4029920,
         "last_payment_date": "2044-02-18"
      }
   }
]
```

**Ошибка (400 Bad Request):**
```json
{
   "error": "empty cache"
}
```

## Установка и запуск

### Требования
- Go 1.21+
- Docker (для сборки контейнера)

### Сборка и запуск

1. Клонировать репозиторий:
```bash
git clone https://github.com/yourusername/mortgage-calculator.git
cd mortgage-calculator
```

2. Запустить с помощью Makefile:
```bash
make run
```

Или вручную:
```bash
go run cmd/main.go
```

### Сборка Docker-образа
```bash
make build
```

### Запуск контейнера
```bash
make start
```

### Остановка контейнера
```bash
make stop
```

## Тестирование

Запуск тестов:
```bash
make test
```

Запуск проверки конфигурацией линтера:
```bash
make lint
```

## Конфигурация

Настройки сервиса хранятся в файле `./internal/config/config.yml`:
```yaml
s
port: 8080
```

## Технические детали

- Используется стандартный кэш в памяти (не требует внешних БД)
- Реализован middleware для логирования запросов
- Поддержка graceful shutdown
- Оптимизированный Docker-образ (<30MB)
- Полное покрытие unit-тестами (>80%)
- Проверка кода golangci-lint

