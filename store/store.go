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

// GetService fetches a single service by a given service ID.
func (s *Store) GetService(requestedId string) ([]models.Service, error) {
	// I could have used sql.QueryRow (which returns at most one row) but elected to use QueryRows so I could reuse my scan fx.
	row, err := s.DB.Query("SELECT * FROM servicetable WHERE id = $1;", requestedId)
	if err != nil {
		return nil, fmt.Errorf("Invalid query")
	}
	defer row.Close()

	results, err := scanResults(row)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// SearchServices fetches all results from the services table that are "close enough" to a user's search term.
// The term is compared against a contatenation of the name and description. This will break in the event name is allowed to be nil in the future!
// The precision is a runtime config to support lower-lift product tweaking.
func (s *Store) SearchServices(searchTerm string, sort string) ([]models.Service, error) {
	queryStmt := fmt.Sprintf(`SELECT * FROM servicetable WHERE SIMILARITY((name || ' ' || description), '%s') > %f ORDER BY id %s limit %d;`, searchTerm, s.Precision, sort, s.Limit)
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

// ListServices fetches all services from the database up to the preconfigured limit.
// Note: Given time, I would add more advanced support in the query for fetching results after a certain date etc.
func (s *Store) ListServices(offset int, sort string) ([]models.Service, error) {
	queryStmt := fmt.Sprintf("SELECT * FROM servicetable ORDER BY id %s LIMIT %d ", sort, s.Limit)

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

// scanResults traverses/transforms multirow database response into service objects
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
