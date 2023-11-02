package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"goAlertStore/data"
	"goAlertStore/models"
	"goAlertStore/utils"
)

type AlertHandler struct {
	db *data.Database
}

func NewAlertHandler(db *data.Database) *AlertHandler {
	return &AlertHandler{db}
}

func (h *AlertHandler) CreateAlert(w http.ResponseWriter, r *http.Request) {
	var alert models.Alert

	// Decoding the request body into the alert struct
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if alert.AlertID == "" || alert.ServiceID == "" || alert.AlertType == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Incomplete request payload")
		return
	}

	// Save the alert data to the chosen data storage (e.g., a database)
	if err := h.db.CreateAlert(&alert); err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{
			"alert_id": alert.AlertID,
			"error":    err.Error(),
		})
		return
	}

	// Return the response with the alert_id from the request and an empty error
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"alert_id": alert.AlertID,
		"error":    "",
	})
}

func (h *AlertHandler) GetAlertsByServiceIdAndTime(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the service ID from the URL path
	serviceID := r.URL.Query().Get("service_id")
	if serviceID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Service ID is required")
		return
	}

	// Parse and validate the start and end timestamps from query parameters
	startTS, err := strconv.ParseInt(r.URL.Query().Get("start_ts"), 10, 64)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid start_ts")
		return
	}
	endTS, err := strconv.ParseInt(r.URL.Query().Get("end_ts"), 10, 64)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid end_ts")
		return
	}

	if endTS < startTS {
		utils.RespondWithError(w, http.StatusBadRequest, "start_ts cannot be greater than end_ts")
		return
	}

	// Retrieving alerts from the database based on service_id and time range
	alerts, err := h.db.GetAlertsByServiceAndTime(serviceID, startTS, endTS)
	if len(alerts) == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "No alerts found")
		return
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving alerts")
		return
	}
	var alertsData []models.AlertData

	for _, alert := range alerts {
		alertDataForMapping := models.AlertData{
			AlertID:   alert.AlertID,
			Model:     alert.Model,
			AlertType: alert.AlertType,
			AlertTS:   alert.AlertTS,
			Severity:  alert.Severity,
			TeamSlack: alert.TeamSlack,
		}

		alertsData = append(alertsData, alertDataForMapping)
	}

	response := models.AlertDataNested{
		ServiceID:   serviceID,
		ServiceName: alerts[0].ServiceName, // Assumed service names are same for service_IDs with equal values, hence used service_name from one of the filtered alerts
		Alerts:      alertsData,
	}
	utils.RespondWithJSONusingCustomResponse(w, &response, http.StatusOK)
}
