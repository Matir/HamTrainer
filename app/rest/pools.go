package rest

import (
	"app/pools"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

func poolHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pool, err := pools.GetPool(vars["class"])
	if err != nil {
		http.Error(w, "Unable to find pool.", 404)
		return
	}
	json.NewEncoder(w).Encode(pool)
}
