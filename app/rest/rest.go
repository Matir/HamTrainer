// RESTful interface for HamTrainer
package rest

import (
	"fmt"
	"app/logging"
	"net/http"
	"github.com/gorilla/mux"
)

func init() {
	SetCSRFKey("TODO_load_key_from_models")

	// Setup handlers
	r := mux.NewRouter()
	r.HandleFunc("/api/ok", ok)
	r.HandleFunc("/api/pool/{class}", poolHandler)
	http.Handle("/", r)
}

// Sanity test handler
func ok(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
	logging.Infof(r, "OK requested.")
}
