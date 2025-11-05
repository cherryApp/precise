package debug

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/crush/internal/llm/provider"
)

// ExecutionMetrics tracks the performance metrics of the last LLM execution
type ExecutionMetrics struct {
	// Timing information
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration

	// Token information
	InputTokens  int64
	OutputTokens int64
	CacheCreationTokens int64
	CacheReadTokens int64

	// Performance metrics
	TokensPerSecond float64

	// Model information
	ModelName string
	Provider  string

	// Request metadata
	RequestID string
	FinishReason string

	// Cost information
	Cost float64
}

// LogEntry represents a single log entry for the debug panel
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Source    string
	Message   string
	Data      map[string]interface{}
}

// DebugInfo holds all debug information for display
type DebugInfo struct {
	LastExecution *ExecutionMetrics
	Logs          []LogEntry
	MaxLogs       int
}

// NewDebugInfo creates a new DebugInfo instance
func NewDebugInfo(maxLogs int) *DebugInfo {
	return &DebugInfo{
		Logs:    make([]LogEntry, 0, maxLogs),
		MaxLogs: maxLogs,
	}
}

// UpdateExecutionMetrics updates the last execution metrics
func (d *DebugInfo) UpdateExecutionMetrics(
	startTime, endTime time.Time,
	usage provider.TokenUsage,
	modelName, provider, finishReason string,
	cost float64,
) {
	duration := endTime.Sub(startTime)
	totalOutputTokens := usage.OutputTokens + usage.CacheReadTokens

	// Calculate tokens per second (only for output tokens)
	var tps float64
	if totalOutputTokens > 0 && duration.Seconds() > 0 {
		tps = float64(totalOutputTokens) / duration.Seconds()
	}

	d.LastExecution = &ExecutionMetrics{
		StartTime:           startTime,
		EndTime:             endTime,
		Duration:            duration,
		InputTokens:         usage.InputTokens,
		OutputTokens:        usage.OutputTokens,
		CacheCreationTokens: usage.CacheCreationTokens,
		CacheReadTokens:     usage.CacheReadTokens,
		TokensPerSecond:     tps,
		ModelName:           modelName,
		Provider:            provider,
		RequestID:           fmt.Sprintf("req_%d", endTime.Unix()),
		FinishReason:        finishReason,
		Cost:                cost,
	}

	// Add a log entry for this execution
	d.AddLog("INFO", "EXECUTION", fmt.Sprintf("Completed request in %v (%.2f TPS)",
		duration, tps), map[string]interface{}{
		"duration_ms":     duration.Milliseconds(),
		"input_tokens":    usage.InputTokens,
		"output_tokens":   totalOutputTokens,
		"tokens_per_sec":  tps,
		"model":          modelName,
		"provider":       provider,
		"finish_reason":  finishReason,
	})
}

// AddLog adds a new log entry, maintaining the maximum number of logs
func (d *DebugInfo) AddLog(level, source, message string, data map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Source:    source,
		Message:   message,
		Data:      data,
	}

	d.Logs = append(d.Logs, entry)

	// Keep only the latest MaxLogs entries
	if len(d.Logs) > d.MaxLogs {
		d.Logs = d.Logs[len(d.Logs)-d.MaxLogs:]
	}
}

// FormatMetrics returns a formatted string of the last execution metrics
func (d *DebugInfo) FormatMetrics() string {
	if d.LastExecution == nil {
		return "No execution data available"
	}

	m := d.LastExecution
	var builder strings.Builder

	// Header
	builder.WriteString("🚀 Last Execution Metrics\n\n")

	// Timing
	builder.WriteString(fmt.Sprintf("⏱️  Duration: %v\n", m.Duration))
	builder.WriteString(fmt.Sprintf("📊 TPS: %.2f tokens/sec\n\n", m.TokensPerSecond))

	// Tokens
	totalInput := m.InputTokens + m.CacheCreationTokens
	totalOutput := m.OutputTokens + m.CacheReadTokens
	builder.WriteString(fmt.Sprintf("📥 Input: %d tokens", totalInput))
	if m.CacheCreationTokens > 0 {
		builder.WriteString(fmt.Sprintf(" (+%d cached)", m.CacheCreationTokens))
	}
	builder.WriteString("\n")

	builder.WriteString(fmt.Sprintf("📤 Output: %d tokens", totalOutput))
	if m.CacheReadTokens > 0 {
		builder.WriteString(fmt.Sprintf(" (+%d cached)", m.CacheReadTokens))
	}
	builder.WriteString("\n\n")

	// Model info
	builder.WriteString(fmt.Sprintf("🤖 Model: %s\n", m.ModelName))
	builder.WriteString(fmt.Sprintf("🏢 Provider: %s\n", m.Provider))
	builder.WriteString(fmt.Sprintf("🏁 Finish: %s\n", m.FinishReason))

	// Cost
	if m.Cost > 0 {
		builder.WriteString(fmt.Sprintf("💰 Cost: $%.6f\n", m.Cost))
	}

	return builder.String()
}

// FormatLogs returns a formatted string of recent logs
func (d *DebugInfo) FormatLogs(maxEntries int) string {
	if len(d.Logs) == 0 {
		return "No logs available"
	}

	var builder strings.Builder
	builder.WriteString("📋 Execution Logs\n\n")

	// Show the most recent logs (up to maxEntries)
	start := 0
	if len(d.Logs) > maxEntries {
		start = len(d.Logs) - maxEntries
	}

	for i := start; i < len(d.Logs); i++ {
		entry := d.Logs[i]
		timeStr := entry.Timestamp.Format("15:04:05")

		// Format based on log level
		levelIcon := "ℹ️"
		switch entry.Level {
		case "ERROR":
			levelIcon = "❌"
		case "WARN":
			levelIcon = "⚠️"
		case "DEBUG":
			levelIcon = "🔍"
		}

		builder.WriteString(fmt.Sprintf("%s [%s] %s: %s\n",
			levelIcon, timeStr, entry.Source, entry.Message))
	}

	return builder.String()
}