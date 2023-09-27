package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// ReadReqBody - {out} have to be *ptr
func ReadReqBody(r *http.Request, out any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}
	return nil
}

func WriteOKRes(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
