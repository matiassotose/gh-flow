package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	checkMark         = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).SetString("✓")
	uncheckedMark     = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).SetString("○")
)

// RepoSelectionModel for selecting repositories
type RepoSelectionModel struct {
	repos         []string
	selected      map[int]bool
	cursor        int
	SelectedRepos []string
	quitting      bool
}

func NewRepoSelectionModel(repos []string) RepoSelectionModel {
	return RepoSelectionModel{
		repos:    repos,
		selected: make(map[int]bool),
		cursor:   0,
	}
}

func NewRepoSelectionModelWithDefaults(repos, defaults []string) RepoSelectionModel {
	selected := make(map[int]bool)
	for i, repo := range repos {
		for _, def := range defaults {
			if repo == def {
				selected[i] = true
				break
			}
		}
	}
	return RepoSelectionModel{
		repos:    repos,
		selected: selected,
		cursor:   0,
	}
}

func (m RepoSelectionModel) Init() tea.Cmd {
	return nil
}

func (m RepoSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.repos)-1 {
				m.cursor++
			}

		case " ", "enter":
			m.selected[m.cursor] = !m.selected[m.cursor]

		case "ctrl+d":
			// Done selecting
			m.SelectedRepos = make([]string, 0)
			for i, repo := range m.repos {
				if m.selected[i] {
					m.SelectedRepos = append(m.SelectedRepos, repo)
				}
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m RepoSelectionModel) View() string {
	if m.quitting && len(m.SelectedRepos) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("Select repositories (Space/Enter to toggle, Ctrl+D when done):"))
	b.WriteString("\n\n")

	for i, repo := range m.repos {
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		checked := uncheckedMark.String()
		if m.selected[i] {
			checked = checkMark.String()
		}

		item := fmt.Sprintf("%s%s %s", cursor, checked, repo)
		if m.cursor == i {
			b.WriteString(selectedItemStyle.Render(item))
		} else {
			b.WriteString(itemStyle.Render(item))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Press Ctrl+D to confirm selection"))

	return b.String()
}
