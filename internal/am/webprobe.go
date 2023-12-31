package am

import (
	"fmt"
	"net/http"
)

type (
	Probe struct{}
)

var (
	HealthzProbe = Probe{}
)

func (p Probe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, "OK")
	if err != nil {
		// WIP: Implement logging through worker
		return
	}
}
