package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mcp-octoperf/octoperf"

	"github.com/mark3labs/mcp-go/mcp"
)

type Handler struct {
	Client *octoperf.Client
}

func (obj *Handler) GetTestStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	benchResultId, err := request.RequireString("benchResultId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp, err := obj.Client.GetTestStatus(ctx, benchResultId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error checking status: %v", err)), nil
	}

	resultData := map[string]interface{}{
		"status":   "retrieved",
		"response": resp,
	}

	return respondAsJsonText(resultData), nil
}

func (obj *Handler) RunOctoPerfTest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	runtimeId, err := request.RequireString("runtimeId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp, err := obj.Client.RunOctoPerfTest(ctx, runtimeId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error running test: %v", err)), nil
	}

	resultData := map[string]interface{}{
		"status":   "test_started",
		"response": resp,
	}

	return respondAsJsonText(resultData), nil
}

func (obj *Handler) GetReportDetails(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	reportId, err := request.RequireString("reportId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp, err := obj.Client.GetReportDetails(ctx, reportId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving report: %v", err)), nil
	}

	resultData := map[string]interface{}{
		"status":   "report_retrieved",
		"response": resp,
	}

	return respondAsJsonText(resultData), nil
}

func (obj *Handler) GetMetricDetail(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	var inputs struct {
		BenchResultId string   `json:"benchResultId"`
		MetricIds     []string `json:"metricIds"`
	}

	if err := request.BindArguments(&inputs); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp, err := obj.Client.GetMetricDetail(ctx, inputs.BenchResultId, inputs.MetricIds)

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving metrics: %v", err)), nil
	}
	resultData := map[string]interface{}{
		"status":   "metrics_retrieved",
		"response": resp,
	}

	return respondAsJsonText(resultData), nil
}

func respondAsJsonText(data any) *mcp.CallToolResult {
	jsonResult, err := json.Marshal(data)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error during JSON serialization: %v", err))
	}

	return mcp.NewToolResultText(string(jsonResult))
}

func (obj *Handler) GetCurrentUserWorkspaces(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	resp, err := obj.Client.GetCurrentUserWorkspaces(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving user workspaces: %v", err)), nil
	}
	resultData := map[string]interface{}{
		"status":   "retrieved",
		"response": resp,
	}

	return respondAsJsonText(resultData), nil
}

func (obj *Handler) GetProjectsByWorkspaceId(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workspaceId, err := request.RequireString("workspaceId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp, err := obj.Client.GetProjectsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving projects for workspace %s: %v", workspaceId, err)), nil
	}
	resultData := map[string]interface{}{
		"status":      "retrieved",
		"workspaceId": workspaceId,
		"response":    resp,
	}
	return respondAsJsonText(resultData), nil
}

func (obj *Handler) GetRuntimeIds(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectId, err := request.RequireString("projectId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resp, err := obj.Client.GetRuntimeIds(ctx, projectId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error retrieving Runtime IDs for project %s: %v", projectId, err)), nil
	}
	resultData := map[string]interface{}{
		"status":    "retrieved",
		"projectId": projectId,
		"response":  resp,
	}
	return respondAsJsonText(resultData), nil
}
