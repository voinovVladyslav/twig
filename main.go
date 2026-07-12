package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices  []branch
	errors   map[branch]string
	cursor   int
	selected map[int]struct{}
}

type branch struct {
	name    string
	current bool
	remote  bool
}

var (
	selectedStyle = lipgloss.NewStyle().Bold(true)
	remoteStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#a00000"))
	currentStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#0a0"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffef00"))
)

func initialModel() model {
	return model{
		choices:  make([]branch, 0),
		cursor:   0,
		errors:   make(map[branch]string),
		selected: make(map[int]struct{}),
	}
}

type getBranchesMsg struct {
	branches []branch
}

func getBranchesCmd() tea.Msg {
	branches := getBranches()
	ev := getBranchesMsg{branches}
	return ev
}

func deleteBranch(b *branch) error {
	if b.current {
		return errors.New("cannot delete current branch")
	}
	if b.name == "main" {
		return errors.New("cannot delete main branch, it's main")
	}

	flags := "-D"
	if b.remote {
		flags += "r"
	}
	cmd := exec.Command("git", "branch", flags, b.name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		lines := strings.Split(string(output), "\n")
		return errors.New(strings.Join(lines, " "))
	}
	return nil
}

type deleteBranchesMsg struct {
	msg string
	b   branch
}
type deleteBranchesMsgs struct {
	messages []deleteBranchesMsg
}

func deleteBranchesCmd(branches []branch) tea.Cmd {
	return func() tea.Msg {
		resp := make([]deleteBranchesMsg, 0)
		for _, b := range branches {
			err := deleteBranch(&b)
			if err != nil {
				resp = append(resp, deleteBranchesMsg{err.Error(), b})
			}
		}
		return deleteBranchesMsgs{resp}
	}
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

		case "r":
			clear(m.selected)
			clear(m.errors)
			clear(m.choices)
			return m, getBranchesCmd
		case "enter":
			toDelete := make([]branch, 0)
			for i := range m.selected {
				toDelete = append(toDelete, m.choices[i])
			}
			if len(toDelete) == 0 {
				return m, nil
			}
			clear(m.selected)
			m.cursor = 0
			return m, tea.Sequence(deleteBranchesCmd(toDelete), getBranchesCmd)
		}
	case getBranchesMsg:
		m.choices = msg.branches
	case deleteBranchesMsgs:
		clear(m.errors)
		for _, message := range msg.messages {
			m.errors[message.b] = message.msg
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	output := make([]string, 3)
	output = append(output, "Press Space to select branch\n")
	output = append(output, "Press Enter to Delete selected branches\n")
	output = append(output, "Press r to Delete selected branches\n\n")
	output = append(output, "Branches:\n")
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		line := choice.name
		selected := " "
		_, ok := m.selected[i]
		if ok {
			selected = "x"
			line = selectedStyle.Render(line)
		}

		if choice.remote {
			line = remoteStyle.Render(line)
		} else if choice.current {
			line = currentStyle.Render(line)
		}

		composedLine := fmt.Sprintf("%s [%s] %s", cursor, selected, line)
		errMsg, ok := m.errors[choice]
		if ok {
			composedLine += "   Err: " + errorStyle.Render(errMsg) + "\n"
		} else {
			composedLine += "\n"
		}
		output = append(output, composedLine)

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
