package util

import (
	"log/slog"
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type Cursor interface {
	Cursor() *tea.Cursor
}

type Model interface {
	tea.Model
	tea.ViewModel
}

func CmdHandler(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func ReportError(err error) tea.Cmd {
	slog.Error("Error reported", "error", err)
	return CmdHandler(InfoMsg{
		Type: InfoTypeError,
		Msg:  err.Error(),
	})
}

type InfoType int

const (
	InfoTypeInfo InfoType = iota
	InfoTypeWarn
	InfoTypeError
)

func ReportInfo(info string) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeInfo,
		Msg:  info,
	})
}

func ReportWarn(warn string) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeWarn,
		Msg:  warn,
	})
}

type (
	InfoMsg struct {
		Type InfoType
		Msg  string
		TTL  time.Duration
	}
	ClearStatusMsg struct{}
	ChatFocusedMsg struct {
		Focused bool
	}
	ReloadLastPromptMsg  struct{}
	SwitchSessionsMsg    struct{}
	NewSessionsMsg       struct{}
	SwitchModelMsg       struct{}
	QuitMsg              struct{}
	OpenFilePickerMsg    struct{}
	ToggleHelpMsg        struct{}
	ToggleCompactModeMsg struct{}
	ToggleThinkingMsg    struct{}
	ToggleYoloModeMsg    struct{}
	CompactMsg           struct {
		SessionID string
	}
	ClearContextMsg struct {
		SessionID string
	}
	CommandRunCustomMsg struct {
		Content string
	}
	ShowArgumentsDialogMsg struct {
		CommandID string
		Content   string
		ArgNames  []string
	}
	CloseArgumentsDialogMsg struct {
		Submit    bool
		CommandID string
		Content   string
		Args      map[string]string
	}
	CancelTimerExpiredMsg struct{}
	ExecutionStartMsg struct {
		SessionID string
	}
)

func Clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}
