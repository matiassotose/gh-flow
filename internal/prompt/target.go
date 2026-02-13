package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// TargetBranchModel for selecting target branches when auto-detection fails
type TargetBranchModel struct {
	repo             string
	cursor           int
	selected         map[int]bool
	choices          []string
	SelectedBranches []string
}

func NewTargetBranchModel(repo string) TargetBranchModel {
	return TargetBranchModel{
		repo:     repo,
		cursor:   0,
		selected: make(map[int]bool),
		choices:  []string{"dev", "main"},
	}
}

func (m TargetBranchModel) Init() tea.Cmd {
	return nil
}

func (m TargetBranchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		case " ", "enter":
			m.selected[m.cursor] = !m.selected[m.cursor]

		case "ctrl+d":
			// Done selecting
			m.SelectedBranches = make([]string, 0)
			for i, branch := range m.choices {
				if m.selected[i] {
					m.SelectedBranches = append(m.SelectedBranches, branch)
				}
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m TargetBranchModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("Select target branches for %s:", m.repo)))
	b.WriteString("\n\n")
	b.WriteString(itemStyle.Render("These branches will receive PRs"))
	b.WriteString("\n\n")

	for i, branch := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		checked := uncheckedMark.String()
		if m.selected[i] {
			checked = checkMark.String()
		}

		item := fmt.Sprintf("%s%s %s", cursor, checked, branch)
		if m.cursor == i {
			b.WriteString(selectedItemStyle.Render(item))
		} else {
			b.WriteString(itemStyle.Render(item))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Press Ctrl+D to confirm"))

	return b.String()
}
