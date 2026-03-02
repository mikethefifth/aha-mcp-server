package mcputil

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewCallToolResultForAny(msg string, isError bool) *mcp.CallToolResult {
	result := &mcp.CallToolResult{
		IsError: isError,
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
	}
	var structured interface{}
	if err := json.Unmarshal([]byte(msg), &structured); err == nil {
		result.StructuredContent = structured
	}
	return result
}
