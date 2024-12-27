package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 1

const (
	todo status = iota
	inProgress
	done
)

// STYLING

var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// CUSTOM ITEM
type Task struct {
	status      status
	title       string
	description string
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t *Task) FilterValue() string {
	return t.title
}

func (t *Task) Title() string {
	return t.title
}

func (t *Task) Description() string {
	return t.description
}

// MAIN MODEL

type Model struct {
	focused  status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

// TODO: call this on tea.WindowSizeMsg
func (m *Model) initLists(width, height int) {
	frameWidth, frameHeight := columnStyle.GetFrameSize()
	m.lists = []list.Model{
		list.New([]list.Item{}, list.NewDefaultDelegate(), (width-frameWidth)/divisor, (height-frameHeight)/2),
		list.New([]list.Item{}, list.NewDefaultDelegate(), (width-frameWidth)/divisor, (height-frameHeight)/2),
		list.New([]list.Item{}, list.NewDefaultDelegate(), (width-frameWidth)/divisor, (height-frameHeight)/2),
	}

	m.lists[todo].SetShowHelp(false)
	m.lists[inProgress].SetShowHelp(false)
	m.lists[done].SetShowHelp(false)

	// Init To dos
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		&Task{status: todo, title: "Buy Milk", description: "Stawberry milk"},
		&Task{status: todo, title: "Eat Sushi", description: "nigiri"},
		&Task{status: todo, title: "Fold laundry", description: "or wear wrinkly t shirt"},
	})

	// Init in Progress
	m.lists[inProgress].Title = "In Progress"
	m.lists[inProgress].SetItems([]list.Item{
		&Task{status: inProgress, title: "Write Code", description: "In GO"},
		&Task{status: inProgress, title: "KanbanCLI", description: "3/5 Progress"},
	})

	// Init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		&Task{status: done, title: "Stay Cool", description: "as a cucumber"},
		&Task{status: done, title: "Basic Structure", description: "Styling and switching between multiple list"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New() *Model {
	return &Model{}
}

func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	if selectedItem == nil {
		return nil
	}

	selectedTask := selectedItem.(*Task)
	currentIndex := m.lists[m.focused].Index()

	// Remove from current list
	m.lists[m.focused].RemoveItem(currentIndex)

	// Update status and add to next list
	oldStatus := selectedTask.status
	selectedTask.Next()

	// Only move if status actually changed
	if oldStatus != selectedTask.status {
		m.lists[selectedTask.status].InsertItem(0, selectedTask)
	} else {
		// If status didn't change, put it back in the original list
		m.lists[m.focused].InsertItem(currentIndex, selectedTask)
	}

	return nil
}

// TODO: go to next list
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

// TODO: goto prev list
func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		}
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgressView := m.lists[inProgress].View()
		doneView := m.lists[done].View()
		switch m.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(lipgloss.Left, columnStyle.Render(todoView), focusedStyle.Render(inProgressView), columnStyle.Render(doneView))
		case done:
			return lipgloss.JoinHorizontal(lipgloss.Left, columnStyle.Render(todoView), columnStyle.Render(inProgressView), focusedStyle.Render(doneView))
		default:
			return lipgloss.JoinHorizontal(lipgloss.Left, focusedStyle.Render(todoView), columnStyle.Render(inProgressView), columnStyle.Render(doneView))
		}
	} else {
		return "Loading..."
	}
}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("\nAn Error Occured: %v", err)
		os.Exit(1)
	}
}
