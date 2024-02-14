package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	workers   = 10    // Number of workers for parallel processing
	batchSize = 1000  // Batch size for inserts
	batchlog  = 10000 // Batch log
)

func main() {
	// Read SQL commands from the file
	// Add pathfile
	data, err := os.ReadFile("pathfile")
	if err != nil {
		panic(err)
	}

	commands := strings.Split(string(data), ";")

	// Register TLS config
	mysql.RegisterTLSConfig("custom", &tls.Config{
		InsecureSkipVerify: true,
	})

	// Open database connection
	db, err := sql.Open("mysql", "username:password@tcp(url)/dbname?tls=custom")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connected to the database")

	// Use a WaitGroup for parallel processing
	var wg sync.WaitGroup

	// Create a channel to communicate with workers
	commandsChan := make(chan string, len(commands))

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processCommands(db, commandsChan)
		}()
	}

	// Send commands to workers
	for _, cmd := range commands {
		if strings.TrimSpace(cmd) != "" {
			commandsChan <- cmd
		}
	}
	close(commandsChan)

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("All commands executed successfully")
}

// processCommands processes SQL commands from the channel
func processCommands(db *sql.DB, commandsChan <-chan string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	batchCount := 0
	for cmd := range commandsChan {
		_, err := tx.Exec(cmd)
		if err != nil {
			log.Printf("Failed to execute command: %s\n", cmd)
			log.Fatal(err)
		}

		batchCount++
		if batchCount%batchSize == 0 {
			if err := tx.Commit(); err != nil {
				log.Fatal(err)
			}
			tx, err = db.Begin()
			if err != nil {
				log.Fatal(err)
			}
		}

		if batchCount%batchlog == 0 {
			log.Printf("Executing command %d\n", batchCount)
		}
	}

	// Commit the remaining transactions
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
