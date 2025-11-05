package debug

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
	"github.com/charmbracelet/lipgloss/v2"
)

// Panel represents the debug panel component
type Panel interface {
	util.Model
	SetSize(width, height int)
	GetDebugInfo() *DebugInfo
	SetTab(tab TabType)
	GetCurrentTab() TabType
}

type TabType string

const (
	TabMetrics TabType = "metrics"
	TabLogs    TabType = "logs"
)

type panel struct {
	width, height int
	debugInfo     *DebugInfo
	currentTab    TabType
	maxLogEntries int
}

// NewPanel creates a new debug panel
func NewPanel() Panel {
	return &panel{
		debugInfo:     NewDebugInfo(100), // Keep last 100 log entries
		currentTab:    TabMetrics,
		maxLogEntries: 15, // Show up to 15 log entries in the UI
	}
}

func (p *panel) Init() tea.Cmd {
	return nil
}

func (p *panel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle tab switching if needed in the future
	return p, nil
}

func (p *panel) View() string {
	if p.width == 0 || p.height == 0 {
		return ""
	}

	t := styles.CurrentTheme()

	// Create tab headers
	tabStyle := t.S().Base.
		Padding(0, 2).
		Margin(0, 1, 0, 0)

	activeTabStyle := tabStyle.Copy().
		Background(t.Primary).
		Foreground(t.BgBase).
		Bold(true)

	inactiveTabStyle := tabStyle.Copy().
		Background(t.FgMuted).
		Foreground(t.BgBase)

	var metricsTab, logsTab string
	if p.currentTab == TabMetrics {
		metricsTab = activeTabStyle.Render("Metrics")
		logsTab = inactiveTabStyle.Render("Logs")
	} else {
		metricsTab = inactiveTabStyle.Render("Metrics")
		logsTab = activeTabStyle.Render("Logs")
	}

	// Add key binding hint
	keyHintStyle := t.S().Base.
		Foreground(t.FgMuted).
		Padding(0, 1).
		Margin(0, 0, 0, 2)
	keyHint := keyHintStyle.Render("shift+tab")

	tabs := lipgloss.JoinHorizontal(lipgloss.Left, metricsTab, logsTab, keyHint)

	// Create content based on current tab
	contentHeight := p.height - 3 // Account for tabs and borders
	contentWidth := p.width - 4   // Account for padding

	var content string
	switch p.currentTab {
	case TabMetrics:
		content = p.debugInfo.FormatMetrics()
	case TabLogs:
		content = p.debugInfo.FormatLogs(p.maxLogEntries)
	}

	// Wrap content and ensure it fits in the available space
	contentStyle := t.S().Base.
		Width(contentWidth).
		Height(contentHeight).
		Padding(1)

	wrappedContent := contentStyle.Render(content)

	// Combine tabs and content
	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabs,
		t.S().Base.Foreground(t.Border).Render(strings.Repeat("─", p.width-2)),
		wrappedContent,
	)
}

func (p *panel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *panel) GetDebugInfo() *DebugInfo {
	return p.debugInfo
}

func (p *panel) SetTab(tab TabType) {
	p.currentTab = tab
}

func (p *panel) GetCurrentTab() TabType {
	return p.currentTab
}