package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	choices  []branch
	cursor   int
	selected map[int]struct{}
}

type branch struct {
	name    string
	current bool
	remote  bool
}

type getBranchesMsg struct {
	branches []branch
}

func initialModel() model {
	return model{
		choices:  make([]branch, 0),
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

func getBranchesCmd() tea.Msg {
	branches := getBranches()
	ev := getBranchesMsg{branches}
	return ev
}

func (m model) Init() tea.Cmd {
	return getBranchesCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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
		case "space":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
		}
	case getBranchesMsg:
		m.choices = msg.branches
	}
	return m, nil
}

func (m model) View() tea.View {
	output := make([]string, 3)
	output = append(output, "Select Branch to Delete\n\n")
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		selected := " "
		_, ok := m.selected[i]
		if ok {
			selected = "x"
		}

		output = append(output,
			fmt.Sprintf("%s [%s] %s\n", cursor, selected, choice.name),
		)

	}
	output = append(output, "\nPress q or ctrl+c to exit\n")
	return tea.NewView(strings.Join(output, ""))
}

func getBranches() []branch {
	branches := make([]branch, 0)
	cmd := exec.Command("git", "branch", "-a")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return branches
	}
	for line := range strings.SplitSeq(string(output), "\n") {
		if len(strings.Trim(line, "\n")) == 0 {
			continue
		}
		isCurrent := strings.Contains(line, "*")
		name := strings.Trim(line, " *")
		remote := strings.HasPrefix(name, "remotes/")
		if remote {
			name = name[len("remotes/"):]
		}
		branches = append(branches, branch{name, isCurrent, remote})
	}
	return branches
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There is error: %v", err)
		os.Exit(1)
	}
}
