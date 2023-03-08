package main

import (
	"net/http"
)

type status struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	s := status{
		Status:      "available",
		Environment: app.config.env,
		Version:     version,
	}

	err := app.writeJSON(w, http.StatusOK, s, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
