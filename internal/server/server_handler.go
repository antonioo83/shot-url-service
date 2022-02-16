package server

import (
	"github.com/antonioo83/shot-url-service/internal/handlers"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/localcache"
	"net/http"
	"strings"
)

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		code := strings.Replace(r.RequestURI, "/", "", 1)
		model, err := localcache.FindByCode(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", model.OriginalUrl)
		w.WriteHeader(307)
		//w.Write([]byte(model.OriginalUrl))
	case http.MethodPost:
		originalUrl, err := handlers.GetUrlParameter(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shotUrl, code, err := handlers.GetShortUrl(originalUrl, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var shortUrl models.ShortUrl
		shortUrl.Code = code
		shortUrl.OriginalUrl = originalUrl
		shortUrl.ShortUrl = shotUrl
		localcache.SaveUrl(shortUrl)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shotUrl))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.WriteHeader(400)
		w.Write([]byte("error request"))
	}
}
