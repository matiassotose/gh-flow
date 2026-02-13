package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// StashModel for handling uncommitted changes
type StashModel struct {
	repos         []string
	selectedRepos map[int]bool
	cursor        int
	mode          string // "select", "action"
	ShouldStash   bool
	Aborted       bool
	SelectedRepos []string // Exported field with final selection
}

func NewStashModel(repos []string) StashModel {
	selected := make(map[int]bool)
	for i := range repos {
		selected[i] = true // Default: all selected
	}
	return StashModel{
		repos:         repos,
		selectedRepos: selected,
		cursor:        0,
		mode:          "select",
	}
}

func (m StashModel) Init() tea.Cmd {
	return nil
}

func (m StashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case "select":
			switch msg.String() {
			case "ctrl+c":
				m.Aborted = true
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < len(m.repos) {
					m.cursor++
				}

			case " ", "enter":
				if m.cursor < len(m.repos) {
					m.selectedRepos[m.cursor] = !m.selectedRepos[m.cursor]
				}

			case "a":
				// Select all
				for i := range m.repos {
					m.selectedRepos[i] = true
				}

			case "n":
				// Select none
				for i := range m.repos {
					m.selectedRepos[i] = false
				}

			case "ctrl+d":
				m.mode = "action"
				m.cursor = 0
			}

		case "action":
			switch msg.String() {
			case "ctrl+c", "q":
				m.Aborted = true
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < 2 {
					m.cursor++
				}

			case "enter":
				switch m.cursor {
				case 0: // Stash and continue
					m.ShouldStash = true
					// Fill SelectedRepos based on selection
					m.SelectedRepos = make([]string, 0)
					for i, repo := range m.repos {
						if m.selectedRepos[i] {
							m.SelectedRepos = append(m.SelectedRepos, repo)
						}
					}
					return m, tea.Quit
				case 1: // Continue without stashing
					m.ShouldStash = false
					return m, tea.Quit
				case 2: // Abort
					m.Aborted = true
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
}

func (m StashModel) View() string {
	var b strings.Builder

	switch m.mode {
	case "select":
		b.WriteString(titleStyle.Render("⚠️  Uncommitted changes detected"))
		b.WriteString("\n\n")
		b.WriteString(itemStyle.Render("Select repositories to stash:"))
		b.WriteString("\n\n")

		// Add "Select All" option
		cursor := "  "
		if m.cursor == len(m.repos) {
			cursor = "> "
		}
		b.WriteString(fmt.Sprintf("%s[Press 'a' for all, 'n' for none]\n\n", cursor))

		for i, repo := range m.repos {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}

			checked := uncheckedMark.String()
			if m.selectedRepos[i] {
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
		b.WriteString(itemStyle.Render("Press Ctrl+D to continue"))

	case "action":
		b.WriteString(titleStyle.Render("Choose action:"))
		b.WriteString("\n\n")

		actions := []string{
			"Stash changes and continue",
			"Continue without stashing (discard changes)",
			"Abort",
		}

		for i, action := range actions {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}

			item := cursor + action
			if m.cursor == i {
				b.WriteString(selectedItemStyle.Render(item))
			} else {
				b.WriteString(itemStyle.Render(item))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}
