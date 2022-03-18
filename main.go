package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const CF_BASE = "https://api.cloudflare.com/client/v4/"

func getPublicIP(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://v4.ident.me", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request to v4.ident.me: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code [%d]", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return strings.TrimSpace(string(content)), nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
	cfToken := os.Getenv("CF_API_TOKEN")
	domain := os.Getenv("DOMAIN")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Println("Getting public IP...")
	ip, err := getPublicIP(ctx)
	if err != nil {
		log.Fatalf("failed to get public IP: %v", err)
	}
	log.Printf("Found public IP: %s", ip)

	log.Printf("Updating %s A record to %s\n", domain, ip)
}
