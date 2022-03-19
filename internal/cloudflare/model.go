package cloudflare

import "time"

type CloudflareResponse struct {
	Result     *Result       `json:"result"`
	Success    bool          `json:"success"`
	Errors     []Errors      `json:"errors"`
	Messages   []interface{} `json:"messages"`
	ResultInfo ResultInfo    `json:"result_info"`
}

type Result struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	TTL        int       `json:"ttl"`
	Priority   int       `json:"priority"`
	Locked     bool      `json:"locked"`
	ZoneID     string    `json:"zone_id"`
	ZoneName   string    `json:"zone_name"`
	ModifiedOn time.Time `json:"modified_on"`
	CreatedOn  time.Time `json:"created_on"`
}

type Errors struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
}

type DnsUpdateRequestData struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`

	Proxied bool `json:"proxied"`
}
