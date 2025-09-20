package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"task-manager/backend/internal/repositories"

	"gorm.io/gorm"
)

func main() {
	// Connect to database
	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get command line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run migrate.go [up|down]")
	}

	command := os.Args[1]

	switch command {
	case "up":
		err = migrateUp(db)
	case "down":
		err = migrateDown(db)
	default:
		log.Fatal("Invalid command. Use 'up' or 'down'")
	}

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Migration completed successfully!")
}

func connectDB() (*gorm.DB, error) {
	dbCfg := repositories.NewDatabaseConfig()
	return dbCfg.Connect()
}

func migrateUp(db *gorm.DB) error {
	migrationDir := "migrations/up"
	files, err := getMigrationFiles(migrationDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Printf("Running migration: %s\n", file)
		
		content, err := ioutil.ReadFile(filepath.Join(migrationDir, file))
		if err != nil {
			return err
		}

		err = db.Exec(string(content)).Error
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	return nil
}

func migrateDown(db *gorm.DB) error {
	migrationDir := "migrations/down"
	files, err := getMigrationFiles(migrationDir)
	if err != nil {
		return err
	}

	// Run down migrations in reverse order
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		fmt.Printf("Rolling back migration: %s\n", file)
		
		content, err := ioutil.ReadFile(filepath.Join(migrationDir, file))
		if err != nil {
			return err
		}

		err = db.Exec(string(content)).Error
		if err != nil {
			return fmt.Errorf("failed to rollback migration %s: %w", file, err)
		}
	}

	return nil
}

func getMigrationFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	// Sort files to ensure correct order
	sort.Strings(sqlFiles)
	return sqlFiles, nil
}
