package handler

import (
	"encoding/json"
	"net/http"

	"db-designer-vkr/internal/model"
	"db-designer-vkr/internal/service"
	"db-designer-vkr/internal/storage"
)

var resultStore *storage.Store

type AnalyzeRequest struct {
	Text     string         `json:"text"`
	Database model.Database `json:"database,omitempty"`
}

type GenerateSQLRequest struct {
	Database  model.Database   `json:"database"`
	Entities  []model.Entity   `json:"entities"`
	Relations []model.Relation `json:"relations"`
}

type GenerateSQLResponse struct {
	SQL             string                     `json:"sql"`
	Diagnostics     []model.Diagnostic         `json:"diagnostics"`
	Transformations []model.TransformationStep `json:"transformations"`
	StorageKey      string                     `json:"storageKey,omitempty"`
}

func SetStore(store *storage.Store) {
	resultStore = store
}

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	result, err := service.AnalyzeTextWithDatabase(req.Text, req.Database)
	if err != nil {
		if service.IsInputError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result.StorageKey = saveAnalysisResult(req.Text, inputType(req.Text), result.Database, result.Entities, result.Relations, result.SQL, result.Diagnostics, result.Transformations, &result.Diagnostics)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GenerateSQLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateSQLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	result, diagnostics, transformations := service.GenerateSQLWithMetadata(req.Database, req.Entities, req.Relations)
	storageKey := saveAnalysisResult(modelAsInput(req), "manual_edit", req.Database, req.Entities, req.Relations, result, diagnostics, transformations, &diagnostics)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenerateSQLResponse{SQL: result, Diagnostics: diagnostics, Transformations: transformations, StorageKey: storageKey})
}

func saveAnalysisResult(
	inputText string,
	inputType string,
	database model.Database,
	entities []model.Entity,
	relations []model.Relation,
	sql string,
	diagnostics []model.Diagnostic,
	transformations []model.TransformationStep,
	responseDiagnostics *[]model.Diagnostic,
) string {
	if resultStore == nil {
		return ""
	}

	key, err := resultStore.SaveAnalysis(storage.AnalysisRecord{
		InputText:       inputText,
		InputType:       inputType,
		Database:        database,
		Entities:        entities,
		Relations:       relations,
		SQL:             sql,
		Diagnostics:     diagnostics,
		Transformations: transformations,
	})
	if err != nil {
		*responseDiagnostics = append(*responseDiagnostics, model.Diagnostic{
			Level:   "warning",
			Message: "Analysis result was not saved to PostgreSQL storage: " + err.Error(),
		})
		return ""
	}

	*responseDiagnostics = append(*responseDiagnostics, model.Diagnostic{
		Level:   "info",
		Message: "Analysis result was saved to PostgreSQL storage with key " + key + ".",
	})
	return key
}

func inputType(text string) string {
	if service.LooksLikeJSON(text) {
		return "json"
	}
	if service.LooksLikeCSV(text) {
		return "csv"
	}
	return "text"
}

func modelAsInput(req GenerateSQLRequest) string {
	value, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	return string(value)
}
