package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnalyzeHandlerReturnsModel(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

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
	if !strings.Contains(body, `"entities"`) ||
		!strings.Contains(body, `"explanation"`) ||
		!strings.Contains(body, `"transformations"`) ||
		!strings.Contains(body, "FOREIGN KEY") {
		t.Fatalf("expected response body to contain entities, explanation, transformations and SQL, got: %s", body)
	}
}

func TestAnalyzeHandlerUsesDatabaseFromRequest(t *testing.T) {
	t.Setenv("NLP_SERVICE_URL", "http://127.0.0.1:1/analyze")

	request := httptest.NewRequest(
		http.MethodPost,
		"/analyze",
		strings.NewReader(`{"text":"Reader has name email. Book has title author. Loan belongs book.","database":{"name":"Library"}}`),
	)
	recorder := httptest.NewRecorder()

	AnalyzeHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"database":{"name":"Library","domain":"Library"}`) ||
		!strings.Contains(body, "CREATE DATABASE library") {
		t.Fatalf("expected request database to be used in response and SQL, got: %s", body)
	}
}

func TestAnalyzeHandlerReturnsBadRequestForEmptyText(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/analyze", strings.NewReader(`{"text":"   "}`))
	recorder := httptest.NewRecorder()

	AnalyzeHandler(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "domain description is empty") {
		t.Fatalf("expected input validation message, got: %s", recorder.Body.String())
	}
}

func TestAnalyzeHandlerReturnsBadRequestForInvalidStructuredJSON(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/analyze", strings.NewReader(`{"text":"{}"}`))
	recorder := httptest.NewRecorder()

	AnalyzeHandler(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "structured JSON model must contain entities") {
		t.Fatalf("expected structured JSON validation message, got: %s", recorder.Body.String())
	}
}

func TestAnalyzeHandlerAcceptsStructuredCSVInput(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/analyze",
		strings.NewReader(`{"text":"kind,database,entity,attribute,type,required,unique,from,to,relation_type,cardinality\nattribute,Shop,Customer,email,VARCHAR(255),true,true,,,,\nattribute,Shop,Order,amount,\"NUMERIC(12,2)\",false,false,,,,\nrelation,Shop,,,,,,Order,Customer,belongs_to,many-to-one"}`),
	)
	recorder := httptest.NewRecorder()

	AnalyzeHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	body := recorder.Body.String()
	for _, fragment := range []string{
		`"database":{"name":"Shop"}`,
		"Structured CSV input was parsed directly",
		"CREATE TABLE customer",
		"FOREIGN KEY (customer_id) REFERENCES customer(id)",
	} {
		if !strings.Contains(body, fragment) {
			t.Fatalf("expected response body to contain %q, got: %s", fragment, body)
		}
	}
}

func TestGenerateSQLHandlerReturnsSQLForEditedModel(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPost,
		"/generate-sql",
		strings.NewReader(`{
			"database":{"name":"Shop"},
			"entities":[
				{"name":"Customer","attributes":[{"name":"email","type":"VARCHAR(255)","required":true}]},
				{"name":"Order","attributes":[{"name":"amount","type":"NUMERIC(12,2)","required":false}]}
			],
			"relations":[{"from":"Order","to":"Customer","type":"belongs_to","cardinality":"many-to-one"}]
		}`),
	)
	recorder := httptest.NewRecorder()

	GenerateSQLHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	body := recorder.Body.String()
	if !strings.Contains(body, "CREATE DATABASE shop") ||
		!strings.Contains(body, "CREATE TABLE order_table") ||
		!strings.Contains(body, "customer_id INTEGER") ||
		!strings.Contains(body, `"diagnostics"`) ||
		!strings.Contains(body, `"transformations"`) {
		t.Fatalf("expected generated SQL, diagnostics and transformations for edited model, got: %s", body)
	}
}
