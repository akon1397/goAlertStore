package api

import "github.com/go-chi/chi"

func RegisterRoutes(r *chi.Mux, handler *AlertHandler) {
	r.Post("/alerts", handler.CreateAlert)
	r.Get("/alerts", handler.GetAlertsByServiceIdAndTime)
}
