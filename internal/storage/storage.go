package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"db-designer-vkr/internal/model"
)

type Store struct {
	databaseURL string
	psqlPath    string
}

type AnalysisRecord struct {
	InputText       string
	InputType       string
	Database        model.Database
	Entities        []model.Entity
	Relations       []model.Relation
	SQL             string
	Diagnostics     []model.Diagnostic
	Transformations []model.TransformationStep
}

func NewFromEnv() (*Store, bool) {
	databaseURL := strings.TrimSpace(os.Getenv("DB_DESIGNER_DATABASE_URL"))
	if databaseURL == "" {
		databaseURL = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	}
	if databaseURL == "" {
		return nil, false
	}

	return &Store{
		databaseURL: databaseURL,
		psqlPath:    resolvePSQLPath(),
	}, true
}

func (s *Store) Init() error {
	if s == nil {
		return nil
	}
	return s.runSQL(schemaSQL)
}

func (s *Store) SaveAnalysis(record AnalysisRecord) (string, error) {
	if s == nil {
		return "", nil
	}

	key, err := newStorageKey()
	if err != nil {
		return "", err
	}

	sql, err := buildSaveSQL(key, record)
	if err != nil {
		return "", err
	}
	if err := s.runSQL(sql); err != nil {
		return "", err
	}

	return key, nil
}

func (s *Store) runSQL(sql string) error {
	if strings.TrimSpace(s.psqlPath) == "" {
		return fmt.Errorf("psql executable was not found")
	}

	path := filepath.Join(os.TempDir(), "db-designer-storage-*.sql")
	file, err := os.CreateTemp("", filepath.Base(path))
	if err != nil {
		return err
	}
	name := file.Name()
	defer os.Remove(name)

	if _, err := file.WriteString(sql); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	cmd := exec.Command(s.psqlPath, "-X", "-v", "ON_ERROR_STOP=1", s.databaseURL, "-f", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("psql failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func buildSaveSQL(key string, record AnalysisRecord) (string, error) {
	modelJSON, err := json.Marshal(struct {
		Database  model.Database   `json:"database"`
		Entities  []model.Entity   `json:"entities"`
		Relations []model.Relation `json:"relations"`
	}{
		Database:  record.Database,
		Entities:  record.Entities,
		Relations: record.Relations,
	})
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString("BEGIN;\n")
	builder.WriteString("INSERT INTO text_description (storage_key, content, input_type, database_name, model_json) VALUES (")
	builder.WriteString(sqlString(key))
	builder.WriteString(", ")
	builder.WriteString(sqlString(record.InputText))
	builder.WriteString(", ")
	builder.WriteString(sqlString(defaultString(record.InputType, "text")))
	builder.WriteString(", ")
	builder.WriteString(sqlString(record.Database.Name))
	builder.WriteString(", ")
	builder.WriteString(sqlString(string(modelJSON)))
	builder.WriteString("::jsonb);\n")

	for _, entity := range record.Entities {
		builder.WriteString("INSERT INTO entity (storage_key, description_id, name, entity_type) VALUES (")
		builder.WriteString(sqlString(key))
		builder.WriteString(", (SELECT id FROM text_description WHERE storage_key = ")
		builder.WriteString(sqlString(key))
		builder.WriteString("), ")
		builder.WriteString(sqlString(entity.Name))
		builder.WriteString(", 'detected');\n")

		for _, attribute := range entity.Attributes {
			builder.WriteString("INSERT INTO attribute (storage_key, entity_id, name, data_type, required, unique_value) VALUES (")
			builder.WriteString(sqlString(key))
			builder.WriteString(", (SELECT id FROM entity WHERE storage_key = ")
			builder.WriteString(sqlString(key))
			builder.WriteString(" AND name = ")
			builder.WriteString(sqlString(entity.Name))
			builder.WriteString(" LIMIT 1), ")
			builder.WriteString(sqlString(attribute.Name))
			builder.WriteString(", ")
			builder.WriteString(sqlString(attribute.Type))
			builder.WriteString(", ")
			builder.WriteString(sqlBool(attribute.Required))
			builder.WriteString(", ")
			builder.WriteString(sqlBool(attribute.Unique))
			builder.WriteString(");\n")
		}
	}

	for _, relation := range record.Relations {
		builder.WriteString("INSERT INTO relation (storage_key, description_id, source_entity_id, target_entity_id, source_name, target_name, relation_type, cardinality) VALUES (")
		builder.WriteString(sqlString(key))
		builder.WriteString(", (SELECT id FROM text_description WHERE storage_key = ")
		builder.WriteString(sqlString(key))
		builder.WriteString("), (SELECT id FROM entity WHERE storage_key = ")
		builder.WriteString(sqlString(key))
		builder.WriteString(" AND name = ")
		builder.WriteString(sqlString(relation.From))
		builder.WriteString(" LIMIT 1), (SELECT id FROM entity WHERE storage_key = ")
		builder.WriteString(sqlString(key))
		builder.WriteString(" AND name = ")
		builder.WriteString(sqlString(relation.To))
		builder.WriteString(" LIMIT 1), ")
		builder.WriteString(sqlString(relation.From))
		builder.WriteString(", ")
		builder.WriteString(sqlString(relation.To))
		builder.WriteString(", ")
		builder.WriteString(sqlString(relation.Type))
		builder.WriteString(", ")
		builder.WriteString(sqlString(relation.Cardinality))
		builder.WriteString(");\n")
	}

	builder.WriteString("INSERT INTO database_schema (storage_key, description_id, sql_text, diagnostics_json, transformations_json) VALUES (")
	builder.WriteString(sqlString(key))
	builder.WriteString(", (SELECT id FROM text_description WHERE storage_key = ")
	builder.WriteString(sqlString(key))
	builder.WriteString("), ")
	builder.WriteString(sqlString(record.SQL))
	builder.WriteString(", ")
	diagnosticsJSON, _ := json.Marshal(record.Diagnostics)
	transformationsJSON, _ := json.Marshal(record.Transformations)
	builder.WriteString(sqlString(string(diagnosticsJSON)))
	builder.WriteString("::jsonb, ")
	builder.WriteString(sqlString(string(transformationsJSON)))
	builder.WriteString("::jsonb);\n")
	builder.WriteString("COMMIT;\n")

	return builder.String(), nil
}

func newStorageKey() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[:]), nil
}

func resolvePSQLPath() string {
	if path := strings.TrimSpace(os.Getenv("PSQL_PATH")); path != "" {
		return path
	}
	if path, err := exec.LookPath("psql"); err == nil {
		return path
	}

	candidates := []string{
		`C:\Program Files\PostgreSQL\18\bin\psql.exe`,
		`C:\Program Files\PostgreSQL\17\bin\psql.exe`,
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func sqlString(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func sqlBool(value bool) string {
	if value {
		return "TRUE"
	}
	return "FALSE"
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
