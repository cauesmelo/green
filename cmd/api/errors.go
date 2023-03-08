package main

import (
	"net/http"

	"github.com/cauesmelo/green/internal/data"
)

func (app *application) logError(err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	data := data.Err{Error: message}
	err := app.writeJSON(w, status, data, nil)
	if err != nil {
		app.logError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(err)

	app.errorResponse(w, r, http.StatusInternalServerError, "Internal server error")
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "Not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
