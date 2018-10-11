package handlers

import (
	"net/http"
)

func Github(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
