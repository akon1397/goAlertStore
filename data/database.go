package data

import (
	"database/sql"
	"goAlertStore/models"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db     *sql.DB
	dbPath string
}

func NewDatabase(dbPath string) *Database {
	return &Database{db: nil, dbPath: "alert1.db"}
}

func (d *Database) Open() error {
	db, err := sql.Open("sqlite3", d.dbPath)
	if err != nil {
		return err
	}
	d.db = db

	// Create the "alerts" table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS alerts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alert_id TEXT NOT NULL,
		service_id TEXT,
		service_name TEXT,
		model TEXT,
		alert_type TEXT,
		alert_ts TEXT,
		severity TEXT,
		team_slack TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = d.db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d *Database) CreateAlert(alert *models.Alert) error {
	insertSQL := `
		INSERT INTO alerts (alert_id, service_id, service_name, model, alert_type, alert_ts, severity, team_slack)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := d.db.Exec(
		insertSQL,
		alert.AlertID,
		alert.ServiceID,
		alert.ServiceName,
		alert.Model,
		alert.AlertType,
		alert.AlertTS,
		alert.Severity,
		alert.TeamSlack,
	)
	return err
}

func (d *Database) GetAlertsByServiceAndTime(serviceID string, startTS, endTS int64) ([]models.Alert, error) {
	var alerts []models.Alert
	// Database query to retrieve alerts based on service_id and time range
	rows, err := d.db.Query("SELECT alert_id, service_id, service_name, model, alert_type, alert_ts, severity, team_slack FROM alerts WHERE service_id = ? AND alert_ts >= ? AND alert_ts <= ?", serviceID, startTS, endTS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alert models.Alert
		// Mapping the database columns to the alert struct
		if err := rows.Scan(&alert.AlertID, &alert.ServiceID, &alert.ServiceName, &alert.Model, &alert.AlertType, &alert.AlertTS, &alert.Severity, &alert.TeamSlack); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}
