package api

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, data envelope, statusCode int) error {

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	w.WriteHeader(statusCode)
	w.Write(js)

	return nil
}
