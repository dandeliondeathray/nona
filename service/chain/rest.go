package chain

import (
	"fmt"
	"log"
	"net/http"
)

func readiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ListenAndServe starts the HTTP service.
//
// Blocks until the server fails, therefore you must call it as a goroutine.
func (s *Service) ListenAndServe(port int) {
	http.HandleFunc("/readiness", readiness)
	http.HandleFunc("/liveness", liveness)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
