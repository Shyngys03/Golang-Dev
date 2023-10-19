package main

import (
	"fmt"
	"net/http"
	"electricaltools/internal/data"
	"electricaltools/internal/validator"
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
	
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	drill := &data.Drill{
		ID: input.ID,
		Weight: input.Weight,
		Name: input.Name,
		CableLength: input.CableLength,
		Worktime: input.Worktime,
		ChuckDiameter: input.ChuckDiameter,
	}

	v := validator.New()

	if data.ValidateDrill(v, drill); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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
