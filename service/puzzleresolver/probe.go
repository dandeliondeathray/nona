package puzzleresolver

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

// StartProbes starts readiness and liveness probes.
//
// Blocks until the server fails, therefore you must call it as a goroutine.
func StartProbes(port int) {
	http.HandleFunc("/readiness", readiness)
	http.HandleFunc("/liveness", liveness)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
