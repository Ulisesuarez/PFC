package util

import (
	"encoding/json"
	"net/http"
	"strconv"
)

//JSONResponse ...
func JSONResponse(response interface{}, w http.ResponseWriter, statusCode int) {

	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if statusCode <= 0 {
		statusCode = http.StatusOK
	}

	if statusCode == http.StatusNoContent {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(statusCode)
	w.Write(data)
}
