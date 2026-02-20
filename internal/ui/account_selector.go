// Package ui
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	fizzy "github.com/rogeriopvl/fizzy-go"
)

type accountModel struct {
	accounts []fizzy.Account
	cursor   int
}

func (m accountModel) Init() tea.Cmd {
	return nil
}

func (m accountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.accounts)-1 {
				m.cursor++
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m accountModel) View() string {
	s := "Select an account:\n\n"
	for i, account := range m.accounts {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		s += fmt.Sprintf("%s%s (%s)\n", cursor, account.Name, account.Slug)
	}
	s += "\nUse ↑/↓ or k/j to navigate, Enter to select, q to quit"
	return s
}

func SelectAccount(accounts []fizzy.Account) (fizzy.Account, error) {
	m := accountModel{accounts: accounts, cursor: 0}
	p, err := tea.NewProgram(m).Run()
	if err != nil {
		return fizzy.Account{}, err
	}

	finalModel := p.(accountModel)
	return finalModel.accounts[finalModel.cursor], nil
}
