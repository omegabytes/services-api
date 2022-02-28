package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/omegabytes/services-api/models"
)

type Store struct {
	DB        *sql.DB
	Limit     uint16
	Precision float32
}

func (s *Store) GetService(requestedId string) ([]models.Service, error) {
	results := []models.Service{}
	// todo: sqlinjection guard
	// tested on curl localhost:8080/services/105%20or%201%3D1
	row := s.DB.QueryRow("SELECT * FROM servicetable WHERE id = $1;", requestedId)

	var id int
	var description sql.NullString
	var name string
	var versions interface{}

	switch err := row.Scan(&id, &name, &description, &versions); err {
	case nil:
		s := models.Service{
			Id:          id,
			Name:        name,
			Description: description.String,
		}

		if versions != nil {
			var v = []models.ServiceVersion{}
			json.Unmarshal([]byte(versions.([]uint8)), &v)
			s.Versions = v
		}
		results = append(results, s)
	default:
		return nil, err
	}
	return results, nil
}

func (s *Store) SearchServices(searchTerm string) ([]models.Service, error) {
	queryStmt := fmt.Sprintf(`SELECT * FROM servicetable WHERE SIMILARITY((name || ' ' || description), '%s') > %f limit %d;`, searchTerm, s.Precision, s.Limit)
	rows, err := s.DB.Query(queryStmt)
	if err != nil {
		return nil, fmt.Errorf("Invalid query")
	}
	defer rows.Close()

	results, err := scanResults(rows)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Store) ListServices(offset int) ([]models.Service, error) {
	queryStmt := fmt.Sprintf("SELECT * FROM servicetable LIMIT %d", s.Limit)

	if offset != 0 {
		queryStmt = fmt.Sprintf("%s OFFSET %d", queryStmt, offset)
	}

	rows, err := s.DB.Query(queryStmt)
	if err != nil {
		return nil, fmt.Errorf("Invalid query")
	}
	defer rows.Close()

	results, err := scanResults(rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func scanResults(rows *sql.Rows) ([]models.Service, error) {
	results := []models.Service{}
	for rows.Next() {
		var id int
		var description sql.NullString
		var name string
		var versions interface{}

		switch err := rows.Scan(&id, &name, &description, &versions); err {
		case nil:
			s := models.Service{
				Id:          id,
				Name:        name,
				Description: description.String,
			}

			if versions != nil {
				var v = []models.ServiceVersion{}
				json.Unmarshal([]byte(versions.([]uint8)), &v)
				s.Versions = v
			}

			results = append(results, s)
		default:
			return nil, err
		}
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}
	return results, nil
}
