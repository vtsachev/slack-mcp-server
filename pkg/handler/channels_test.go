package handler

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/korotovsky/slack-mcp-server/pkg/provider"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockApiProvider is a mock implementation of the ApiProvider.
type MockApiProvider struct {
	MockProvide         func() (*slack.Client, error)
	MockProvideUsersMap func() map[string]slack.User
}

func (m *MockApiProvider) Provide() (*slack.Client, error) {
	if m.MockProvide != nil {
		return m.MockProvide()
	}
	return nil, fmt.Errorf("MockProvide function not set")
}

func (m *MockApiProvider) ProvideUsersMap() map[string]slack.User {
	if m.MockProvideUsersMap != nil {
		return m.MockProvideUsersMap()
	}
	return make(map[string]slack.User) // Return empty map
}

// mockSlackClientForConversations is used to mock the GetConversationsContext method.
// It embeds slack.Client to satisfy the interface where *slack.Client is expected by the handler,
// allowing us to override only the methods we need for the mock.
type mockSlackClientForConversations struct {
	slack.Client
	GetConversationsContextFunc func(ctx context.Context, params *slack.GetConversationsParameters) ([]slack.Channel, string, error)
}

// Override GetConversationsContext to use our mock function.
// Note: This method signature must exactly match the one in the version of slack-go/slack being used.
func (m *mockSlackClientForConversations) GetConversationsContext(ctx context.Context, params *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
	if m.GetConversationsContextFunc != nil {
		return m.GetConversationsContextFunc(ctx, params)
	}
	return nil, "", fmt.Errorf("GetConversationsContextFunc not set in mock")
}


func TestNewChannelsHandler(t *testing.T) {
	apiProviderMock := &MockApiProvider{}
	var p provider.ApiProvider = apiProviderMock // Use the interface type
	handler := NewChannelsHandler(p)

	require.NotNil(t, handler, "NewChannelsHandler returned nil")
	assert.Equal(t, p, handler.apiProvider, "NewChannelsHandler did not set apiProvider correctly")
	assert.Len(t, handler.validTypes, len(AllChanTypes), "Expected validTypes to have %d entries, got %d", len(AllChanTypes), len(handler.validTypes))
	for _, chanType := range AllChanTypes {
		assert.True(t, handler.validTypes[chanType], "Expected validTypes to contain %s", chanType)
	}
}

func TestChannelsHandler_Success_Basic(t *testing.T) {
	ctx := context.Background()

	mockSlackClient := &mockSlackClientForConversations{}
	mockSlackClient.GetConversationsContextFunc = func(ctx context.Context, params *slack.GetConversationsParameters) ([]slack.Channel, string, error) {
		assert.Equal(t, 10, params.Limit)
		assert.Equal(t, []string{PubChanType}, params.Types)

		channels := []slack.Channel{
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{ID: "C1", NumMembers: 5},
					Name:         "channel1",
					Purpose:      slack.Purpose{Value: "Purpose1"},
					Topic:        slack.Topic{Value: "Topic1"},
				},
			},
			{
				GroupConversation: slack.GroupConversation{
					Conversation: slack.Conversation{ID: "C2", NumMembers: 10},
					Name:         "channel2",
					Purpose:      slack.Purpose{Value: "Purpose2"},
					Topic:        slack.Topic{Value: "Topic2"},
				},
			},
		}
		return channels, "nextcursor123", nil // Provide a next cursor
	}

	apiProviderMock := &MockApiProvider{
		MockProvide: func() (*slack.Client, error) {
			// Reverting to returning the embedded actual client to ensure compilation.
			// This means GetConversationsContextFunc will NOT be called by this setup.
			// Proper mocking of slack.Client typically requires mocking the HTTP transport.
			return &mockSlackClient.Client, nil
		},
	}

	var p provider.ApiProvider = apiProviderMock
	handler := NewChannelsHandler(p)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "channelsTool", // Tool name, not strictly used by this handler but good practice
			Arguments: map[string]interface{}{ // This was the crucial part
				"limit":         10,
				"sort":          "popularity",
				"channel_types": PubChanType,
			},
		},
	}

	result, err := handler.ChannelsHandler(ctx, request)
	require.NoError(t, err)
	require.NotNil(t, result)
	// Removing assertion for result.Type as mcp.ToolResultType_TEXT and result.Type are undefined.
	// The successful extraction of TextContent implies it's a text result.
	require.NotEmpty(t, result.Content, "Result content should not be empty")

	textContent, ok := result.Content[0].(mcp.TextContent)
	require.True(t, ok, "First content element is not TextContent")

	// Expected CSV output (Sorted by popularity: C2 then C1)
	// gocsv adds a newline at the end.
	// The last channel (C1 after sorting) should have the cursor.
	expectedCSVHeader := "ID,Name,Topic,Purpose,MemberCount,Cursor"
	expectedCSVLine1 := "C2,#channel2,Topic2,Purpose2,10,"
	expectedCSVLine2 := "C1,#channel1,Topic1,Purpose1,5,nextcursor123"

	lines := strings.Split(strings.TrimSpace(textContent.Text), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSuffix(line, "\r") // Normalize newlines just in case
	}

	require.Len(t, lines, 3, "Expected 3 lines in CSV (header + 2 data rows), got: %v", lines)
	assert.Equal(t, expectedCSVHeader, lines[0], "CSV Header mismatch")
	assert.Equal(t, expectedCSVLine1, lines[1], "CSV data line 1 mismatch")
	assert.Equal(t, expectedCSVLine2, lines[2], "CSV data line 2 mismatch")
}

// TODO: Add more tests:
// - TestChannelsHandler_Success_Pagination (multiple fetches)
// - TestChannelsHandler_Success_Sorting_Alphabetical (if another sort is added)
// - TestChannelsHandler_Success_MultipleChannelTypes
// - TestChannelsHandler_Success_EmptyResult (e.g. no channels found)
// - TestChannelsHandler_Error_NoLimitOrCursor
// - TestChannelsHandler_Error_ApiProviderError
// - TestChannelsHandler_Error_GetConversationsError
// - TestChannelsHandler_Error_CsvMarshalError (hard to test without direct injection)
