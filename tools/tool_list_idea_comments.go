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

type ListIdeaCommentsParams struct {
	IdeaID  string `json:"idea_id" description:"Idea ID to get comments for (e.g., 'JPRO-I-2025')"`
	Page    *int32 `json:"page,omitempty" description:"Page number"`
	PerPage *int32 `json:"per_page,omitempty" description:"Results per page"`
}

func (tc *ToolsClient) ListIdeaComments(ctx context.Context, req *mcp.CallToolRequest, params ListIdeaCommentsParams) (*mcp.CallToolResult, any, error) {
	q := url.Values{}
	if params.Page != nil {
		q.Set("page", strconv.Itoa(int(*params.Page)))
	}
	if params.PerPage != nil {
		q.Set("per_page", strconv.Itoa(int(*params.PerPage)))
	}

	apiURL := fmt.Sprintf("/api/v1/ideas/%s/idea_comments", url.PathEscape(params.IdeaID))
	if len(q) > 0 {
		apiURL += "?" + q.Encode()
	}

	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    apiURL,
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error listing idea comments: %v", err), true), nil, err
	} else if commentsJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	} else if jsonData, err := json.Marshal(map[string]any{
		"idea_comments": json.RawMessage(commentsJSON),
		"status_code":   resp.StatusCode,
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil, nil
	}
}

func ListIdeaCommentsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_idea_comments",
		Description: "List comments on an Aha idea",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"idea_id": {
					Type:        "string",
					Description: "Idea ID to get comments for (e.g., 'JPRO-I-2025')",
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
			Required: []string{"idea_id"},
		},
	}
}
