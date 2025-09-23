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

func TestGLMToolCallParsing(t *testing.T) {
	// Create a mock OpenAI client with GLM configuration
	opts := providerClientOptions{
		baseURL: "https://api.glm.ai/v1",
		config: config.ProviderConfig{
			ID: "glm",
		},
		model: func(config.SelectedModelType) (m catwalk.Model) {
			return catwalk.Model{ID: "glm-4"}
		},
	}

	client := &openaiClient{
		providerOptions: opts,
	}

	t.Run("should detect GLM model correctly", func(t *testing.T) {
		assert.True(t, client.isGLMModel())
	})

	t.Run("should parse basic GLM tool call", func(t *testing.T) {
		content := `<think>
I need to view a file to understand its contents.
</think>

<tool_call>
view
<arg_key>file_path</arg_key>
<arg_value>/Volumes/Work/Dev/totaltel-manager/COMPONENTS.md</arg_value>
</tool_call>`

		toolCalls := client.parseGLMToolCalls(content)

		assert.Len(t, toolCalls, 1)
		assert.Equal(t, "view", toolCalls[0].Name)
		assert.Equal(t, "glm_call_0", toolCalls[0].ID)
		assert.Equal(t, "function", toolCalls[0].Type)
		assert.True(t, toolCalls[0].Finished)
		assert.Contains(t, toolCalls[0].Input, "file_path")
		assert.Contains(t, toolCalls[0].Input, "/Volumes/Work/Dev/totaltel-manager/COMPONENTS.md")
	})

	t.Run("should parse GLM tool call with multiple arguments", func(t *testing.T) {
		content := `<tool_call>
edit
<arg_key>file_path</arg_key>
<arg_value>/test/file.txt</arg_value>
<arg_key>old_string</arg_key>
<arg_value>old content</arg_value>
<arg_key>new_string</arg_key>
<arg_value>new content</arg_value>
</tool_call>`

		toolCalls := client.parseGLMToolCalls(content)

		assert.Len(t, toolCalls, 1)
		assert.Equal(t, "edit", toolCalls[0].Name)
		assert.Contains(t, toolCalls[0].Input, "file_path")
		assert.Contains(t, toolCalls[0].Input, "old_string")
		assert.Contains(t, toolCalls[0].Input, "new_string")
		assert.Contains(t, toolCalls[0].Input, "/test/file.txt")
		assert.Contains(t, toolCalls[0].Input, "old content")
		assert.Contains(t, toolCalls[0].Input, "new content")
	})

	t.Run("should parse multiple GLM tool calls", func(t *testing.T) {
		content := `<think>
I need to view a file and then edit it.
</think>

<tool_call>
view
<arg_key>file_path</arg_key>
<arg_value>/test1.txt</arg_value>
</tool_call>

<tool_call>
edit
<arg_key>file_path</arg_key>
<arg_value>/test2.txt</arg_value>
<arg_key>old_string</arg_key>
<arg_value>test</arg_value>
<arg_key>new_string</arg_key>
<arg_value>updated</arg_value>
</tool_call>`

		toolCalls := client.parseGLMToolCalls(content)

		assert.Len(t, toolCalls, 2)
		assert.Equal(t, "view", toolCalls[0].Name)
		assert.Equal(t, "edit", toolCalls[1].Name)
		assert.Equal(t, "glm_call_0", toolCalls[0].ID)
		assert.Equal(t, "glm_call_1", toolCalls[1].ID)
	})

	t.Run("should handle GLM tool call with JSON value", func(t *testing.T) {
		content := `<tool_call>
config
<arg_key>settings</arg_key>
<arg_value>{"debug": true, "port": 8080}</arg_value>
</tool_call>`

		toolCalls := client.parseGLMToolCalls(content)

		assert.Len(t, toolCalls, 1)
		assert.Equal(t, "config", toolCalls[0].Name)
		assert.Contains(t, toolCalls[0].Input, "settings")
		assert.Contains(t, toolCalls[0].Input, "debug")
		assert.Contains(t, toolCalls[0].Input, "port")
	})

	t.Run("should handle GLM tool call with only think block", func(t *testing.T) {
		content := `<think>
Just thinking about the problem, no tool calls needed.
</think>`

		toolCalls := client.parseGLMToolCalls(content)

		assert.Len(t, toolCalls, 0)
	})

	t.Run("should handle empty GLM tool call", func(t *testing.T) {
		content := `<tool_call>
</tool_call>`

		toolCalls := client.parseGLMToolCalls(content)

		assert.Len(t, toolCalls, 0)
	})

	t.Run("should handle GLM tool call without args", func(t *testing.T) {
		content := `<tool_call>
help
</tool_call>`

		toolCalls := client.parseGLMToolCalls(content)

		assert.Len(t, toolCalls, 1)
		assert.Equal(t, "help", toolCalls[0].Name)
		assert.Equal(t, "{}", toolCalls[0].Input)
	})
}