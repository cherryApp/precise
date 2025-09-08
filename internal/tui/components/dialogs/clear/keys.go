package clear

import "github.com/charmbracelet/bubbles/v2/key"

// KeyMap defines the keybindings for the clear context dialog.
type KeyMap struct {
	EnterSpace key.Binding
	LeftRight  key.Binding
	Tab        key.Binding
	Yes        key.Binding
	No         key.Binding
	Close      key.Binding
}

// DefaultKeymap returns the default keybindings for the clear context dialog.
func DefaultKeymap() KeyMap {
	return KeyMap{
		EnterSpace: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "confirm"),
		),
		LeftRight: key.NewBinding(
			key.WithKeys("left", "right"),
			key.WithHelp("←/→", "navigate"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "navigate"),
		),
		Yes: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "yes"),
		),
		No: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "no"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "close"),
		),
	}
}
