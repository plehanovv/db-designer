package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnalyzeHandlerReturnsModel(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/analyze",
		strings.NewReader(`{"text":"Customer has name email phone. Order has number date amount. Order belongs customer."}`),
	)
	recorder := httptest.NewRecorder()

	AnalyzeHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"entities"`) || !strings.Contains(body, "FOREIGN KEY") {
		t.Fatalf("expected response body to contain entities and SQL, got: %s", body)
	}
}
