package server

import (
	"context"
	"fmt"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers"
	authFactory "github.com/antonioo83/shot-url-service/internal/handlers/auth/factory"
	"github.com/antonioo83/shot-url-service/internal/handlers/generators"
	"github.com/antonioo83/shot-url-service/internal/repositories/factory"
	"github.com/antonioo83/shot-url-service/internal/utils"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v4/pgxpool"
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

	var tokenAuth *jwtauth.JWTAuth
	var pool *pgxpool.Pool
	context := context.Background()
	configFromFile, err := config.LoadConfigFile("config.json")
	if err != nil {
		fmt.Println("i can't load configuration file:" + err.Error())
	}
	config := config.GetConfigSettings(configFromFile)
	if config.IsUseDatabase {
		pool, _ := pgxpool.Connect(context, config.DatabaseDsn) //databaseRepository.Connect(context)
		defer pool.Close()
	}

	userRepository := factory.GetUserRepository(context, pool, config)
	r := GetRouters(RouteParameters{
		Config:             config,
		ShotURLRepository:  factory.GetRepository(context, pool, config),
		UserRepository:     userRepository,
		DatabaseRepository: factory.GetDatabaseRepository(config),
		UserAuthHandler:    authFactory.NewAuthHandler(tokenAuth, userRepository, config),
		Generator:          generators.NewShortLinkDefaultGenerator(),
	})
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		request, err := handlers.GetJSONRequest("url", tt.originalURL)
		assert.NoError(t, err)

		jsonRequest := getJSONPostRequest(t, ts, "/api/shorten", strings.NewReader(string(request)))
		resp, jsonResponse := sendRequest(t, jsonRequest)
		resultParameter, err := handlers.GetResultParameter(jsonResponse)
		require.NoError(t, err)
		assert.Equal(t, tt.wantPost.httpStatus, resp.StatusCode)
		assert.Equal(t, ts.URL+"/"+tt.code, resultParameter)
		if err := resp.Body.Close(); err != nil {
			assert.NoError(t, err)
		}

		postRequest := getPostRequest(t, ts, "/", strings.NewReader(tt.originalURL))
		resp, body := sendRequest(t, postRequest)
		assert.Equal(t, tt.wantPost.httpStatus, resp.StatusCode)
		assert.Equal(t, ts.URL+"/"+tt.code, body)
		if err := resp.Body.Close(); err != nil {
			assert.NoError(t, err)
		}

		getGetRequest := getGetRequest(t, ts, "/"+tt.code, nil)
		resp, body = sendRequest(t, getGetRequest)
		assert.Equal(t, tt.wantGet.httpStatus, resp.StatusCode)
		assert.Equal(t, tt.originalURL, body)
		if err := resp.Body.Close(); err != nil {
			assert.NoError(t, err)
		}
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

	defer utils.ResourceClose(resp.Body)

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

func getJSONPostRequest(t *testing.T, ts *httptest.Server, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", ts.URL+path, body)
	req.Header.Add("Content-Type", "application/json")
	require.NoError(t, err)

	return req
}
