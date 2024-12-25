package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

// CUSTOM ITEM
type Task struct {
	status      status
	title       string
	description string
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

// MAIN MODEL

type Model struct {
	list list.Model
	err  error
}

// TODO: call this on tea.WindowSizeMsg
func (m *Model) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "To Do"
	m.list.SetItems([]list.Item{
		Task{status: todo, title: "Buy Milk", description: "Stawberry milk"},
		Task{status: todo, title: "Eat Sushi", description: "nigiri"},
		Task{status: todo, title: "Fold laundry", description: "or wear wrinkly t shirt"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func New() *Model {
	return &Model{}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

func main() {
	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("\nAn Error Occured: %v", err)
		os.Exit(1)
	}
}
