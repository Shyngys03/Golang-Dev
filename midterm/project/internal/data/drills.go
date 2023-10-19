package data

type Drill struct {
	ID int64			`json: "id"`
	Weight float32		`json: "weight"`
	Name string			`json: "name"`
	CableLength float32 `json: "length"`
	Worktime int32 		`json: "worktime"`
	ChuckDiameter int32 `json: "diameter"`
}
