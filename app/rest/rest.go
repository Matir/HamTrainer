// RESTful interface for HamTrainer
package rest

import (
	"fmt"
	"app/logging"
	"net/http"
)

func init() {
	http.HandleFunc("/api/ok", ok)
}

// Sanity test handler
func ok(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
	logging.Infof(r, "OK requested.")
}
