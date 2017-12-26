package control

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//go:generate mockgen -destination=../mock/mock_control.go -package=mock github.com/dandeliondeathray/nona/control Target

// Target is the interface between the control layer and the rest of the game.
type Target interface {
	NewRound(seed int64)
}

// NewRouter creates a router for all supported HTTP endpoints in the control layer
func NewRouter(target Target) *mux.Router {
	r := mux.NewRouter()
	c := controlHandler{target}
	r.HandleFunc("/round/{seed:[0-9]+}", c.handleNewRound)
	return r
}

type controlHandler struct {
	target Target
}

func (c *controlHandler) handleNewRound(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	seedString := vars["seed"]

	seed, err := strconv.ParseInt(seedString, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.target.NewRound(seed)
	w.WriteHeader(http.StatusOK)
}

// StartControl listens to HTTP requests to the control layer.
func StartControl(target Target) {
	r := NewRouter(target)
	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}
