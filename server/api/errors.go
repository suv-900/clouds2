package api

import "net/http"

func (app *application) sendResponse(w http.ResponseWriter, statusCode int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, env, statusCode)
	if err != nil {
		//err while sending
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request) {
	msg := "server error."
	app.sendResponse(w, http.StatusInternalServerError, msg)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
