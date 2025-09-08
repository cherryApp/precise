package clear

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	question                              = "Are you sure you want to clear the current context?"
	ClearContextDialogID dialogs.DialogID = "clear_context"
)

// ClearContextDialog represents a confirmation dialog for clearing context.
type ClearContextDialog interface {
	dialogs.DialogModel
}

type clearContextDialogCmp struct {
	wWidth  int
	wHeight int

	selectedNo bool // true if "No" button is selected
	keymap     KeyMap
	sessionID  string
}

// NewClearContextDialog creates a new clear context confirmation dialog.
func NewClearContextDialog(sessionID string) ClearContextDialog {
	return &clearContextDialogCmp{
		selectedNo: true, // Default to "No" for safety
		keymap:     DefaultKeymap(),
		sessionID:  sessionID,
	}
}

func (c *clearContextDialogCmp) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input for the clear context dialog.
func (c *clearContextDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.wWidth = msg.Width
		c.wHeight = msg.Height
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, c.keymap.LeftRight, c.keymap.Tab):
			c.selectedNo = !c.selectedNo
			return c, nil
		case key.Matches(msg, c.keymap.EnterSpace):
			if !c.selectedNo {
				return c, tea.Sequence(
					util.CmdHandler(dialogs.CloseDialogMsg{}),
					util.CmdHandler(util.ClearContextMsg{SessionID: c.sessionID}),
				)
			}
			return c, util.CmdHandler(dialogs.CloseDialogMsg{})
		case key.Matches(msg, c.keymap.Yes):
			return c, tea.Sequence(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				util.CmdHandler(util.ClearContextMsg{SessionID: c.sessionID}),
			)
		case key.Matches(msg, c.keymap.No, c.keymap.Close):
			return c, util.CmdHandler(dialogs.CloseDialogMsg{})
		}
	}
	return c, nil
}

// View renders the clear context dialog with Yes/No buttons.
func (c *clearContextDialogCmp) View() string {
	t := styles.CurrentTheme()
	baseStyle := t.S().Base
	yesStyle := t.S().Text
	noStyle := yesStyle

	if c.selectedNo {
		noStyle = noStyle.Foreground(t.White).Background(t.Secondary)
		yesStyle = yesStyle.Background(t.BgSubtle)
	} else {
		yesStyle = yesStyle.Foreground(t.White).Background(t.Secondary)
		noStyle = noStyle.Background(t.BgSubtle)
	}

	const horizontalPadding = 3
	yesButton := yesStyle.PaddingLeft(horizontalPadding).Underline(true).Render("Y") +
		yesStyle.PaddingRight(horizontalPadding).Render("es")
	noButton := noStyle.PaddingLeft(horizontalPadding).Underline(true).Render("N") +
		noStyle.PaddingRight(horizontalPadding).Render("o")

	buttons := baseStyle.Width(lipgloss.Width(question)).Align(lipgloss.Right).Render(
		lipgloss.JoinHorizontal(lipgloss.Center, yesButton, "  ", noButton),
	)

	content := baseStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			question,
			"",
			buttons,
		),
	)

	clearDialogStyle := baseStyle.
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus)

	return clearDialogStyle.Render(content)
}

func (c *clearContextDialogCmp) Position() (int, int) {
	row := c.wHeight / 2
	row -= 7 / 2
	col := c.wWidth / 2
	col -= (lipgloss.Width(question) + 4) / 2

	return row, col
}

func (c *clearContextDialogCmp) ID() dialogs.DialogID {
	return ClearContextDialogID
}
