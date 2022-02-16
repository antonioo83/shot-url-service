package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUrlHandler(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code        int
		response    string
		contentType string
	}
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name        string
		urlId       string
		originalUrl string
		wantPost    want
		wantGet     want
	}{
		// определяем все тесты
		{
			name:        "positive test #1",
			urlId:       "H1P4S7fw",
			originalUrl: "https://stackoverflow.com/questions/15240884/how-can-i-handle-http-requests-of-different-methods-to-in-go",
			wantPost: want{
				code:        201,
				response:    "http://example.com/H1P4S7fw",
				contentType: "",
			},
			wantGet: want{
				code:        307,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			postRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.originalUrl))
			postRequest.Header.Add("Content-Type", "text/plain")

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(UrlHandler)
			// запускаем сервер
			h.ServeHTTP(w, postRequest)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.wantPost.code {
				t.Errorf("Expected status code %d, got %d", tt.wantPost.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.wantPost.response {
				t.Errorf("Expected body %s, got %s", tt.wantPost.response, w.Body.String())
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.wantPost.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.wantPost.contentType, res.Header.Get("Content-Type"))
			}

			getRequest := httptest.NewRequest(http.MethodGet, "/"+tt.urlId, nil)
			// создаём новый Recorder
			w2 := httptest.NewRecorder()
			// определяем хендлер
			h2 := http.HandlerFunc(UrlHandler)
			// запускаем сервер
			h2.ServeHTTP(w2, getRequest)
			res2 := w2.Result()
			// проверяем код ответа
			if res2.StatusCode != tt.wantGet.code {
				t.Errorf("Expected status code %d, got %d", tt.wantPost.code, w.Code)
			}

			// заголовок ответа
			if res2.Header.Get("Location") != tt.originalUrl {
				t.Errorf("Expected Content-Type %s, got %s", tt.wantPost.contentType, res.Header.Get("Location"))
			}
		})
	}
}
