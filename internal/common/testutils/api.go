package testutils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func DoRequest(
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
