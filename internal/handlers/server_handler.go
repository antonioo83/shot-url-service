package handlers

import "net/http"

func UrlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, err := GetQuery("id", r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", "https://www.rbc.ru")
		w.WriteHeader(307)
		w.Write([]byte(id))
	case http.MethodPost:
		originalUrl, err := GetUrlParameter(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shotUrl, err := GetShortUrl(originalUrl, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shotUrl))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.WriteHeader(400)
		w.Write([]byte("error request"))
	}
}
