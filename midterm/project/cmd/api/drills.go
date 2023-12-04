package main

import (
	"electricaltools/internal/data"
	"electricaltools/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createDrillHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
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

	err = app.models.Drills.Insert(drill)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/drills/%d", drill.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"drills": drill}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showDrillHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	drill, err := app.models.Drills.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"drill": drill}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateDrillHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	drill, err := app.models.Drills.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Weight *float32		`json: "weight"`
		Name *string			`json: "name"`
		CableLength *float32 `json: "length"`
		Worktime *int32 		`json: "worktime"`
		ChuckDiameter *int32 `json: "diameter"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Weight != nil {
		drill.Weight = *input.Weight
	}
	if input.Name != nil {
		drill.Name = *input.Name
	}
	if input.CableLength != nil {
		drill.CableLength = *input.CableLength
	}
	if input.Worktime != nil {
		drill.Worktime = *input.Worktime
	}
	if input.ChuckDiameter != nil {
		drill.ChuckDiameter = *input.ChuckDiameter
	}
	

	v := validator.New()

	if data.ValidateDrill(v, drill); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Drills.Update(drill)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"drill": drill}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteDrillHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.models.Drills.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "drill successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDrillsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Weight float32
		Name string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Weight = app.readFloat(qs, "weight", 1.0, v)
	input.Name = app.readString(qs, "name", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "weight", "name", "-id", "-weight", "-name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	drills, metadata, err := app.models.Drills.GetAll(input.Weight, input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"drills": drills, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}