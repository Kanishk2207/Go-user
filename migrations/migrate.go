// migrations/migrate.go
package migrations

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const migrationsDir = "./migrations"

func Run(db *sql.DB) error {
    err := filepath.Walk(migrationsDir, func(path string, info fs.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() || !strings.HasSuffix(info.Name(), ".sql") {
            return nil
        }
        script, err := os.ReadFile(path)
        if err != nil {
            return err
        }
        _, err = db.Exec(string(script))
        if err != nil {
            return fmt.Errorf("failed to execute migration %s: %v", info.Name(), err)
        }
        log.Printf("Executed migration: %s\n", info.Name())
        return nil
    })
    if err != nil {
        return err
    }

    log.Println("Database migration completed successfully.")
    return nil
}
