package data

import (
	"electricaltools/internal/validator"
)

type Drill struct {
	ID int64			`json: "id"`
	Weight float32		`json: "weight"`
	Name string			`json: "name"`
	CableLength float32 `json: "length"`
	Worktime int32 		`json: "worktime"`
	ChuckDiameter int32 `json: "diameter"`
}

func ValidateDrill(v *validator.Validator, drill *Drill) {

	v.Check(drill.ID > 0, "id", "must not be equal to 0")
	v.Check(len(drill.Name) <= 20, "name", "must not be more than 20 bytes long")
	v.Check(drill.Weight > 0, "weight", "must have a weight")
	v.Check(drill.CableLength <= 10.0, "length", "must be less than 10 meters")
	v.Check(drill.Worktime <= 20, "worktime", "must not be very long")
	v.Check(drill.ChuckDiameter > 5, "diameter", "must be greater than 5 mm")
	v.Check(drill.Worktime > 0, "worktime", "must work")
}