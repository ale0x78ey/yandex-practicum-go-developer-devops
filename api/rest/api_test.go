package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/service/server"
)

func TestInit(t *testing.T) {
	assert.Nil(t, Init(nil))
	assert.NotNil(t, Init(&server.Server{}))
}

func doRequest(
	t *testing.T,
	server *httptest.Server,
	method, path string,
) int {
	request, err := http.NewRequest(method, server.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	// TODO: statusCode, bodyContent
	return response.StatusCode
}
