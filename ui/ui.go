package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustmason/asteroids/world"
	"strings"
	"time"
)

type UIModel struct {
	world      *world.World
	playerID   string
	playerName string
	termWidth  int
	termHeight int
	keys       keyMap
}

func NewUIModel(w *world.World, playerID, playerName string, termWidth, termHeight int) UIModel {
	return UIModel{
		world:      w,
		playerID:   playerID,
		playerName: playerName,
		termWidth:  termWidth,
		termHeight: termHeight,
		keys:       keys,
	}
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Space key.Binding
	Tab   key.Binding
	Esc   key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "x"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h", "a"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l", "d"),
		key.WithHelp("→/l", "move right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type TickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m UIModel) Init() tea.Cmd {
	return doTick()
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		return m, doTick()
	case tea.WindowSizeMsg:
		m.termHeight = msg.Height
		m.termWidth = msg.Width
	}
	return m.handleMessage(msg)
}

func (m UIModel) handleMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			m.world.AcceleratePlayer(-1, m.playerID)
		case key.Matches(msg, m.keys.Down):
			m.world.AcceleratePlayer(1, m.playerID)
		case key.Matches(msg, m.keys.Left):
			m.world.RotatePlayer(-1, m.playerID)
		case key.Matches(msg, m.keys.Right):
			m.world.RotatePlayer(1, m.playerID)
		case key.Matches(msg, m.keys.Space):
			m.world.FirePlayer(m.playerID)
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

var (
	borderedBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#874BFD")).
				Padding(0).
				BorderTop(true).
				BorderLeft(true).
				BorderRight(true).
				BorderBottom(true)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
	docStyle = lipgloss.NewStyle().Padding(0)
)

// todo center the view when terminal size is beyond world size
var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

func (m UIModel) View() string {
	mainWidth := m.mainWidth()
	mainHeight := m.termHeight - 4 - 4

	// local copy, because Width/Height mutate it. this avoids `concurrent map write` panics
	mainStyle := lipgloss.NewStyle().Inherit(borderedBoxStyle).Width(m.world.W).Height(m.world.H)
	sbStyleLeft := lipgloss.NewStyle().Inherit(statusBarStyle).Width(m.termWidth / 2)
	sbStyleRight := lipgloss.NewStyle().Inherit(statusBarStyle).Width(m.termWidth / 2).Align(lipgloss.Right)

	var mainContents string
	mainContents = mainStyle.Render(m.world.Render(m.playerID, m.playerName, m.world.W, m.world.H))

	doc := strings.Builder{}
	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.Place(
			mainWidth, mainHeight, lipgloss.Center, lipgloss.Center, mainContents,
			lipgloss.WithWhitespaceChars("猫咪"),
			lipgloss.WithWhitespaceForeground(subtle),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			sbStyleLeft.Render(m.world.RenderWorldStatus()), // todo use status bar for short help text
			sbStyleRight.Render(m.world.RenderPosition(m.playerID)),
		),
	)
	doc.WriteString(ui)
	return docStyle.Render(doc.String())
}

func (m UIModel) mainWidth() int {
	return m.termWidth - 2 // minus border
}
