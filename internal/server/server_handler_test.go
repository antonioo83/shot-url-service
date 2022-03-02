package server

import (
	"github.com/antonioo83/shot-url-service/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRouters(t *testing.T) {
	type want struct {
		httpStatus int
	}
	tests := []struct {
		name        string
		code        string
		originalURL string
		wantPost    want
		wantGet     want
	}{
		{
			name:        "positive test #1",
			code:        "H1P4S7fw",
			originalURL: "https://stackoverflow.com/questions/15240884/how-can-i-handle-http-requests-of-different-methods-to-in-go",
			wantPost: want{
				httpStatus: 201,
			},
			wantGet: want{
				httpStatus: 307,
			},
		},
	}

	r := GetRouters()
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		jsonData := strings.NewReader(string(handlers.GetJsonRequest("url", tt.originalURL)))
		jsonRequest := getJsonPostRequest(t, ts, "/api/shorten", jsonData)
		resp, jsonResponse := sendRequest(t, jsonRequest)
		resultParameter, err := handlers.GetResultParameter(jsonResponse)
		require.NoError(t, err)
		assert.Equal(t, tt.wantPost.httpStatus, resp.StatusCode)
		assert.Equal(t, ts.URL+"/"+tt.code, resultParameter)
		resp.Body.Close()

		postRequest := getPostRequest(t, ts, "/", strings.NewReader(tt.originalURL))
		resp, body := sendRequest(t, postRequest)
		assert.Equal(t, tt.wantPost.httpStatus, resp.StatusCode)
		assert.Equal(t, ts.URL+"/"+tt.code, body)
		resp.Body.Close()

		getGetRequest := getGetRequest(t, ts, "/"+tt.code, nil)
		resp, body = sendRequest(t, getGetRequest)
		assert.Equal(t, tt.wantGet.httpStatus, resp.StatusCode)
		assert.Equal(t, tt.originalURL, body)
		resp.Body.Close()
	}
}

func sendRequest(t *testing.T, req *http.Request) (*http.Response, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func getPostRequest(t *testing.T, ts *httptest.Server, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", ts.URL+path, body)
	req.Header.Add("Content-Type", "text/plain")
	require.NoError(t, err)

	return req
}

func getGetRequest(t *testing.T, ts *httptest.Server, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", ts.URL+path, body)
	require.NoError(t, err)

	return req
}

func getJsonPostRequest(t *testing.T, ts *httptest.Server, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", ts.URL+path, body)
	req.Header.Add("Content-Type", "application/json")
	require.NoError(t, err)

	return req
}
