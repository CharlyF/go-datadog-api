package datadog

import (
	"fmt"
	"net/http"
	"strings"
)

// The list of Rate Limited Endpoints of the Datadog API.
// https://docs.datadoghq.com/api/?lang=bash#rate-limiting
var (
	rateLimitedEndpoints = map[string]string{
		"/v1/query":               "GET",
		"/v1/input":               "GET",
		"/v1/metrics":             "GET",
		"/v1/events":              "POST",
		"/v1/logs-queries/list":   "POST",
		"/v1/graph/snapshot":      "GET",
		"/v1/logs/config/indexes": "GET",
	}
)

func isRateLimited(method string, endpoint string) (limited bool, shortEndpoint string) {
	for e, m := range rateLimitedEndpoints {
		if strings.HasPrefix(endpoint, e) && m == method {
			return true, e
		}
	}
	return false, ""
}

func (client *Client) updateRateLimits(resp *http.Response, api string) error {
	if resp.Header == nil {
		return fmt.Errorf("header missing from the HTTP response.")
	}
	client.RateLimitingStats[api] = RateLimit{
		limit:     resp.Header.Get("X-RateLimit-Limit"),
		reset:     resp.Header.Get("X-RateLimit-Reset"),
		period:    resp.Header.Get("X-RateLimit-Period"),
		remaining: resp.Header.Get("X-RateLimit-Remaining"),
	}
	return nil
}
