# Semantic Core

The project does not use an LLM. The analysis is an explainable NLP and rule-based semantic pipeline.

## Pipeline

1. Input is accepted as natural text, JSON, or CSV.
2. JSON and CSV are parsed directly into the logical model.
3. Natural text is sent to the NLP service when available.
4. The NLP service uses `ru_core_news_sm` for Russian and `en_core_web_sm` for English.
5. If the NLP service is unavailable or weak, the Go local rule-based fallback is used.
6. Tokens, lemmas, POS tags, dependency metadata, and sentence boundaries are normalized.
7. Rule templates extract entities, attributes, relations, cardinalities, and constraints.
8. The logical model is validated.
9. The logical model is transformed into a physical PostgreSQL schema.
10. Generated SQL receives an internal DDL sanity check.

## Logical Model

The logical model is the database structure before physical SQL generation:

```json
{
  "database": {"name": "Shop"},
  "entities": [
    {
      "name": "Customer",
      "attributes": [
        {"name": "email", "type": "VARCHAR(255)", "required": true, "unique": true}
      ]
    }
  ],
  "relations": [
    {"from": "Order", "to": "Customer", "type": "belongs_to", "cardinality": "many-to-one"}
  ]
}
```

## Physical Model

The physical model is generated PostgreSQL DDL:

- `CREATE DATABASE`;
- `CREATE TABLE`;
- generated `id SERIAL PRIMARY KEY`;
- scalar columns;
- `NOT NULL`;
- `UNIQUE`;
- foreign keys;
- many-to-many junction tables;
- FK indexes.

## Entity Rules

Entities are extracted mostly from nouns and proper nouns, then filtered through:

- service/domain stop words;
- known scalar attribute terms;
- contextual noise rules;
- normalization dictionaries for Russian cases and plurals.

Supported examples:

```text
В системе хранятся клиенты, заказы и товары.
У нас есть товары которые имеют цену скидку количество на складе.
Пациент имеет несколько приемов.
```

## Attribute Rules

Supported attribute patterns:

```text
Клиент имеет имя email телефон.
У клиента есть имя и телефон.
Пользователи у которых есть номер телефона почта дата последней покупки.
Entity has name email phone.
```

Compound date attributes are normalized:

- `дата рождения` -> `data_rozhdenie DATE`;
- `дата создания` -> `data_sozdanie DATE`;
- `дата последней покупки` -> `data_pokupka DATE`;
- `дата начала` -> `data_nachalo DATE`;
- `дата окончания` -> `data_okonchanie DATE`.

## Type Rules

Examples:

- `email`, `почта` -> `VARCHAR(255)`;
- `phone`, `телефон` -> `VARCHAR(20)`;
- `date`, `дата_*` -> `DATE`;
- `time`, `время` -> `TIME`;
- `price`, `amount`, `цена`, `сумма`, `скидка` -> `NUMERIC(12,2)`;
- `number`, `количество`, `год`, `этаж` -> `INTEGER`;
- unknown scalar field -> `TEXT`.

## Constraint Rules

Default and explicit constraints are supported:

```text
Customer has required unique email phone.
Пользователь имеет обязательную уникальную почту телефон.
```

Examples:

- `email`, `почта`, `login`, `логин` are required by default;
- `email`, `почта`, `login`, `код`, `госномер`, `паспорт` are unique by default;
- `required`, `mandatory`, `обязательный` set `NOT NULL`;
- `unique`, `уникальный` set `UNIQUE`.

Constraint marker words are not emitted as columns.

## Relation Rules

Supported relation classes:

- ownership: `заказ принадлежит клиенту`, `order belongs to customer`;
- containment: `заказ содержит товары`, `категория включает товары`;
- action: `клиент оформляет заказы`, `пользователь создает заявки`;
- assignment/association: `заявка назначена сотруднику`, `история связана с пользователем`;
- organizational membership: `сотрудник работает в отделе`;
- categorization: `товар относится к категории`;
- educational many-to-many: `student enrolls in courses`, `teacher teaches courses`, `студент записывается на курсы`, `преподаватель ведет курсы`.

Cardinality mapping:

- `belongs_to`, `работает в`, `относится к` -> `many-to-one`;
- `contains`, `includes`, `оформляет`, `создает`, `покупает` -> `one-to-many`;
- `enrolls`, `teaches`, `записывается`, `ведет` -> `many-to-many`;
- unknown association -> `unspecified`, editable by the user.

## Transformation Report

The response includes `transformations`, explaining the conversion:

- `database_name_to_create_database`;
- `entity_to_table`;
- `attribute_to_column`;
- `relation_to_foreign_key`;
- `many_to_many_to_junction_table`;
- `foreign_key_index`.

This connects the semantic/logical model with the physical SQL schema.

## Validation

Model validation checks:

- empty database or entity names;
- duplicate entities;
- duplicate attributes;
- unsupported SQL types;
- relation endpoints that do not match entities;
- unsupported cardinalities;
- SQL identifier collisions.

SQL sanity validation checks:

- empty SQL;
- unbalanced parentheses;
- missing semicolons;
- duplicate `CREATE TABLE`;
- FK references to missing tables;
- duplicate indexes;
- indexes on missing tables.

## Boundaries

The system supports controlled natural language, not unlimited human-level understanding. The expected input should contain recognizable markers such as:

- `имеет`, `есть`;
- `содержит`, `включает`;
- `связан`, `назначен`;
- `принадлежит`, `относится`;
- `оформляет`, `создает`;
- `записывается`, `ведет`.

The advantage of this approach is that every result is explainable, deterministic, testable, and manually correctable.
