package cgo

import (
	"encoding/json"
	"io"
	"net/http"
)

func ResultOk(w http.ResponseWriter, data string) {
	_, _ = io.WriteString(w, data)
}
func ResultFail(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusBadRequest)
}

func ResultJsonOk(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonString,_ := json.Marshal(data)
	_, _ = w.Write(jsonString)
}
