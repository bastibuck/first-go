package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3" // needed for the side-effects of this package (registering sqlite3 with database/sql)
)

func OpenConnection() *sql.DB {
	// Open a database connection
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	// Get the directory of the current file
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFilePath)

	eventsSql, err := os.ReadFile(filepath.Join(currentDir, "sql/events.sql"))
	if err != nil {
		log.Fatal(err)
	}

	usersSql, err := os.ReadFile(filepath.Join(currentDir, "sql/users.sql"))
	if err != nil {
		log.Fatal(err)
	}

	// Combine both SQL strings
	combinedSql := string(eventsSql) + "\n" + string(usersSql)

	_, err = db.Exec(combinedSql)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection established and tables prepared")

	return db
}
