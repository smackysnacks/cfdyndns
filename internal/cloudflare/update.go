package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) UpdateDnsRecord(ctx context.Context, zoneId string, recordId string, data DnsUpdateRequestData) (*CloudflareResponse, error) {
	b, _ := json.Marshal(data)
	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		fmt.Sprintf("%szones/%s/dns_records/%s", BaseUrl, zoneId, recordId),
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, fmt.Errorf("constructing new request: %v", err)
	}

	c.addStandardHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %v", err)
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status [%d]: %v", resp.StatusCode, string(content))
	}

	var cloudflareResponse CloudflareResponse
	_ = json.Unmarshal(content, &cloudflareResponse)

	return &cloudflareResponse, nil
}
