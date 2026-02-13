package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// UrgencyModel for selecting urgency level
type UrgencyModel struct {
	cursor   int
	choices  []string
	Selected string
}

func NewUrgencyModel() UrgencyModel {
	return UrgencyModel{
		cursor:  0,
		choices: []string{"normal", "urgent"},
	}
}

func (m UrgencyModel) Init() tea.Cmd {
	return nil
}

func (m UrgencyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			m.Selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m UrgencyModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Select urgency level:"))
	b.WriteString("\n\n")

	descriptions := []string{
		"Normal - Branch from dev (standard development)",
		"Urgent  - Branch from main (hotfix/production)",
	}

	for i := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		item := fmt.Sprintf("%s%s", cursor, descriptions[i])
		if m.cursor == i {
			b.WriteString(selectedItemStyle.Render(item))
		} else {
			b.WriteString(itemStyle.Render(item))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Press Enter to select"))

	return b.String()
}
