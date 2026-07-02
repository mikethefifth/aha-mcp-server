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

type ListIdeasParams struct {
	ProductID      string `json:"product_id,omitempty" description:"Filter by product ID or key (e.g., 'EXP'). When provided, only ideas from this product are returned."`
	Q              string `json:"q,omitempty" description:"Search term to match against the idea name"`
	Spam           *bool  `json:"spam,omitempty" description:"When true, shows ideas marked as spam"`
	WorkflowStatus string `json:"workflow_status,omitempty" description:"Filter by workflow status ID or name"`
	Sort           string `json:"sort,omitempty" description:"Sort by: recent, trending, or popular"`
	CreatedBefore  string `json:"created_before,omitempty" description:"UTC timestamp (ISO8601). Only ideas created before this time"`
	CreatedSince   string `json:"created_since,omitempty" description:"UTC timestamp (ISO8601). Only ideas created after this time"`
	UpdatedSince   string `json:"updated_since,omitempty" description:"UTC timestamp (ISO8601). Only ideas updated after this time"`
	Tag            string `json:"tag,omitempty" description:"Filter by tag value"`
	UserID         string `json:"user_id,omitempty" description:"Filter by creator user ID"`
	IdeaUserID     string `json:"idea_user_id,omitempty" description:"Filter by idea user ID"`
	Fields         string `json:"fields,omitempty" description:"Comma-separated list of fields to return (e.g., 'vote_count,endorsements_count,duplicate_of'). Use '*' for all fields."`
	Page           *int32 `json:"page,omitempty" description:"Page number"`
	PerPage        *int32 `json:"per_page,omitempty" description:"Results per page"`
}

func (tc *ToolsClient) ListIdeas(ctx context.Context, req *mcp.CallToolRequest, params ListIdeasParams) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	if params.Q != "" {
		q.Set("q", params.Q)
	}
	if params.Spam != nil {
		q.Set("spam", strconv.FormatBool(*params.Spam))
	}
	if params.WorkflowStatus != "" {
		q.Set("workflow_status", params.WorkflowStatus)
	}
	if params.Sort != "" {
		q.Set("sort", params.Sort)
	}
	if params.CreatedBefore != "" {
		q.Set("created_before", params.CreatedBefore)
	}
	if params.CreatedSince != "" {
		q.Set("created_since", params.CreatedSince)
	}
	if params.UpdatedSince != "" {
		q.Set("updated_since", params.UpdatedSince)
	}
	if params.Tag != "" {
		q.Set("tag", params.Tag)
	}
	if params.UserID != "" {
		q.Set("user_id", params.UserID)
	}
	if params.IdeaUserID != "" {
		q.Set("idea_user_id", params.IdeaUserID)
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

	var apiURL string
	if params.ProductID != "" {
		apiURL = fmt.Sprintf("/api/v1/products/%s/ideas", url.PathEscape(params.ProductID))
	} else {
		apiURL = "/api/v1/ideas"
	}
	if len(q) > 0 {
		apiURL += "?" + q.Encode()
	}

	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    apiURL,
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error listing Ideas: %v", err), true), nil, err
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

func ListIdeasTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_ideas",
		Description: "List ideas with filtering by product, status, tags, dates, or user. For keyword search across name/description, use search_ideas instead.",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"product_id": {
					Type:        "string",
					Description: "Filter by product ID or key (e.g., 'EXP'). When provided, only ideas from this product are returned.",
				},
				"q": {
					Type:        "string",
					Description: "Search term to match against the idea name (name only, not full-text). For searching description/body, use search_ideas tool.",
				},
				"spam": {
					Type:        "boolean",
					Description: "When true, shows ideas marked as spam",
				},
				"workflow_status": {
					Type:        "string",
					Description: "Filter by workflow status ID or name",
				},
				"sort": {
					Type:        "string",
					Description: "Sort by: recent, trending, or popular",
					Enum:        []any{"recent", "trending", "popular"},
				},
				"created_before": {
					Type:        "string",
					Description: "UTC timestamp (ISO8601). Only ideas created before this time",
				},
				"created_since": {
					Type:        "string",
					Description: "UTC timestamp (ISO8601). Only ideas created after this time",
				},
				"updated_since": {
					Type:        "string",
					Description: "UTC timestamp (ISO8601). Only ideas updated after this time",
				},
				"tag": {
					Type:        "string",
					Description: "Filter by tag value",
				},
				"user_id": {
					Type:        "string",
					Description: "Filter by creator user ID",
				},
				"idea_user_id": {
					Type:        "string",
					Description: "Filter by idea user ID",
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
		},
	}
}
