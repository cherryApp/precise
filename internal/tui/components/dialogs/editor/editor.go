package editor

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/crush/internal/tui/components/core"
	"github.com/charmbracelet/crush/internal/tui/components/dialogs"
	"github.com/charmbracelet/crush/internal/tui/exp/list"
	"github.com/charmbracelet/crush/internal/tui/styles"
	"github.com/charmbracelet/crush/internal/tui/util"
	"github.com/charmbracelet/lipgloss/v2"
)

const (
	EditorDialogID dialogs.DialogID = "editor"

	defaultWidth int = 70
)

// EditorOption represents an editor that can be selected
type EditorOption struct {
	Name        string
	Description string
	Command     string
}

// EditorSelectedMsg is sent when an editor is selected
type EditorSelectedMsg struct {
	Editor EditorOption
}

// EditorDialog interface for the editor selection dialog
type EditorDialog interface {
	dialogs.DialogModel
}

type editorDialogCmp struct {
	width   int
	wWidth  int // Width of the terminal window
	wHeight int // Height of the terminal window

	editorList list.FilterableList[list.CompletionItem[EditorOption]]
	keyMap     KeyMap
	help       help.Model
	editors    []EditorOption
}

func NewEditorDialog() EditorDialog {
	keyMap := DefaultEditorDialogKeyMap()

	listKeyMap := list.DefaultKeyMap()
	listKeyMap.Down.SetEnabled(false)
	listKeyMap.Up.SetEnabled(false)
	listKeyMap.DownOneItem = keyMap.Select
	listKeyMap.UpOneItem = key.NewBinding(key.WithKeys("up"))

	t := styles.CurrentTheme()
	inputStyle := t.S().Base.PaddingLeft(1).PaddingBottom(1)

	editors := discoverEditors()

	editorItems := []list.CompletionItem[EditorOption]{}
	for _, editor := range editors {
		editorItems = append(editorItems, list.NewCompletionItem(editor.Name, editor))
	}

	editorList := list.NewFilterableList(
		editorItems,
		list.WithFilterInputStyle(inputStyle),
		list.WithFilterListOptions(
			list.WithKeyMap(listKeyMap),
			list.WithWrapNavigation(),
			list.WithResizeByList(),
		),
	)

	helpModel := help.New()
	helpModel.Styles = t.S().Help

	return &editorDialogCmp{
		editorList: editorList,
		width:      defaultWidth,
		keyMap:     keyMap,
		help:       helpModel,
		editors:    editors,
	}
}

func (e *editorDialogCmp) Init() tea.Cmd {
	return e.editorList.SetSize(e.listWidth(), e.listHeight())
}

func (e *editorDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.wWidth = msg.Width
		e.wHeight = msg.Height
		return e, e.editorList.SetSize(e.listWidth(), e.listHeight())
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, e.keyMap.Select):
			selectedItem := e.editorList.SelectedItem()
			if selectedItem == nil {
				return e, nil // No item selected
			}
			editor := (*selectedItem).Value()
			return e, tea.Sequence(
				util.CmdHandler(dialogs.CloseDialogMsg{}),
				util.CmdHandler(EditorSelectedMsg{Editor: editor}),
			)
		case key.Matches(msg, e.keyMap.Close):
			return e, util.CmdHandler(dialogs.CloseDialogMsg{})
		default:
			u, cmd := e.editorList.Update(msg)
			e.editorList = u.(list.FilterableList[list.CompletionItem[EditorOption]])
			return e, cmd
		}
	}
	return e, nil
}

func (e *editorDialogCmp) View() string {
	t := styles.CurrentTheme()
	listView := e.editorList.View()

	header := t.S().Base.Padding(0, 1, 1, 1).Render(core.Title("Select External Editor", e.width-4))
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		listView,
		"",
		t.S().Base.Width(e.width-2).PaddingLeft(1).AlignHorizontal(lipgloss.Left).Render(e.help.View(e.keyMap)),
	)
	return e.style().Render(content)
}

