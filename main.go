package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rmbreak/cfdyndns/internal/cloudflare"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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

func updateDomainRecord(ctx context.Context, cfToken string, domain string, ip string) error {
	cloudflareClient := cloudflare.New(cfToken)
	resp, err := cloudflareClient.UpdateDnsRecord(ctx, "", "", cloudflare.DnsUpdateRequestData{
		Type:    "A",
		Name:    domain,
		Content: ip,
		Ttl:     600,
	})
	if err != nil {
		return fmt.Errorf("failed to update dns record: %v", err)
	}
	if !resp.Success {
		return fmt.Errorf("received a failure response: %v", resp.Errors)
	}

	b, _ := json.Marshal(*resp)
	log.Debug().Msgf("%s", string(b))

	return nil
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	if err := godotenv.Load(); err != nil {
		log.Fatal().Msgf("%v", err)
	}

	logLevel := zerolog.InfoLevel
	if os.Getenv("LOG_LEVEL") != "" {
		level, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err != nil {
			log.Warn().Msgf("Unable to parse LOG_LEVEL=%s... You must use a valid level string defined in zerolog. Defaulting to 'info'", os.Getenv("LOG_LEVEL"))
			logLevel = zerolog.InfoLevel
		} else {
			logLevel = level
		}
	}
	zerolog.SetGlobalLevel(logLevel)

	cfToken := os.Getenv("CF_API_TOKEN")
	domain := os.Getenv("DOMAIN")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Info().Msg("Getting public IP...")
	ip, err := getPublicIP(ctx)
	if err != nil {
		log.Fatal().Msgf("failed to get public IP: %v", err)
	}
	log.Info().Msgf("Found public IP: %s", ip)

	log.Info().Msgf("Updating %s A record to %s", domain, ip)
	err = updateDomainRecord(ctx, cfToken, domain, ip)
	if err != nil {
		log.Fatal().Msgf("failed to update domain record: %v", err)
	}
	log.Info().Msg("Successfully updated record")
}
