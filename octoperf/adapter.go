package octoperf

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

const baseUrl = "https://api.octoperf.com"

func NewClient() (*Client, error) {
	var apiKey = os.Getenv("OCTOPERF_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OCTOPERF_API_KEY environment variable is not set")
	}

	var client = resty.New().
		SetBaseURL(baseUrl).
		SetAuthToken(apiKey)

	return &Client{client: client}, nil
}

// Function to check the status of an OctoPerf test
func (obj Client) GetTestStatus(ctx context.Context, benchResultId string) (string, error) {
	resp, err := obj.client.R().
		SetContext(ctx).
		SetPathParam("benchResultId", benchResultId).
		Get("/runtime/bench-results/progress/{benchResultId}")
	if err != nil {
		return "", err
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("OctoPerf API error (code: %d): %s", resp.StatusCode(), resp.String())
	}

	return resp.String(), nil
}

// Function to start an OctoPerf test with the runtime ID from the first scenario
func (obj Client) RunOctoPerfTest(ctx context.Context, runtimeId string) (string, error) {

	resp, err := obj.client.R().
		SetContext(ctx).
		SetPathParam("runtimeId", runtimeId).
		Post("/runtime/scenarios/run/{runtimeId}")
	if err != nil {
		return "", err
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("OctoPerf API error (code: %d): %s", resp.StatusCode(), resp.String())
	}

	return resp.String(), nil
}

// Function to retrieve details of an OctoPerf test report
func (obj Client) GetReportDetails(ctx context.Context, reportId string) (string, error) {

	resp, err := obj.client.R().
		SetContext(ctx).
		SetPathParam("reportId", reportId).
		Get("/analysis/bench-reports/{reportId}")
	if err != nil {
		return "", err
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("OctoPerf API error (code: %d): %s", resp.StatusCode(), resp.String())
	}

	return resp.String(), nil
}

// Utility function to determine the metric type based on its ID
func (obj Client) getMetricType(ctx context.Context, metricId string) string {
	// Assign the appropriate type according to the metric ID
	switch metricId {
	case "RESPONSE_TIME_AVG", "LATENCY_STD":
		return "CONTAINER"
	case "RESPONSE_TIME_PERCENTILE_90", "RESPONSE_TIME_PERCENTILE_95":
		return "HIT"
	case "HITS_TOTAL", "ERRORS_TOTAL", "ERRORS_PERCENT", "THROUGHPUT_RATE", "HITS_RATE":
		return "HIT"
	default:
		return "HIT" // Default type
	}
}

// Function to retrieve specific metrics from an OctoPerf benchmark
func (obj Client) GetMetricDetail(ctx context.Context, benchResultId string, metricIds []string) (string, error) {

	// Build the request body with the requested metrics
	metrics := []map[string]interface{}{}

	for _, metricId := range metricIds {
		metric := map[string]interface{}{
			"id":            metricId,
			"type":          obj.getMetricType(ctx, metricId),
			"filters":       []interface{}{},
			"benchResultId": benchResultId,
			"configs":       []interface{}{},
		}
		metrics = append(metrics, metric)
	}

	reqBody := map[string]interface{}{
		"@type":   "SummaryReportItem",
		"metrics": metrics,
		"id":      "",
		"name":    "Statistics summary",
	}

	log.Println("Sending metrics request with body:", reqBody)

	resp, err := obj.client.R().
		SetContext(ctx).
		SetBody(reqBody).
		Post("/analysis/metrics/summary")

	if err != nil {
		return "", err
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("OctoPerf API error (code: %d): %s", resp.StatusCode(), resp.String())
	}

	return resp.String(), nil
}

// Function to retrieve workspaces of the current user
func (obj Client) GetCurrentUserWorkspaces(ctx context.Context) (string, error) {

	resp, err := obj.client.R().
		SetContext(ctx).
		Get("/workspaces/member-of")

	if err != nil {
		return "", fmt.Errorf("error retrieving workspaces: %v", err)
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("OctoPerf API error (code: %d): %s", resp.StatusCode(), resp.String())
	}

	return string(resp.Body()), nil
}

// Function to retrieve projects by workspace ID
func (obj Client) GetProjectsByWorkspaceId(ctx context.Context, workspaceId string) (string, error) {

	resp, err := obj.client.R().
		SetContext(ctx).
		SetPathParam("workspaceId", workspaceId).
		Get("/design/projects/by-workspace/{workspaceId}/DESIGN")

	if err != nil {
		return "", fmt.Errorf("error retrieving projects: %v", err)
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("OctoPerf API error (code: %d): %s", resp.StatusCode(), resp.String())
	}

	return string(resp.Body()), nil
}

// Function to retrieve test typologies for a project (runtime IDs)
func (obj Client) GetRuntimeIds(ctx context.Context, projectId string) (string, error) {

	resp, err := obj.client.R().
		SetContext(ctx).
		SetPathParam("projectId", projectId).
		Get("/runtime/scenarios/by-project/{projectId}")

	if err != nil {
		return "", err
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("OctoPerf API error (code: %d): %s", resp.StatusCode(), resp.String())
	}

	return resp.String(), nil
}
