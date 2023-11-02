package utils

import (
	"encoding/json"
	"goAlertStore/models"
	"net/http"
)

// RespondWithError sends an error response in JSON format.
func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{
		"alert_id": "",
		"error":    message})
}

// RespondWithJSON sends a JSON response with the given data and status code.
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "JSON serialization error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func RespondWithJSONusingCustomResponse(w http.ResponseWriter, r *models.AlertDataNested, status int) {
	response, err := json.Marshal(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "JSON serialization error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
