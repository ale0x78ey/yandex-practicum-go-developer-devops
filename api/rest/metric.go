package rest

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

func (api *API) InitMetric() {
	api.Routes.Metric.Route("/update/{metricType}/{metricName}/{metricValue}",
		func(r chi.Router) {
			r.Use(withMetricTypeValidator)
			r.Post("/", updateMetric)
		})

	api.Routes.Metric.Route("/value/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(withMetricTypeValidator)
			r.Get("/", getMetric)
		})

	api.Routes.Metric.Get("/", getMetricList)
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	srv := ctx.Value(ContextServerKey).(*server.Server)
	metricType := model.MetricType(chi.URLParam(r, "metricType"))
	metricName := model.MetricName(chi.URLParam(r, "metricName"))

	metric := &model.Metric{
		Type:        metricType,
		Name:        metricName,
		StringValue: chi.URLParam(r, "metricValue"),
	}

	if err := srv.SaveMetric(ctx, metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do I need close handles in this cases?
	r.Body.Close()

	w.WriteHeader(http.StatusOK)
}

func getMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	srv := ctx.Value(ContextServerKey).(*server.Server)
	metricType := model.MetricType(chi.URLParam(r, "metricType"))
	metricName := model.MetricName(chi.URLParam(r, "metricName"))

	value, err := srv.LoadMetric(ctx, metricType, metricName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	r.Body.Close()

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, value)
}

func getMetricList(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Metrics}}<div>{{ .Name }}: {{ .StringValue }}</div>{{end}}
	</body>
	</html>`

	t, err := template.New("getMetricList").Parse(htmlTemplate)
	if err != nil {
		errCode := http.StatusInternalServerError
		http.Error(w, err.Error(), errCode)
		return
	}

	ctx := r.Context()
	srv := ctx.Value(ContextServerKey).(*server.Server)

	metrics, err := srv.LoadMetricList(ctx)
	if err != nil {
		errCode := http.StatusInternalServerError
		http.Error(w, err.Error(), errCode)
		return
	}

	data := struct {
		Title   string
		Metrics []*model.Metric
	}{
		Title:   "Metric List",
		Metrics: metrics,
	}

	r.Body.Close()

	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	_ = t.Execute(w, data)
}
