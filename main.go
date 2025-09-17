package main

import (
	"context"
	"log/slog"
	"mcp-octoperf/octoperf"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var ctx = context.Background()

	// Add logs for debugging
	slog.InfoContext(ctx, "Starting OctoPerf MCP server...")

	// Create client for OctoPerf API
	octoperfClient, err := octoperf.NewClient()
	if err != nil {
		slog.ErrorContext(ctx, "Error creating OctoPerf client", "error", err)
	}

	// Create handler for OctoPerf API
	var handler = Handler{
		Client: octoperfClient,
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"OctoPerf MCP Server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Tool to launch performance test with OctoPerf
	slog.InfoContext(ctx, "Adding tool", "name", "octoperf_run_test")
	s.AddTool(
		mcp.NewTool("octoperf_run_test",
			mcp.WithDescription("Start an OctoPerf performance test with the specified runtime Id"),
			mcp.WithString("runtimeId",
				mcp.Required(),
				mcp.Description("The runtime Id of the OctoPerf test to start"),
			),
		),
		handler.RunOctoPerfTest,
	)

	// Tool to check the status of an OctoPerf test
	slog.InfoContext(ctx, "Adding tool", "name", "octoperf_status")
	s.AddTool(
		mcp.NewTool("octoperf_status",
			mcp.WithDescription("Check the status of an OctoPerf performance test"),
			mcp.WithString("benchResultId",
				mcp.Required(),
				mcp.Description("The benchmark Id for which to check the status"),
			),
		),
		handler.GetTestStatus,
	)

	// New tool to retrieve test report details
	slog.InfoContext(ctx, "Adding tool", "name", "octoperf_report")
	s.AddTool(
		mcp.NewTool("octoperf_report",
			mcp.WithDescription("Retrieve the details of an OctoPerf test report"),
			mcp.WithString("reportId",
				mcp.Required(),
				mcp.Description("The Id of the test report to retrieve"),
			),
		),
		handler.GetReportDetails,
	)

	// New tool to retrieve specific metrics
	slog.InfoContext(ctx, "Adding tool", "name", "octoperf_get_report_metrics")
	s.AddTool(
		mcp.NewTool("octoperf_get_report_metrics",
			mcp.WithDescription("Retrieve specific metrics from an OctoPerf test report"),
			mcp.WithString("benchResultId",
				mcp.Required(),
				mcp.Description("The benchmark Id for which to retrieve metrics"),
			),
			mcp.WithArray("metricIds",
				mcp.Required(),
				mcp.Description("List of metric Ids to retrieve (e.g., HITS_TOTAL, ERRORS_TOTAL, etc.)"),
				mcp.Items(map[string]any{
					"type": "string",
				}),
			),
		),
		handler.GetMetricDetail,
	)

	// New tool to retrieve current user workspaces
	slog.InfoContext(ctx, "Adding tool", "name", "octoperf_get_current_user_workspaces")
	s.AddTool(
		mcp.NewTool("octoperf_get_current_user_workspaces",
			mcp.WithDescription("Retrieve the workspaces of the current user"),
		),
		handler.GetCurrentUserWorkspaces,
	)

	// New tool to retrieve projects by workspace Id
	slog.InfoContext(ctx, "Adding tool", "name", "get_project_by_workspace_id")
	s.AddTool(
		mcp.NewTool("get_project_by_workspace_id",
			mcp.WithDescription("Retrieve the projects linked to a specific workspace"),
			mcp.WithString("workspaceId",
				mcp.Required(),
				mcp.Description("The workspace Id for which to retrieve projects"),
			),
		),
		handler.GetProjectsByWorkspaceId,
	)

	// New tool to retrieve Runtime Ids for a project
	slog.InfoContext(ctx, "Adding tool", "name", "get_runtime_id")
	s.AddTool(
		mcp.NewTool("get_runtime_id",
			mcp.WithDescription("Retrieve the available runtime Ids for a project"),
			mcp.WithString("projectId",
				mcp.Description("The project Id for which to retrieve runtime Ids (will use the default PROJECT_ID if not specified)"),
			),
		),
		handler.GetRuntimeIds,
	)

	// Start the MCP server
	slog.InfoContext(ctx, "Starting MCP server...")
	if err := server.ServeStdio(s); err != nil {
		slog.ErrorContext(ctx, "Server error", "error", err)
	}
}
