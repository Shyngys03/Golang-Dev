package data

import (
	"database/sql"

	"electricaltools/internal/validator"

	"context"
	"errors"
	"time"
	"fmt"

	_ "github.com/lib/pq"
)

type Drill struct {
	ID int64			`json:"id"`
	Weight float32		`json:"weight"`
	Name string			`json:"name"`
	CableLength float32 `json:"length"`
	Worktime int32 		`json:"worktime"`
	ChuckDiameter int32 `json:"diameter"`
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


type DrillModel struct {
	DB *sql.DB
}

func (d DrillModel) Insert(drill *Drill) error {
	query := `
			INSERT INTO (name, work_time, chuck_diameter)
			VALUES ($1, $2, $3)
			RETURNING (id, weight, cable_length)`
	
	args := []interface{}{drill.Name, drill.Worktime, drill.ChuckDiameter}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&drill.ID, &drill.Weight, &drill.CableLength)
}


func (d DrillModel) Get(id int64) (*Drill, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
			SELECT id, weight, name, cable_length, work_time, chuck_diameter
			FROM drills
			WHERE id = $1`
	
	var drill Drill

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := d.DB.QueryRowContext(ctx, query, id).Scan(
		&drill.ID,
		&drill.Weight,
		&drill.Name,
		&drill.CableLength,
		&drill.Worktime,
		&drill.ChuckDiameter,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &drill, nil
}


func (d DrillModel) Update(drill *Drill) error {
	query := `
			UPDATE drills
			SET weight = $1, cable_length = $2, work_time = $3, chuck_diameter = $4
			WHERE id = $5 and name = $6
			RETURNING name`
	
	args := []interface{}{
		drill.Weight,
		drill.CableLength,
		drill.Worktime,
		drill.ChuckDiameter,
		drill.ID,
		drill.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := d.DB.QueryRowContext(ctx, query, args...).Scan(&drill.Name)
	if err != nil {
		switch{
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}


func (d DrillModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
			DELETE FROM drills
			WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := d.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (d DrillModel) GetAll(weight float32, name string, filters Filters) ([]*Drill, Metadata, error) {

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, weight, name, cable_length, work_time, chuck_diameter 
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (weight @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{weight, name, filters.limit(), filters.offset()}

	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	drills := []*Drill{}

	for rows.Next() {
		var drill Drill
		err := rows.Scan(
			
			&drill.ID,
			&drill.Weight,
			&drill.Name,
			&drill.CableLength,
			&drill.Worktime,
			&drill.ChuckDiameter,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		drills = append(drills, &drill)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return drills, metadata, nil
}