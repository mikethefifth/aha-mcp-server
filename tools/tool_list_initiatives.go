package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/aha-mcp-server/mcputil"
)

type ListInitiativesParams struct {
	Q            string `json:"q,omitempty" description:"Search term to match against initiative name"`
	UpdatedSince string `json:"updated_since,omitempty" description:"UTC timestamp (ISO8601). Only initiatives updated after this time"`
	Page         *int32 `json:"page,omitempty" description:"Page number"`
	PerPage      *int32 `json:"per_page,omitempty" description:"Results per page"`
}

func (tc *ToolsClient) ListInitiatives(ctx context.Context, req *mcp.CallToolRequest, params ListInitiativesParams) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	if params.Q != "" {
		q.Set("q", params.Q)
	}
	if params.UpdatedSince != "" {
		q.Set("updated_since", params.UpdatedSince)
	}
	if params.Page != nil {
		q.Set("page", strconv.Itoa(int(*params.Page)))
	}
	if params.PerPage != nil {
		q.Set("per_page", strconv.Itoa(int(*params.PerPage)))
	}

	apiURL := "/api/v1/initiatives"
	if len(q) > 0 {
		apiURL += "?" + q.Encode()
	}

	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    apiURL,
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error listing Initiatives: %v", err), true), nil, err
	} else if initiativesJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	} else if jsonData, err := json.Marshal(map[string]any{
		"initiatives": json.RawMessage(initiativesJSON),
		"status_code": resp.StatusCode,
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil, nil
	}
}

func ListInitiativesTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_initiatives",
		Description: "List initiatives from Aha with optional filtering and pagination",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"q": {
					Type:        "string",
					Description: "Search term to match against initiative name",
				},
				"updated_since": {
					Type:        "string",
					Description: "UTC timestamp (ISO8601). Only initiatives updated after this time",
				},
				"page": {
					Type:        "integer",
					Description: "Page number",
				},
				"per_page": {
					Type:        "integer",
					Description: "Results per page",
				},
			},
		},
	}
}
