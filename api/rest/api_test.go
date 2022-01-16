package rest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func newTestHandler(t *testing.T, metricStorage storage.MetricStorage) *Handler {
	srv, err := server.NewServer(metricStorage)
	require.NoError(t, err)

	h, err := NewHandler(srv)
	require.NoError(t, err)

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
	require.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	responseBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	defer response.Body.Close()

	return response.StatusCode, string(responseBody)
}
