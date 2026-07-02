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
	IdeaID      string `json:"idea_id" description:"Idea ID to get comments for (e.g., 'JPRO-I-2025')"`
	CommentType string `json:"comment_type,omitempty" description:"Type of comments to retrieve: 'all' (default, both internal and portal), 'internal' (internal comments only), or 'portal' (portal-visible comments only)"`
	Page        *int32 `json:"page,omitempty" description:"Page number"`
	PerPage     *int32 `json:"per_page,omitempty" description:"Results per page"`
}

func (tc *ToolsClient) ListIdeaComments(ctx context.Context, req *mcp.CallToolRequest, params ListIdeaCommentsParams) (*mcp.CallToolResult, any, error) {
	commentType := params.CommentType
	if commentType == "" {
		commentType = "all"
	}

	q := url.Values{}
	if params.Page != nil {
		q.Set("page", strconv.Itoa(int(*params.Page)))
	}
	if params.PerPage != nil {
		q.Set("per_page", strconv.Itoa(int(*params.PerPage)))
	}
	queryString := ""
	if len(q) > 0 {
		queryString = "?" + q.Encode()
	}

	result := map[string]any{
		"idea_id": params.IdeaID,
	}

	// Fetch internal comments
	if commentType == "all" || commentType == "internal" {
		internalURL := fmt.Sprintf("/api/v1/ideas/%s/comments%s", url.PathEscape(params.IdeaID), queryString)
		if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
			Method: http.MethodGet,
			URL:    internalURL,
		}); err != nil {
			result["internal_comments_error"] = fmt.Sprintf("error fetching internal comments: %v", err)
		} else if commentsJSON, err := io.ReadAll(resp.Body); err != nil {
			result["internal_comments_error"] = fmt.Sprintf("error reading internal comments: %v", err)
		} else {
			result["internal_comments"] = json.RawMessage(commentsJSON)
			result["internal_status_code"] = resp.StatusCode
		}
	}

	// Fetch portal comments
	if commentType == "all" || commentType == "portal" {
		portalURL := fmt.Sprintf("/api/v1/ideas/%s/idea_comments%s", url.PathEscape(params.IdeaID), queryString)
		if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
			Method: http.MethodGet,
			URL:    portalURL,
		}); err != nil {
			result["portal_comments_error"] = fmt.Sprintf("error fetching portal comments: %v", err)
		} else if commentsJSON, err := io.ReadAll(resp.Body); err != nil {
			result["portal_comments_error"] = fmt.Sprintf("error reading portal comments: %v", err)
		} else {
			result["portal_comments"] = json.RawMessage(commentsJSON)
			result["portal_status_code"] = resp.StatusCode
		}
	}

	if jsonData, err := json.Marshal(result); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil, nil
	}
}

func ListIdeaCommentsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_idea_comments",
		Description: "List comments on an Aha idea. Returns both internal (team-only) and portal-visible comments by default. Use comment_type to filter.",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"idea_id": {
					Type:        "string",
					Description: "Idea ID to get comments for (e.g., 'JPRO-I-2025')",
				},
				"comment_type": {
					Type:        "string",
					Description: "Type of comments to retrieve: 'all' (default, both internal and portal), 'internal' (internal/team comments only), or 'portal' (portal-visible comments only)",
					Enum:        []any{"all", "internal", "portal"},
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
