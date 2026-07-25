package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// ==============================================================================
// IMPORTANT: UPDATE THESE URLs FIRST BEFORE RUNNING!
// ==============================================================================
const (
	DATABASE_URL = "POSTGRES_CONNECTION_URL"
	REDIS_URL    = "REDIS_CONNECTION_URL"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=======================================================")
		fmt.Println("⚠️  Make sure you have updated the URLs in the code first! And also add user IP in trusted sources.")
		fmt.Println("=======================================================")
		fmt.Println("Which experiment would you like to run?")
		fmt.Println("1) Database Connection Exhaustion (PgBouncer vs Direct)")
		fmt.Println("2) Redis Memory Limit Test (LRU Eviction)")
		fmt.Println("3) Exit")
		fmt.Print("\nEnter choice (1, 2, or 3): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			runPostgresExperiment()
		case "2":
			runRedisExperiment()
		case "3":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter 1, 2, or 3.")
		}
	}
}

func runPostgresExperiment() {
	fmt.Println("\n--- Starting Database Connection Exhaustion Test ---")
	if DATABASE_URL == "" {
		fmt.Println("Error: DATABASE_URL is empty. Please set it at the top of the file.")
		return
	}

	var wg sync.WaitGroup
	conns := 1000 // Attempt to open 1000 connections

	fmt.Printf("Attempting to open %d connections to: %s\n", conns, DATABASE_URL)

	for i := 0; i < conns; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			db, err := sql.Open("postgres", DATABASE_URL)
			if err != nil {
				return
			}
			// Keep the connection open by executing a long sleep
			_, err = db.Exec("SELECT pg_sleep(30)")
			if err != nil {
				fmt.Printf("Conn %d failed: %v\n", id, err)
			} else {
				fmt.Printf("Conn %d succeeded\n", id)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("--- Database Test Finished ---")
}

func runRedisExperiment() {
	fmt.Println("\n--- Starting Redis Memory Limit Test ---")
	if REDIS_URL == "" {
		fmt.Println("Error: REDIS_URL is empty. Please set it at the top of the file.")
		return
	}

	opt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		fmt.Printf("Error parsing REDIS_URL: %v\n", err)
		return
	}
	client := redis.NewClient(opt)
	ctx := context.Background()

	// Generate 1MB of random data
	data := make([]byte, 1024*1024)
	rand.Read(data)

	fmt.Println("Starting to fill Redis with 1MB blocks...")
	for i := 0; i < 2000; i++ { // Attempt to write 2GB total
		key := fmt.Sprintf("junk_data_%d", i)
		err := client.Set(ctx, key, data, 0).Err()
		if err != nil {
			fmt.Printf("Failed at %d MB: %v\n", i, err)
			break
		}
		if i%100 == 0 {
			fmt.Printf("Inserted %d MB\n", i)
		}
	}
	fmt.Println("--- Redis Test Finished ---")
}
