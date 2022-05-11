package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	mw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/middleware"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/internal/model"
)

type Handler struct {
	Server *Server
	Router *chi.Mux
}

func NewHandler(server *Server) (*Handler, error) {
	if server == nil {
		return nil, errors.New("invalid server value: nil")
	}

	router := chi.NewRouter()

	h := &Handler{
		Server: server,
		Router: router,
	}

	logger := httplog.NewLogger("http-request-logger", httplog.Options{
		JSON: true,
	})

	h.Router.Use(mw.Recoverer)
	h.Router.Use(httplog.RequestLogger(logger))
	h.Router.Use(middleware.GzipDecoder())
	h.Router.Use(middleware.GzipEncoder())

	h.Router.Route("/update/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
		r.Post("/", h.updateMetricWithURL)
	})

	h.Router.Post("/update/", h.updateMetricWithBody)

	h.Router.Post("/updates/", h.updateMetricListWithBody)

	h.Router.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
		r.Get("/", h.getMetricWithURL)
	})

	h.Router.Post("/value/", h.getMetricWithBody)

	h.Router.Get("/", h.getMetricList)

	h.Router.Get("/ping", h.heartbeat)

	return h, nil
}

func (h *Handler) updateMetricWithURL(w http.ResponseWriter, r *http.Request) {
	metricType := model.MetricType(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricStringValue := chi.URLParam(r, "metricValue")

	if err := metricType.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	metric, err := model.MetricFromString(metricName, metricType, metricStringValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Server.PushMetric(r.Context(), metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updateMetricWithBody(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := metric.MType.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	if err := metric.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	valid, err := h.Server.ValidateHash(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(w, "invalid hash value", http.StatusBadRequest)
		return
	}

	if err := h.Server.PushMetric(r.Context(), metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(struct{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

func (h *Handler) updateMetricListWithBody(w http.ResponseWriter, r *http.Request) {
	var metrics []model.Metric
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, metric := range metrics {
		if err := metric.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		valid, err := h.Server.ValidateHash(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "invalid hash value", http.StatusBadRequest)
			return
		}
	}

	if err := h.Server.PushMetricList(r.Context(), metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(struct{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

func (h *Handler) getMetricWithURL(w http.ResponseWriter, r *http.Request) {
	metricType := model.MetricType(chi.URLParam(r, "metricType"))
	metricName := model.MetricName(chi.URLParam(r, "metricName"))

	metric := model.Metric{
		ID:    metricName,
		MType: metricType,
	}

	if err := metric.MType.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	if err := metric.ID.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := h.Server.LoadMetric(r.Context(), metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if m == nil {
		http.Error(w, fmt.Sprintf("Metric %s not found", metric.ID), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, m.String())
}

func (h *Handler) getMetricWithBody(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := metric.MType.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	if err := metric.ID.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := h.Server.LoadMetric(r.Context(), metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if m == nil {
		http.Error(w, fmt.Sprintf("Metric %s not found", metric.ID), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

func (h *Handler) getMetricList(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Metrics}}<div>{{ .ID }}: {{ .String }}</div>{{end}}
	</body>
	</html>`

	t, err := template.New("getMetricList").Parse(htmlTemplate)
	if err != nil {
		errCode := http.StatusInternalServerError
		http.Error(w, err.Error(), errCode)
		return
	}

	metrics, err := h.Server.LoadMetricList(r.Context())
	if err != nil {
		errCode := http.StatusInternalServerError
		http.Error(w, err.Error(), errCode)
		return
	}

	data := struct {
		Title   string
		Metrics []model.Metric
	}{
		Title:   "Metric List",
		Metrics: metrics,
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_ = t.Execute(w, data)
}

func (h *Handler) heartbeat(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
	defer cancel()

	if err := h.Server.Heartbeat(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
