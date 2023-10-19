package main

import (
	"electricaltools/internal/data"
	"fmt"
	"net/http"
	"encoding/json"
)


func (app *application) createDrillHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ID int64			`json: "id"`
		Weight float32		`json: "weight"`
		Name string			`json: "name"`
		CableLength float32 `json: "length"`
		Worktime int32 		`json: "worktime"`
		ChuckDiameter int32 `json: "diameter"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	
	fmt.Fprintf(w, "%+v\n", input)
}


func (app *application) showDrillHandler(w http.ResponseWriter, r *http.Request) {
	
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	drill := data.Drill{
		ID: id,
		Weight: 2.5,
		Name: "Bosch",
		CableLength: 3.2,
		Worktime: 5,
		ChuckDiameter: 15,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"drill": drill}, nil)
	
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
