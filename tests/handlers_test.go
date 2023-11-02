package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"goAlertStore/api"
	"goAlertStore/data"
	"goAlertStore/models"
)

func TestCreateAlert_ValidAlert(t *testing.T) {
	db := data.NewDatabase("test.db")
	err := db.Open()
	if err != nil {
		t.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close()

	alertHandler := api.NewAlertHandler(db)

	alert := models.Alert{
		AlertID:     "b950482e9911ec7e41f7ca5e5d9a424f",
		ServiceID:   "my_test_service_id",
		ServiceName: "my_test_service",
		Model:       "my_test_model",
		AlertType:   "anomaly",
		AlertTS:     "1695644160",
		Severity:    "warning",
		TeamSlack:   "slack_ch",
	}

	alertJSON, _ := json.Marshal(alert)

	req, _ := http.NewRequest("POST", "/alerts", bytes.NewReader(alertJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	alertHandler.CreateAlert(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, but got %d", rr.Code)
	}

}

func TestGetAlertsByServiceIdAndTime_ValidParams(t *testing.T) {
	// Create a test database
	db := data.NewDatabase("test.db")
	err := db.Open()
	if err != nil {
		t.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close()

	alertHandler := api.NewAlertHandler(db)

	req, _ := http.NewRequest("GET", "/alerts?service_id=my_test_service_id&start_ts=1695644160&end_ts=1695644162", nil)
	rr := httptest.NewRecorder()

	alertHandler.GetAlertsByServiceIdAndTime(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", rr.Code)
	}
}

func TestCreateAlert_InvalidJSON(t *testing.T) {
	db := data.NewDatabase("test.db")
	err := db.Open()
	if err != nil {
		t.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close()

	alertHandler := api.NewAlertHandler(db)

	invalidJSON := []byte(`{"alert_id": "invalid"}`)

	req, _ := http.NewRequest("POST", "/alerts", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	alertHandler.CreateAlert(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, but got %d", rr.Code)
	}
}

func TestGetAlertsByServiceIdAndTime_InvalidParams(t *testing.T) {
	db := data.NewDatabase("test.db")
	err := db.Open()
	if err != nil {
		t.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close()

	alertHandler := api.NewAlertHandler(db)

	req, _ := http.NewRequest("GET", "/alerts?service_id=&start_ts=invalid&end_ts=invalid", nil)
	rr := httptest.NewRecorder()

	alertHandler.GetAlertsByServiceIdAndTime(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, but got %d", rr.Code)
	}
}
