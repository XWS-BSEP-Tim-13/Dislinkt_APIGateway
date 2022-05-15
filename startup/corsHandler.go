package startup

import "net/http"

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
	}
}
