package handlers

import (
	"net/http"
	"strings"
)

var database = make(map[string]string)

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		urlId := strings.Replace(r.RequestURI, "/", "", 1)
		redirectUrl, ok := database[urlId]
		if ok == false {
			http.Error(w, string("Original Url not exist"), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", redirectUrl)
		w.WriteHeader(307)
		//w.Write([]byte())
	case http.MethodPost:
		originalUrl, err := GetUrlParameter(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shotUrl, urlId, err := GetShortUrl(originalUrl, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		database[urlId] = originalUrl
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shotUrl))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.WriteHeader(400)
		w.Write([]byte("error request"))
	}
}
