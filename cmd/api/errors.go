package main

import (
	"net/http"

	"github.com/cauesmelo/green/internal/data"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, errors any) {
	data := data.Err{Error: errors}
	err := app.writeJSON(w, status, data, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

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

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusConflict, "Unable to update due to conflict")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded")
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusUnauthorized, "invalid credentials")
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	app.errorResponse(w, r, http.StatusUnauthorized, "invalid or missing authentication token")
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusUnauthorized, "you must be authenticated to access this resource")
}

func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusForbidden, "your user account must be activated to access this resource")
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusForbidden, "your user account doesn't have the necessary permissions to access this resource")
}