func (e *editorDialogCmp) Cursor() *tea.Cursor {
	if cursor, ok := e.editorList.(util.Cursor); ok {
		cursor := cursor.Cursor()
		if cursor != nil {
			cursor = e.moveCursor(cursor)
		}
		return cursor
	}
	return nil
}

func (e *editorDialogCmp) listWidth() int {
	return defaultWidth - 2 // 4 for padding
}

func (e *editorDialogCmp) listHeight() int {
	return min(len(e.editors)+2, e.wHeight/2)
}

func (e *editorDialogCmp) moveCursor(cursor *tea.Cursor) *tea.Cursor {
	row, col := e.Position()
	offset := row + 3
	cursor.Y += offset
	cursor.X = cursor.X + col + 2
	return cursor
}

func (e *editorDialogCmp) style() lipgloss.Style {
	t := styles.CurrentTheme()
	return t.S().Base.
		Width(e.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus)
}

func (e *editorDialogCmp) Position() (int, int) {
	row := e.wHeight/4 - 2 // just a bit above the center
	col := e.wWidth / 2
	col -= e.width / 2
	return row, col
}

func (e *editorDialogCmp) ID() dialogs.DialogID {
	return EditorDialogID
}

// discoverEditors discovers available editors based on the platform
func discoverEditors() []EditorOption {
	var editors []EditorOption

	// Check if EDITOR environment variable is already set
	if editor := os.Getenv("EDITOR"); editor != "" {
		editors = append(editors, EditorOption{
			Name:        "Current EDITOR (" + editor + ")",
			Description: "Use the currently configured EDITOR",
			Command:     editor,
		})
	}

	// Platform-specific editors
	var candidates []string
	if runtime.GOOS == "windows" {
		candidates = []string{
			"code",      // VS Code
			"code-insiders",
			"notepad++",
			"sublime_text",
			"atom",
			"notepad",
			"notepad.exe",
		}
	} else {
		candidates = []string{
			"code",      // VS Code
			"code-insiders",
			"subl",      // Sublime Text
			"atom",      // Atom
			"nano",
			"vim",
			"vi",
			"emacs",
			"nvim",      // Neovim
			"hx",        // Helix
			"zed",       // Zed
		}
	}

	// Check which editors are available
	for _, candidate := range candidates {
		if _, err := exec.LookPath(candidate); err == nil {
			var name, description string
			switch candidate {
			case "code":
				name = "Visual Studio Code"
				description = "Popular code editor by Microsoft"
			case "code-insiders":
				name = "VS Code Insiders"
				description = "Insiders build of Visual Studio Code"
			case "subl":
				name = "Sublime Text"
				description = "Fast and lightweight text editor"
			case "atom":
				name = "Atom"
				description = "Hackable text editor by GitHub"
			case "nano":
				name = "GNU nano"
				description = "Simple terminal-based editor"
			case "vim":
				name = "Vim"
				description = "Highly configurable text editor"
			case "vi":
				name = "Vi"
				description = "Classic Unix text editor"
			case "emacs":
				name = "Emacs"
				description = "Extensible text editor"
			case "nvim":
				name = "Neovim"
				description = "Hyperextensible Vim-based editor"
			case "hx":
				name = "Helix"
				description = "Post-modern modal text editor"
			case "zed":
				name = "Zed"
				description = "High-performance code editor"
			case "notepad++":
				name = "Notepad++"
				description = "Free source code editor"
			case "sublime_text":
				name = "Sublime Text"
				description = "Sophisticated text editor"
			case "notepad", "notepad.exe":
				name = "Notepad"
				description = "Basic Windows text editor"
			default:
				name = candidate
				description = "Text editor"
			}
			editors = append(editors, EditorOption{
				Name:        name,
				Description: description,
				Command:     candidate,
			})
		}
	}

	// Add option to manually specify an editor
	editors = append(editors, EditorOption{
		Name:        "Custom Editor",
		Description: "Specify a custom editor command",
		Command:     "custom",
	})

	return editors
}