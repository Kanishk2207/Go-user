// database/database.go
package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {

    
    url := "mysql:abc@tcp(0.0.0.0:3306)/bugsmirror"


    // Connect to database
    db, err := sql.Open("mysql", url)
    if err != nil {
        return nil, err
    }

    // Ping database to verify connection
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}
