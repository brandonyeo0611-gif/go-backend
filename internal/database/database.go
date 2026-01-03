package database
import (
	"database/sql"
	_ "github.com/lib/pq"
)
type Database struct {
	Conn *sql.DB
}

func GetDB() (*Database, error) {
	connStr := "postgres://postgres.vdqfthvqkysiaixcqdkq:U1wpBBBD64VRDn6Y@aws-1-ap-northeast-1.pooler.supabase.com:6543/postgres"
	
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
