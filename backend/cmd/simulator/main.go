package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8000"
	}

	log.Printf("Starting Activity Simulator against %s", apiURL)

	// Register a dummy user
	email := fmt.Sprintf("sim_%d@example.com", time.Now().Unix())
	password := "simpassword"
	
	registerPayload, _ := json.Marshal(map[string]string{"email": email, "password": password})
	res, err := http.Post(apiURL+"/api/auth/register", "application/json", bytes.NewBuffer(registerPayload))
	if err != nil {
		log.Fatalf("Error registering: %v", err)
	}
	res.Body.Close()

	// Login to get JWT
	loginPayload, _ := json.Marshal(map[string]string{"email": email, "password": password})
	res, err = http.Post(apiURL+"/api/auth/login", "application/json", bytes.NewBuffer(loginPayload))
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}
	
	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&loginResp); err != nil {
		log.Fatalf("Error decoding login response: %v", err)
	}
	res.Body.Close()

	token := loginResp.Token
	log.Println("Successfully authenticated simulator user.")

	client := &http.Client{Timeout: 5 * time.Second}

	for {
		// Randomly decide to get tasks or create a task
		if rand.Float32() > 0.5 {
			req, _ := http.NewRequest(http.MethodGet, apiURL+"/api/tasks", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := client.Do(req)
			if err == nil {
				log.Println("GET /api/tasks ->", resp.Status)
				resp.Body.Close()
			}
		} else {
			taskPayload, _ := json.Marshal(map[string]string{"title": fmt.Sprintf("Random Task %d", rand.Intn(1000))})
			req, _ := http.NewRequest(http.MethodPost, apiURL+"/api/tasks", bytes.NewBuffer(taskPayload))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err == nil {
				log.Println("POST /api/tasks ->", resp.Status)
				resp.Body.Close()
			}
		}
		
		time.Sleep(2 * time.Second)
	}
}
