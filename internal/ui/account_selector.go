// Package ui
package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rogeriopvl/fizzy/internal/api"
)

type accountModel struct {
	accounts []api.Account
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

func SelectAccount(accounts []api.Account) (api.Account, error) {
	m := accountModel{accounts: accounts, cursor: 0}
	p, err := tea.NewProgram(m).Run()
	if err != nil {
		return api.Account{}, err
	}

	finalModel := p.(accountModel)
	return finalModel.accounts[finalModel.cursor], nil
}
