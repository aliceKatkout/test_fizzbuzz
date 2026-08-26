package api

import (
	"net/http"

	"github.com/alicekatkout/fizzbuzz-api/internal/stats"
)

// NewRouter builds the complete HTTP handler for the API, including
// logging and panic-recovery middleware.
func NewRouter(store *stats.Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /fizzbuzz", FizzBuzzHandler(store))
	mux.HandleFunc("GET /statistics", StatisticsHandler(store))
	mux.HandleFunc("GET /health", healthHandler)

	return withRecovery(withLogging(mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
