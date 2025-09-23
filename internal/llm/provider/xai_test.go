package provider

import (
	"testing"

	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestXAIToolCallParsing(t *testing.T) {
	// Create a mock OpenAI client with x.ai configuration
	opts := providerClientOptions{
		baseURL: "https://api.x.ai/v1",
		config: config.ProviderConfig{
			ID: "xai",
		},
		model: func(config.SelectedModelType) (m catwalk.Model) {
			return catwalk.Model{ID: "grok-2"}
		},
	}

	client := &openaiClient{
		providerOptions: opts,
	}

	t.Run("should detect x.ai model correctly", func(t *testing.T) {
		assert.True(t, client.isXAIModel())
	})

	t.Run("should parse basic x.ai tool call", func(t *testing.T) {
		content := `<xai:function_call name="multiedit"> [{"old_string": "test", "new_string": "updated"}]`

		toolCalls := client.parseXAIToolCalls(content)

		assert.Len(t, toolCalls, 1)
		assert.Equal(t, "multiedit", toolCalls[0].Name)
		assert.Equal(t, "xai_call_0", toolCalls[0].ID)
		assert.Equal(t, "function", toolCalls[0].Type)
		assert.True(t, toolCalls[0].Finished)
	})

	t.Run("should parse x.ai tool call with incomplete array", func(t *testing.T) {
		content := `<xai:function_call name="multiedit"> [{"old_string": "test", "new_string": "updated"}`

		toolCalls := client.parseXAIToolCalls(content)

		assert.Len(t, toolCalls, 1)
		assert.Equal(t, "multiedit", toolCalls[0].Name)
	})

	t.Run("should parse multiple x.ai tool calls", func(t *testing.T) {
		content := `<xai:function_call name="view"> {"file_path": "/test"} </xai:function_call>
		<xai:function_call name="edit"> {"file_path": "/test", "old_string": "a", "new_string": "b"}`

		toolCalls := client.parseXAIToolCalls(content)

		assert.Len(t, toolCalls, 2)
		assert.Equal(t, "view", toolCalls[0].Name)
		assert.Equal(t, "edit", toolCalls[1].Name)
	})

	t.Run("should handle non-JSON arguments", func(t *testing.T) {
		content := `<xai:function_call name="bash"> ls -la`

		toolCalls := client.parseXAIToolCalls(content)

		assert.Len(t, toolCalls, 1)
		assert.Equal(t, "bash", toolCalls[0].Name)
		assert.Equal(t, `"ls -la"`, toolCalls[0].Input)
	})
}

func TestMultiEditXAIArgsParsing(t *testing.T) {
	client := &openaiClient{}

	t.Run("should parse valid JSON array", func(t *testing.T) {
		rawArgs := `[{"old_string": "test", "new_string": "updated"}]`

		result := client.parseMultiEditXAIArgs(rawArgs)

		assert.Contains(t, result, "file_path")
		assert.Contains(t, result, "edits")
		assert.Contains(t, result, "old_string")
		assert.Contains(t, result, "new_string")
	})

	t.Run("should handle incomplete JSON array", func(t *testing.T) {
		rawArgs := `[{"old_string": "test", "new_string": "updated"`

		result := client.parseMultiEditXAIArgs(rawArgs)

		// Should return the original if it can't be fixed
		assert.Equal(t, rawArgs, result)
	})

	t.Run("should return as-is for non-array format", func(t *testing.T) {
		rawArgs := `{"file_path": "/test", "edits": []}`

		result := client.parseMultiEditXAIArgs(rawArgs)

		assert.Equal(t, rawArgs, result)
	})
}