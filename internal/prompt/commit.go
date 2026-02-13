package prompt

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CommitModel for entering commit message
type CommitModel struct {
	input   textinput.Model
	Message string
	Aborted bool
}

func NewCommitModel() CommitModel {
	ti := textinput.New()
	ti.Placeholder = "feat: add awesome feature"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 70

	return CommitModel{
		input: ti,
	}
}

func (m CommitModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CommitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.Aborted = true
			return m, tea.Quit

		case "enter":
			m.Message = m.input.Value()
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m CommitModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Enter commit message:"))
	b.WriteString("\n\n")
	b.WriteString(itemStyle.Render("This message will be used for all repositories"))
	b.WriteString("\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(itemStyle.Render("Press Enter to confirm"))

	return b.String()
}
