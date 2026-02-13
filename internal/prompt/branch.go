package prompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// BranchModel for configuring the new branch
type BranchModel struct {
	branchInput textinput.Model
	cursor      int
	choices     []string
	BranchType  string
	BranchName  string
	Aborted     bool
	step        int // 0 = type, 1 = name
}

func NewBranchModel() BranchModel {
	ti := textinput.New()
	ti.Placeholder = "enter-branch-name"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return BranchModel{
		branchInput: ti,
		cursor:      0,
		choices:     []string{"feat", "hotfix"},
		step:        0,
	}
}

func (m BranchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m BranchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case 0: // Selecting type
			switch msg.String() {
			case "ctrl+c", "q":
				m.Aborted = true
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
				m.BranchType = m.choices[m.cursor]
				m.step = 1
				return m, nil
			}

		case 1: // Entering name
			switch msg.String() {
			case "ctrl+c":
				m.Aborted = true
				return m, tea.Quit

			case "esc":
				m.step = 0
				return m, nil

			case "enter":
				m.BranchName = m.branchInput.Value()
				if m.BranchName == "" {
					return m, nil // Don't proceed with empty name
				}
				return m, tea.Quit
			}

			m.branchInput, cmd = m.branchInput.Update(msg)
			return m, cmd
		}
	}

	return m, cmd
}

func (m BranchModel) View() string {
	var b strings.Builder

	switch m.step {
	case 0:
		b.WriteString(titleStyle.Render("Select branch type:"))
		b.WriteString("\n\n")

		for i, choice := range m.choices {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}

			desc := ""
			switch choice {
			case "feat":
				desc = "Feature - New functionality"
			case "hotfix":
				desc = "Hotfix - Critical bug fix"
			}

			item := fmt.Sprintf("%s%s", cursor, desc)
			if m.cursor == i {
				b.WriteString(selectedItemStyle.Render(item))
			} else {
				b.WriteString(itemStyle.Render(item))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(itemStyle.Render("Press Enter to continue"))

	case 1:
		b.WriteString(titleStyle.Render(fmt.Sprintf("Enter branch name (%s/):", m.BranchType)))
		b.WriteString("\n\n")
		b.WriteString(itemStyle.Render("Name: "))
		b.WriteString(m.branchInput.View())
		b.WriteString("\n\n")
		b.WriteString(itemStyle.Render("Press Enter to confirm, Esc to go back"))
	}

	return b.String()
}
