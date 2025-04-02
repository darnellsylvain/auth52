package api

import (
	"encoding/json"
	"net/http"
)

func sendJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	d, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(d)
	return err
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {	
		return err
	}

	return nil
}