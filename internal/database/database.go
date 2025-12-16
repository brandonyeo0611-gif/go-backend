package database
import (
	"database/sql"
	_ "github.com/lib/pq"
)
type Database struct {
	Conn *sql.DB
}

func GetDB() (*Database, error) {
	connStr := "host=localhost port=5432 user=postgres password=Benbryan1 dbname=yap sslmode=disable"
	
	// accessing database
	db, err := sql.Open("postgres", connStr)

	// checking if got errors
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Conn: db (need remember)

	return &Database{Conn: db}, nil
}
