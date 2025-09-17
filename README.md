# MCP OctoPerf

[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/tpit59/octoperf-mcp)

A Model Context Protocol (MCP) server for OctoPerf to launch and manage performance tests via AI assistants like Claude. ( Proof of Concept )

## 🏗️ Architecture

This project implements an MCP server that bridges AI assistants and the OctoPerf API:

```
┌──────────────────┐    MCP Protocol    ┌──────────────────┐    HTTP API     ┌─────────────────┐
│  AI/Claude/[...] │◄──────────────────►│   MCP Server     │◄───────────────►│  OctoPerf API   │
│                  │     JSON-RPC       │  (this project)  │    REST calls   │                 │
└──────────────────┘                    └──────────────────┘                 └─────────────────┘
```

### Project structure

```
/workspaces/mvp-mcp-octoperf/
├── main.go              # Entry point + MCP tools declaration
├── handler.go           # MCP handlers (interface layer)
├── octoperf/
│   ├── adapter.go       # HTTP client for OctoPerf API
│   └── dto.go          # Data structures
├── .vscode/
│   └── mcp.json        # MCP configuration (ready to use)
├── .devcontainer/      # Dev Container configuration
└── go.mod             # Go dependencies
```

## ⚙️ Configuration

### MCP Configuration with mcp.json (Recommended Method)

The MCP server uses a configuration file `mcp.json` located in the `.vscode` directory. This file contains connection information and necessary environment variables.

#### Configuration Steps:

1. **Create the .vscode directory** (if not already present):
   ```bash
   mkdir -p .vscode
   ```

2. **Create the mcp.json file**:
   
   Create a `.vscode/mcp.json` configuration file in your project.

3. **Understanding the mcp.json structure**:
   
   The `.vscode/mcp.json` file uses VS Code's input prompt system:
   
   ```json
   {
     "inputs": [
       {
         "type": "promptString",
         "id": "octoperf-api-key",
         "description": "OctoPerf API Key",
         "password": true
       }
     ],
     "servers": {
       "octoperf": {
         "command": "go",
         "args": [
            "run",
            "."
         ],
         "type": "stdio",
         "env": {
           "OCTOPERF_API_KEY": "${input:octoperf-api-key}"
         }
       }
     }
   }
   ```

4. **How the API Key Prompt Works**:
   - **First startup**: VS Code will automatically prompt you to enter your OctoPerf API Key
   - **Secure storage**: The API key is stored locally as a password in VS Code's secure storage
   - **Subsequent startups**: The stored key is automatically retrieved, no prompt needed
   - **Security**: The API key is never stored in plain text in configuration files

