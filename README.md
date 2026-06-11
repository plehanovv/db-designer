# DB Designer VKR

Прототип программной подсистемы для автоматизированного проектирования структуры базы данных на основе описания предметной области.

Система принимает текстовое описание, JSON-модель или CSV-модель и формирует первичную логическую структуру базы данных: сущности, атрибуты, связи, ER-представление, диагностические сообщения и SQL DDL-код для PostgreSQL.

## Назначение

Проект предназначен для демонстрации начального этапа автоматизированного проектирования базы данных:

1. пользователь задает название базы данных в отдельном поле интерфейса;
2. пользователь вводит описание предметной области;
3. подсистема выделяет сущности, атрибуты и связи;
4. результат отображается в виде модели, ER-схемы и SQL-кода;
5. пользователь может вручную скорректировать модель и повторно сгенерировать SQL.

Важно: демонстрационные тексты не должны начинаться со строки `База данных ...`, потому что название базы данных задается отдельно в поле **База данных**.

## Состав проекта

```text
cmd/server/                 запуск Go web-сервера
internal/handler/           HTTP-обработчики
internal/service/           сервисный слой и сценарии анализа
internal/nlp/               клиент Python NLP-сервиса и локальный fallback
internal/analyzer/          rule-based анализ: сущности, атрибуты, связи, нормализация
internal/generator/         генерация PostgreSQL DDL
internal/storage/           опциональное сохранение результатов в PostgreSQL
web/                        веб-интерфейс и демонстрационные профили
examples/                   рабочие входные примеры для показа и тестов
docs/                       дополнительные схемы и пояснения
```

## Требования

Минимально для запуска веб-приложения:

- Go 1.22 или новее;
- любой современный браузер.

Для полного режима с NLP-сервисом:

- Python 3.10 или новее;
- зависимости из `D:\Python\nlp-service\requirements.txt`;
- модели spaCy `ru_core_news_sm` и `en_core_web_sm`.

Опционально:

- PostgreSQL, если нужно сохранять результаты анализа в базу данных.

Без Python NLP-сервиса приложение все равно запускается: Go backend использует локальный rule-based fallback. Для защиты лучше запускать Python NLP-сервис, чтобы в диагностике отображалось использование spaCy.

## Быстрый запуск без Python NLP

Этот вариант подходит, если нужно быстро открыть интерфейс на ПК.

```powershell
cd D:\go\proejcts\db-designer-vkr
go run -buildvcs=false ./cmd/server
```

Открыть в браузере:

```text
http://localhost:8080
```

Если порт `8080` занят:

```powershell
$env:PORT="8090"
go run -buildvcs=false ./cmd/server
```

Открыть:

```text
http://localhost:8090
```

## Запуск с Python NLP-сервисом

Сначала запустить NLP-сервис в отдельном терминале:

```powershell
cd D:\Python\nlp-service
.\.venv\Scripts\python.exe -m uvicorn main:app --reload --port 8000
```

Проверить, что сервис отвечает:

```powershell
Invoke-RestMethod http://localhost:8000/health
```

После этого запустить Go-приложение во втором терминале:

```powershell
cd D:\go\proejcts\db-designer-vkr
go run -buildvcs=false ./cmd/server
```

Если Python-сервис запущен на другом адресе:

```powershell
$env:NLP_SERVICE_URL="http://localhost:8000/analyze"
go run -buildvcs=false ./cmd/server
```

## Установка Python NLP-сервиса с нуля

Если на ПК нет готового виртуального окружения:

```powershell
cd D:\Python\nlp-service
python -m venv .venv
.\.venv\Scripts\python.exe -m pip install --upgrade pip
.\.venv\Scripts\python.exe -m pip install -r requirements.txt
.\.venv\Scripts\python.exe -m spacy download ru_core_news_sm
.\.venv\Scripts\python.exe -m spacy download en_core_web_sm
```

Запуск:

```powershell
.\.venv\Scripts\python.exe -m uvicorn main:app --reload --port 8000
```

## Сборка исполняемого файла

Чтобы не запускать проект через `go run`, можно собрать `.exe`:

```powershell
cd D:\go\proejcts\db-designer-vkr
go build -buildvcs=false -o .\db-designer-vkr.exe .\cmd\server
.\db-designer-vkr.exe
```

После запуска открыть:

```text
http://localhost:8080
```

## Переменные окружения

| Переменная | Назначение | Значение по умолчанию |
|---|---|---|
| `PORT` | порт Go web-сервера | `8080` |
| `NLP_SERVICE_URL` | адрес Python NLP-сервиса | `http://localhost:8000/analyze` |
| `DATABASE_URL` | строка подключения PostgreSQL | пусто |
| `DB_DESIGNER_DATABASE_URL` | альтернативная строка подключения PostgreSQL | пусто |
| `PSQL_PATH` | путь к `psql.exe` для проверки SQL | поиск в `PATH` |

PostgreSQL не обязателен. Если строка подключения не задана, приложение работает без сохранения истории анализа.

## Демонстрационные примеры

Основные примеры находятся в `web/domain-examples.json` и отображаются кнопками в интерфейсе:

- **Библиотека** (`LibraryDemo`);
- **Учебный процесс** (`UniversityProcess`);
- **Интернет-магазин** (`CommerceDemo`).

Те же рабочие данные лежат в папке `examples/`:

```text
examples/domain_description.txt      текстовый пример библиотеки
examples/university_process.txt      текстовый пример учебного процесса
examples/structured_model.json       структурированная JSON-модель библиотеки
examples/structured_model.csv        структурированная CSV-модель интернет-магазина
```

Пример текстового ввода:

```text
Читатель имеет имя email телефон.
Автор имеет имя страну.
Книга имеет название год isbn.
Категория имеет название код.
Выдача имеет дату выдачи дату возврата статус.
Бронирование имеет дату статус.
Штраф имеет сумму дату статус.
Автор пишет книги.
Категория включает книги.
Читатель имеет несколько выдач.
Выдача связана с книгой.
Читатель бронирует книги.
Бронирование связано с книгой.
Выдача содержит штрафы.
```

Название базы данных для этого примера задается отдельно в интерфейсе: `LibraryDemo`.

## Проверка проекта

Запуск всех Go-тестов:

```powershell
cd D:\go\proejcts\db-designer-vkr
$env:GOCACHE="D:\go\proejcts\db-designer-vkr\.gocache"
go test ./...
Remove-Item -LiteralPath ".\.gocache" -Recurse -Force -ErrorAction SilentlyContinue
```

Контрольная сборка:

```powershell
go build -buildvcs=false -o .\db-designer-vkr.exe .\cmd\server
```
