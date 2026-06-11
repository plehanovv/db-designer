package service

import (
	"fmt"
	"regexp"
	"strings"

	"db-designer-vkr/internal/model"
)

var (
	createTablePattern = regexp.MustCompile(`(?im)^\s*CREATE\s+TABLE\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
	createIndexPattern = regexp.MustCompile(`(?im)^\s*CREATE\s+INDEX\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+ON\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
	referencesPattern  = regexp.MustCompile(`(?im)\bREFERENCES\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
)

func ValidateSQL(sql string) []model.Diagnostic {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return []model.Diagnostic{warning("Generated SQL is empty.")}
	}

	var diagnostics []model.Diagnostic
	diagnostics = append(diagnostics, validateSQLParentheses(sql)...)
	diagnostics = append(diagnostics, validateSQLSemicolons(sql)...)
	diagnostics = append(diagnostics, validateSQLTablesAndReferences(sql)...)
	diagnostics = append(diagnostics, validateSQLIndexes(sql)...)

	if !hasDiagnosticLevel(diagnostics, "error") {
		diagnostics = append(diagnostics, model.Diagnostic{
			Level:   "info",
			Message: "Generated SQL passed internal DDL sanity checks.",
		})
	}

	return diagnostics
}

func validateSQLParentheses(sql string) []model.Diagnostic {
	balance := 0
	for _, r := range sql {
		switch r {
		case '(':
			balance++
		case ')':
			balance--
			if balance < 0 {
				return []model.Diagnostic{errorDiagnostic("Generated SQL has a closing parenthesis without a matching opening parenthesis.")}
			}
		}
	}
	if balance != 0 {
		return []model.Diagnostic{errorDiagnostic("Generated SQL has unbalanced parentheses.")}
	}
	return nil
}

func validateSQLSemicolons(sql string) []model.Diagnostic {
	var diagnostics []model.Diagnostic
	statements := strings.Split(sql, "\n\n")
	for _, statement := range statements {
		statement = strings.TrimSpace(removeSQLLineComments(statement))
		if statement == "" {
			continue
		}
		if !strings.HasSuffix(statement, ";") {
			diagnostics = append(diagnostics, errorDiagnostic("Generated SQL statement does not end with a semicolon: "+firstLine(statement)))
		}
	}
	return diagnostics
}

func validateSQLTablesAndReferences(sql string) []model.Diagnostic {
	var diagnostics []model.Diagnostic
	tables := make(map[string]bool)

	for _, match := range createTablePattern.FindAllStringSubmatch(sql, -1) {
		table := match[1]
		if tables[table] {
			diagnostics = append(diagnostics, errorDiagnostic(fmt.Sprintf("Generated SQL contains duplicate CREATE TABLE for %q.", table)))
		}
		tables[table] = true
	}

	for _, match := range referencesPattern.FindAllStringSubmatch(sql, -1) {
		table := match[1]
		if !tables[table] {
			diagnostics = append(diagnostics, errorDiagnostic(fmt.Sprintf("Foreign key references missing table %q.", table)))
		}
	}

	return diagnostics
}

func validateSQLIndexes(sql string) []model.Diagnostic {
	var diagnostics []model.Diagnostic
	indexes := make(map[string]bool)
	tables := make(map[string]bool)
	for _, match := range createTablePattern.FindAllStringSubmatch(sql, -1) {
		tables[match[1]] = true
	}

	for _, match := range createIndexPattern.FindAllStringSubmatch(sql, -1) {
		indexName := match[1]
		tableName := match[2]
		if indexes[indexName] {
			diagnostics = append(diagnostics, errorDiagnostic(fmt.Sprintf("Generated SQL contains duplicate CREATE INDEX %q.", indexName)))
		}
		if !tables[tableName] {
			diagnostics = append(diagnostics, errorDiagnostic(fmt.Sprintf("Index %q targets missing table %q.", indexName, tableName)))
		}
		indexes[indexName] = true
	}

	return diagnostics
}

func removeSQLLineComments(statement string) string {
	var lines []string
	for _, line := range strings.Split(statement, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func firstLine(value string) string {
	line := strings.Split(strings.TrimSpace(value), "\n")[0]
	if len(line) > 120 {
		return line[:120] + "..."
	}
	return line
}

func hasDiagnosticLevel(diagnostics []model.Diagnostic, level string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Level == level {
			return true
		}
	}
	return false
}
