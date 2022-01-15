package rest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/storage"
)

func TestNewHandler(t *testing.T) {
	_, err := NewHandler(nil)
	assert.NotNil(t, err)

	srv, err := NewHandler(&server.Server{})
	assert.Nil(t, err)
	assert.NotNil(t, srv)
}

func newTestHandler(t *testing.T, metricStorer storage.MetricStorer) *Handler {
	srv, err := server.NewServer(metricStorer)
	if err != nil {
		t.Fatal(err.Error())
	}

	h, err := NewHandler(srv)
	if err != nil {
		t.Fatal(err.Error())
	}

	return h
}

func doRequest(
	t *testing.T,
	server *httptest.Server,
	method, path string,
	data *[]byte,
) (int, string) {
	var body io.Reader
	if data != nil {
		body = bytes.NewBuffer(*data)
	}
	request, err := http.NewRequest(method, server.URL+path, body)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	return response.StatusCode, string(responseBody)
}
