package storage

const schemaSQL = `
CREATE TABLE IF NOT EXISTS text_description (
    id SERIAL PRIMARY KEY,
    storage_key VARCHAR(64) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    input_type VARCHAR(32) NOT NULL DEFAULT 'text',
    database_name VARCHAR(255),
    model_json JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS entity (
    id SERIAL PRIMARY KEY,
    storage_key VARCHAR(64) NOT NULL,
    description_id INTEGER NOT NULL REFERENCES text_description(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    entity_type VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS attribute (
    id SERIAL PRIMARY KEY,
    storage_key VARCHAR(64) NOT NULL,
    entity_id INTEGER NOT NULL REFERENCES entity(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(100),
    required BOOLEAN NOT NULL DEFAULT FALSE,
    unique_value BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS relation (
    id SERIAL PRIMARY KEY,
    storage_key VARCHAR(64) NOT NULL,
    description_id INTEGER NOT NULL REFERENCES text_description(id) ON DELETE CASCADE,
    source_entity_id INTEGER REFERENCES entity(id) ON DELETE SET NULL,
    target_entity_id INTEGER REFERENCES entity(id) ON DELETE SET NULL,
    source_name VARCHAR(255) NOT NULL,
    target_name VARCHAR(255) NOT NULL,
    relation_type VARCHAR(100),
    cardinality VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS database_schema (
    id SERIAL PRIMARY KEY,
    storage_key VARCHAR(64) NOT NULL,
    description_id INTEGER NOT NULL REFERENCES text_description(id) ON DELETE CASCADE,
    sql_text TEXT NOT NULL,
    diagnostics_json JSONB,
    transformations_json JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_entity_description_id ON entity(description_id);
CREATE INDEX IF NOT EXISTS idx_attribute_entity_id ON attribute(entity_id);
CREATE INDEX IF NOT EXISTS idx_relation_description_id ON relation(description_id);
CREATE INDEX IF NOT EXISTS idx_database_schema_description_id ON database_schema(description_id);
`
