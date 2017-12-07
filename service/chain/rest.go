package chain

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func readiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Service) getPuzzle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team := vars["team"]
	indexArgument := vars["index"]
	index, err := strconv.Atoi(indexArgument)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	puzzle, err := s.teams.GetPuzzle(team, index)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(puzzle))
}

// ListenAndServe starts the HTTP service.
//
// Blocks until the server fails, therefore you must call it as a goroutine.
func (s *Service) ListenAndServe(port int) {
	log.Printf("Starting REST service on port %d", port)
	r := mux.NewRouter()
	r.HandleFunc("/readiness", readiness)
	r.HandleFunc("/liveness", liveness)
	r.HandleFunc("/puzzle/{team}/{index:[0-9]+}", s.getPuzzle)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
