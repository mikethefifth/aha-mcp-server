package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/aha-mcp-server/mcputil"
)

type GetIdeaParams struct {
	IdeaID string `json:"idea_id" description:"Idea ID to get"`
}

func (tc *ToolsClient) GetIdea(ctx context.Context, req *mcp.CallToolRequest, params GetIdeaParams) (*mcp.CallToolResult, any, error) {
	if resp, err := tc.simpleClient.Do(ctx, httpsimple.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("/api/v1/ideas/%s", params.IdeaID),
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("error getting Idea: %v", err), true), nil, err
	} else if ideaJSON, err := io.ReadAll(resp.Body); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error reading API response: %v", err), true), nil, err
	} else if jsonData, err := json.Marshal(map[string]any{
		"idea":        json.RawMessage(ideaJSON),
		"status_code": resp.StatusCode,
	}); err != nil {
		return mcputil.NewCallToolResultForAny(fmt.Sprintf("Error marshaling response: %v", err), true), nil, err
	} else {
		return mcputil.NewCallToolResultForAny(string(jsonData), false), nil, nil
	}
}

func GetIdeaTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_idea",
		Description: "Get Idea from Aha",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"idea_id": {
					Type:        "string",
					Description: "Idea ID to get",
				},
			},
			Required: []string{"idea_id"},
		},
	}
}
