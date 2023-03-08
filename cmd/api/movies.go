package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cauesmelo/green/internal/data"
)

type movieCreate struct {
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Runtime int      `json:"runtime"`
	Genres  []string `json:"genres"`
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var mv movieCreate

	err := app.readJSON(w, r, &mv)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v", mv)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
