# Demonstration Inputs

Use these files for defense and live checks. They are written in a controlled natural-language style: one sentence describes one entity, attribute list, or relation.

## `domain_description.txt`

Domain: library.

Expected entities:

- `Читатель`
- `Автор`
- `Книга`
- `Категория`
- `Выдача`
- `Бронирование`
- `Штраф`

Expected key relations:

- `Автор -> Книга`, `one-to-many`
- `Категория -> Книга`, `one-to-many`
- `Читатель -> Выдача`, `one-to-many`
- `Выдача -> Книга`, `many-to-one`
- `Бронирование -> Книга`, `many-to-one`
- `Выдача -> Штраф`, `one-to-many`

## `university_process.txt`

Domain: university process.

Expected entities:

- `Студент`
- `Преподаватель`
- `Факультет`
- `Кафедра`
- `Курс`
- `Группа`
- `Экзамен`
- `Результат`

Expected key relations:

- `Факультет -> Кафедра`, `one-to-many`
- `Кафедра -> Курс`, `one-to-many`
- `Студент -> Группа`, `many-to-one`
- `Студент -> Курс`, `many-to-many`
- `Преподаватель -> Курс`, `many-to-many`
- `Студент -> Экзамен`, `many-to-many`
- `Экзамен -> Курс`, `many-to-one`
- `Преподаватель -> Результат`, `one-to-many`
- `Результат -> Экзамен`, `many-to-one`

## Structured Inputs

- `structured_model.json`: structured library model, NLP skipped.
- `structured_model.csv`: structured commerce model, NLP skipped.

Structured inputs are the safest fallback if the NLP service is unavailable during a live demonstration.
