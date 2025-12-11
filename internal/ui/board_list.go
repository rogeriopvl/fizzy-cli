package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rogeriopvl/fizzy-cli/internal/api"
)

type boardListModel struct {
	boards []api.Board
	cursor int
}

func (m boardListModel) Init() tea.Cmd {
	return nil
}

func (m boardListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.boards)-1 {
				m.cursor++
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m boardListModel) View() string {
	s := "Boards:\n\n"
	for i, board := range m.boards {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		s += fmt.Sprintf("%s%s\n", cursor, board.Name)
	}
	s += "\nUse ↑/↓ or k/j to navigate, Enter to select, q to quit"
	return s
}

// DisplayBoards shows an interactive list of boards.
func DisplayBoards(boards []api.Board) error {
	m := boardListModel{boards: boards, cursor: 0}
	_, err := tea.NewProgram(m).Run()
	return err
}