5. **Get your OctoPerf API Key**:
   - Log into your OctoPerf account
   - Go to your profile settings
   - Copy or generate your API key
   - Copy the API key (you'll paste it when prompted by VS Code)

6. **No build required**:
   The MCP server runs directly with `go run .` - no compilation step needed.

### Required environment variables

- `OCTOPERF_API_KEY`: API key for OctoPerf authentication (required) - automatically prompted and stored securely

### Alternative configuration methods

1. **System environment variables**
   ```bash
   export OCTOPERF_API_KEY=your_api_key_here
   ```

2. **Direct VS Code configuration**
   - The `mcp.json` file is automatically detected by VS Code
   - Environment variables are injected automatically
   - API key prompting provides enhanced security

## 🚀 Installation and development

### Prerequisites

- Docker (required for Dev Container)
- Go 1.23 or higher
- Valid OctoPerf API key
- VS Code or other IDE

### Local installation

```bash
# Clone the repository
git clone [REPO_URL]
cd octoperf-mcp

# Install dependencies
go mod download
```

You can start the OctoPerf MCP Server through your IDE.

### Development with Dev Container

```bash
# Open in VS Code with Dev Containers
code .
# VS Code will suggest opening in a Dev Container
# Create the .vscode/mcp.json configuration file as described above
```

## 🛠️ Available MCP Tools

The server exposes 7 tools to interact with OctoPerf:

### 📋 Resource Discovery

| Tool | Description | Parameters |
|------|-------------|------------|
| `octoperf_get_current_user_workspaces` | List all accessible workspaces | None |
| `get_project_by_workspace_id` | List projects in a workspace | `workspaceId` (required) |
| `get_runtime_id` | Retrieve runtime IDs for a project | `projectId` (optional) |

### 🚀 Test Execution

| Tool | Description | Parameters |
|------|-------------|------------|
| `octoperf_run_test` | Start a performance test | `runtimeId` (required) |
| `octoperf_status` | Check test status | `benchResultId` (required) |

### 📊 Results Analysis

| Tool | Description | Parameters |
|------|-------------|------------|
| `octoperf_report` | Retrieve report details | `reportId` (required) |
| `octoperf_get_report_metrics` | Retrieve specific metrics | `benchResultId` (required), `metricIds[]` (required) |

### Available metrics

The following metrics can be retrieved with `octoperf_get_report_metrics`:

- **Performance**: `RESPONSE_TIME_AVG`, `RESPONSE_TIME_PERCENTILE_95`, `RESPONSE_TIME_PERCENTILE_99`
- **Volume**: `HITS_TOTAL`, `HITS_SUCCESSFUL_TOTAL`, `THROUGHPUT_TOTAL`
- **Errors**: `ERRORS_TOTAL`, `ERRORS_RATE`, `ERRORS_PERCENT`
- **Network**: `SENT_BYTES_TOTAL`, `CONNECT_TIME_AVG`, `LATENCY_AVG`

## 🔄 Typical usage workflow

### 1. Resource discovery
```bash
# 1. List workspaces
octoperf_get_current_user_workspaces

# 2. List projects in a workspace
get_project_by_workspace_id workspaceId="abc123"

# 3. Retrieve runtime IDs for a project
get_runtime_id projectId="def456"
```

### 2. Test execution
```bash
# 4. Start a test
octoperf_run_test runtimeId="ghi789"

# 5. Check status
octoperf_status benchResultId="jkl012"
```

### 3. Results analysis
```bash
# 6. Retrieve a complete report
octoperf_report reportId="mno345"

# 7. Retrieve specific metrics
octoperf_get_report_metrics benchResultId="jkl012" metricIds=["HITS_TOTAL","ERRORS_TOTAL","RESPONSE_TIME_AVG"]
```

## 🔧 MCP Configuration for VS Code

The `.vscode/mcp.json` file automatically configures the MCP server in VS Code with secure API key handling.

**Detailed structure of .vscode/mcp.json file:**

```json
   {
     "inputs": [
       {
         "type": "promptString",
         "id": "octoperf-api-key",
         "description": "OctoPerf API Key",
         "password": true
       }
     ],
     "servers": {
       "octoperf": {
         "command": "go",
         "args": [
            "run",
            "."
         ],
         "type": "stdio",
         "env": {
           "OCTOPERF_API_KEY": "${input:octoperf-api-key}"
         }
       }
     }
   }
```

**Key Features:**
- **Secure API Key Storage**: Uses VS Code's secure credential storage
- **Automatic Prompting**: API key is requested only on first use
- **No Plain Text**: API keys are never stored in configuration files
- **Seamless Integration**: Works automatically with Claude Dev and other MCP clients

**First-time Setup Flow:**
1. Clone the repository 
2. Create the `.vscode/mcp.json` configuration file as shown above
3. Start VS Code or restart the MCP server
4. VS Code will prompt: "OctoPerf API Key"
5. Enter your API key (input will be masked as it's marked as password)
6. The key is securely stored and automatically used for future sessions

**Alternative for local development:**
```json
{
  "inputs": [
    {
      "type": "promptString",
      "id": "octoperf-api-key",
      "description": "OctoPerf API Key",
      "password": true
    }
  ],
  "servers": {
    "octoperf": {
      "command": "go",
      "args": ["run", "."],
      "cwd": "/path/to/your/mvp-mcp-octoperf",
      "type": "stdio",
      "env": {
        "OCTOPERF_API_KEY": "${input:octoperf-api-key}"
      }
    }
  }
}
```

## 🏛️ Technical architecture

### Abstraction layers

1. **Handler Layer** (`handler.go`): MCP Interface
   - MCP parameter extraction
   - Type validation and conversion
   - Response formatting

2. **Business Layer** (`octoperf/adapter.go`): OctoPerf Client
   - API authentication
   - HTTP request construction
   - API error handling

3. **Transport Layer**: HTTP with Resty
   - Connection management
   - Retry and timeouts
   - Request logging

### Main dependencies

```go
require (
    github.com/mark3labs/mcp-go v0.35.0      // MCP Framework
    github.com/go-resty/resty/v2 v2.16.5     // HTTP Client
)
```

## 🔐 Security

- ✅ API keys are never stored in plain text
- ✅ Uses VS Code's secure credential storage system
- ✅ API key prompting with masked input (password field)
- ✅ `.vscode/mcp.json` contains no sensitive data
- ✅ Automatic Bearer authentication
- ✅ Input parameter validation
- ⚠️ No log encryption (contains API responses)

**Security Best Practices:**
- The `mcp.json` file is safe to commit (contains no secrets)
- API keys are stored in VS Code's secure credential manager
- Regenerate API keys regularly in OctoPerf dashboard
- Create the configuration file manually for enhanced security

## 🚨 Known limitations

- No support for parallel tests
- Fixed API timeout (not configurable)
- Basic network error handling
- No API result caching

## 📄 License

Distributed under the [MIT](./LICENSE) License. 

You are free to use, modify, and redistribute it under the terms of this license.

## 📞 Support

- OctoPerf Documentation: https://doc.octoperf.com/
- MCP Protocol: https://modelcontextprotocol.io/
- GitHub Issues: [Create a ticket](https://github.com/tpit59/octoperf-mcp/issues)