package provider

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/catwalk/pkg/catwalk"
	"github.com/charmbracelet/crush/internal/config"
)

func TestSystemPromptPrefixPath(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a test system prompt file
	systemPromptContent := "You are a helpful AI assistant with expertise in Go programming."
	systemPromptPath := filepath.Join(tempDir, "system_prompt.txt")

	err := os.WriteFile(systemPromptPath, []byte(systemPromptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test system prompt file: %v", err)
	}

	// Create a test configuration
	cfg := config.ProviderConfig{
		ID:                     "test-openai",
		Name:                   "Test OpenAI",
		Type:                   "openai",
		BaseURL:                "https://api.openai.com/v1",
		APIKey:                 "test-key",
		SystemPromptPrefixPath: systemPromptPath,
		Models: []catwalk.Model{
			{
				ID:   "gpt-4",
				Name: "GPT-4",
			},
		},
	}

	// Test creating a provider with the system prompt path
	prov, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// The provider should have been created successfully
	if prov == nil {
		t.Fatal("Provider is nil")
	}

	t.Logf("Provider created successfully with system prompt from file: %s", systemPromptPath)
}