package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseWriter interface{
	WriteError(error)
}

func registerError(err error, w http.ResponseWriter, resp responseWriter) {
	fmt.Printf("Can't get body: %v", err)
	resp.WriteError(err)
	respJson, _ := json.Marshal(resp) // ошибки быть не может
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, string(respJson))
}

func registerSuccess(w http.ResponseWriter, resp responseWriter) {
	respJson, _ := json.Marshal(resp)  // ошибки быть не может
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(respJson))
}
