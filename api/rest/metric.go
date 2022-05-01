package rest

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
)

func (h *Handler) updateMetricWithURL(w http.ResponseWriter, r *http.Request) {
	metricType := model.MetricType(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricStringValue := chi.URLParam(r, "metricValue")

	metric, err := model.MetricFromString(metricName, metricType, metricStringValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Server.PushMetric(r.Context(), metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
