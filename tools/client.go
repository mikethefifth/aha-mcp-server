package tools

import (
	"fmt"

	ao "github.com/grokify/goauth/aha"
	"github.com/grokify/mogo/net/http/httpsimple"
)

type ToolsClient struct {
	simpleClient *httpsimple.Client
}

func NewToolsClient(ahaSubdomain, ahaAPIKey string) (*ToolsClient, error) {
	hc, err := ao.NewClient(ahaSubdomain, ahaAPIKey)
	if err != nil {
		return nil, err
	}
	baseURL := fmt.Sprintf("https://%s.aha.io/", ahaSubdomain)
	sc := httpsimple.NewClient(hc, baseURL)
	return &ToolsClient{
		simpleClient: &sc,
	}, nil
}
