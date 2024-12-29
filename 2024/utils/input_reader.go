package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Config struct {
	SessionToken string `json:"session_token"`
}

func ReadInput(day int) (string, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return "", fmt.Errorf("Error opening config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return "", fmt.Errorf("Error decoding config file: %w", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/2024/day/%d/input", day), nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: config.SessionToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to fetch input data: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %w", err)
	}

	return string(body), nil
}
