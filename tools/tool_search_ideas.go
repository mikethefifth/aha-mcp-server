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

type SearchIdeasParams struct {
	Q             string `json:"q" description:"Query string to search against idea name, description, or ID"`
	IdeaPortalID  string `json:"idea_portal_id,omitempty" description:"Filter results to a specific ideas portal by numeric ID"`
	Fields        string `json:"fields,omitempty" description:"Comma-separated list of fields to return (e.g., 'vote_count,endorsements_count,duplicate_of'). Use '*' for all fields."`
	Page          *int32 `json:"page,omitempty" description:"Page number"`
	PerPage       *int32 `json:"per_page,omitempty" description:"Results per page"`
}

func (tc *ToolsClient) SearchIdeas(ctx context.Context, req *mcp.CallToolRequest, params SearchIdeasParams) (*mcp.CallToolResult, any, error) {
	if params.Q == "" {
		return mcputil.NewCallToolResultForAny("Query parameter 'q' is required", true), nil, fmt.Errorf("query parameter 'q' is required")
	}

	q := url.Values{}
	q.Set("q", params.Q)
	if params.IdeaPortalID != "" {
		q.Set("idea_portal_id", params.IdeaPortalID)
	}
	if params.Fields != "" {
		q.Set("fields", params.Fields)
	}
	if params.Page != nil {
		q.Set("page", strconv.Itoa(int(*params.Page)))
	}
	if params.PerPage != nil {
		q.Set("per_page", strconv.Itoa(int(*params.PerPage)))
	}

	apiURL := "/api/v1/ideas/related?" + q.Encode()

	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    apiURL,
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error searching Ideas: %v", err), true), nil, err
	} else if ideasJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	} else if jsonData, err := json.Marshal(map[string]any{
		"ideas":       json.RawMessage(ideasJSON),
		"status_code": resp.StatusCode,
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil, nil
	}
}

func SearchIdeasTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "search_ideas",
		Description: "Search ideas by keyword across name, description, and ID. Use this for full-text search; use list_ideas for filtering by status, product, tags, etc.",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"q": {
					Type:        "string",
					Description: "Query string to search against idea name, description, or ID",
				},
				"idea_portal_id": {
					Type:        "string",
					Description: "Filter results to a specific ideas portal by numeric ID",
				},
				"fields": {
					Type:        "string",
					Description: "Comma-separated list of fields to return (e.g., 'vote_count,endorsements_count,duplicate_of,created_by'). Use '*' for all fields. Useful for getting votes, endorsements, and merge info.",
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
			Required: []string{"q"},
		},
	}
}
