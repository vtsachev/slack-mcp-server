package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/korotovsky/slack-mcp-server/pkg/handler"
	"github.com/korotovsky/slack-mcp-server/pkg/provider"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server *server.MCPServer
}

func NewMCPServer(provider *provider.ApiProvider) *MCPServer {
	s := server.NewMCPServer(
		"Slack MCP Server",
		"1.0.0",
		server.WithLogging(),
		server.WithRecovery(),
	)

	conversationsHandler := handler.NewConversationsHandler(provider)

	s.AddTool(mcp.NewTool("conversations_history",
		mcp.WithDescription("Get messages from the channel by channel_id, the last row/column in the response is used as 'cursor' parameter for pagination if not empty"),
		mcp.WithString("channel_id",
			mcp.Required(),
			mcp.Description("ID of the channel in format Cxxxxxxxxxx"),
		),
		mcp.WithString("cursor",
			mcp.Description("Cursor for pagination. Use the value of the last row and column in the response as next_cursor field returned from the previous request."),
		),
		mcp.WithString("limit",
			mcp.DefaultString("1d"),
			mcp.Description("Limit of messages to fetch in format of maximum ranges of time (e.g. 1d - 1 day, 30d - 30 days, 90d - 90 days which is a default limit for free tier history) or number of messages (e.g. 50). Must be empty when 'cursor' is provided."),
		),
	), conversationsHandler.ConversationsHistoryHandler)

	channelsHandler := handler.NewChannelsHandler(provider)

	s.AddTool(mcp.NewTool("channels_list",
		mcp.WithDescription("Get list of channels"),
		mcp.WithString("channel_types",
			mcp.Required(),
			mcp.Description("Comma-separated channel types. Allowed values: 'mpim', 'im', 'public_channel', 'private_channel'. Example: 'public_channel,private_channel,im'"),
		),
		mcp.WithString("sort",
			mcp.Description("Type of sorting. Allowed values: 'popularity' - sort by number of members/participants in each channel."),
		),
		mcp.WithNumber("limit",
			mcp.DefaultNumber(100),
			mcp.Description("The maximum number of items to return. Must be an integer under 1000."),
		),
		mcp.WithString("cursor",
			mcp.Description("Cursor for pagination. Use the value of the last row and column in the response as next_cursor field returned from the previous request."),
		),
	), channelsHandler.ChannelsHandler)

	return &MCPServer{
		server: s,
	}
}

func (s *MCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server,
		server.WithBaseURL(fmt.Sprintf("http://%s", addr)),
		server.WithSSEContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			apiKey := os.Getenv("SLACK_MCP_SSE_API_KEY")
			// Use the authFromRequest function (now part of this file) to extract the token
			ctxWithAuth := authFromRequest(ctx, r)
			token, _ := tokenFromContext(ctxWithAuth) // Extract token put in context

			if apiKey != "" { // SLACK_MCP_SSE_API_KEY is set
				// Remove "Bearer " prefix if present
				if len(token) > 7 && token[:7] == "Bearer " {
					token = token[7:]
				}

				if token == "" {
					// API key is configured, but no token was provided.
					// Return a context that indicates this specific unauthorized state.
					// Downstream, the mcp-go server library's SSE handler should ideally use this
					// information to write an HTTP 401 error before trying to establish the SSE stream.
					// For now, we're marking the context.
					return context.WithValue(context.Background(), authKey{}, "unauthorized_sse_token_missing")
				}
				if token != apiKey {
					// API key is configured, and the provided token is invalid.
					// Return a context that indicates this specific unauthorized state.
					return context.WithValue(context.Background(), authKey{}, "unauthorized_sse_token_invalid")
				}
			}
			// If apiKey is not set, or if it is set and token is valid, proceed with the original context
			// which contains the "Authorization" header value (if any).
			return ctxWithAuth
		}),
	)
}

func (s *MCPServer) ServeStdio() error {
	return server.ServeStdio(s.server)
}

// authKey is a custom context key for storing the auth token.
// Moved here from sse_auth.go to be used directly in ServeSSE.
type authKey struct{}

// withAuthKey adds an auth key to the context.
// Moved here from sse_auth.go.
func withAuthKey(ctx context.Context, auth string) context.Context {
	return context.WithValue(ctx, authKey{}, auth)
}

// tokenFromContext extracts the auth token from the context.
// Moved here from sse_auth.go.
func tokenFromContext(ctx context.Context) (string, error) {
	auth, ok := ctx.Value(authKey{}).(string)
	if !ok {
		return "", fmt.Errorf("missing auth")
	}
	return auth, nil
}

// authFromRequest extracts the auth token from the request headers.
// Moved here from sse_auth.go.
func authFromRequest(ctx context.Context, r *http.Request) context.Context {
	return withAuthKey(ctx, r.Header.Get("Authorization"))
}
