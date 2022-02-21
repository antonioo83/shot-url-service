package server

import (
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
		originalUrl string
		wantPost    want
		wantGet     want
	}{
		{
			name:        "positive test #1",
			code:        "H1P4S7fw",
			originalUrl: "https://stackoverflow.com/questions/15240884/how-can-i-handle-http-requests-of-different-methods-to-in-go",
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
		resp, body := sendRequest(t, ts, "POST", "/", strings.NewReader(tt.originalUrl))
		assert.Equal(t, tt.wantPost.httpStatus, resp.StatusCode)
		assert.Equal(t, ts.URL+"/"+tt.code, body)
		resp.Body.Close()

		resp, body = sendRequest(t, ts, "GET", "/"+tt.code, nil)
		assert.Equal(t, tt.wantGet.httpStatus, resp.StatusCode)
		assert.Equal(t, tt.originalUrl, body)
		resp.Body.Close()
	}
}

func sendRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
