package googleBotMigrate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	googleBotMigrateDeploy "github.com/mcphee11/mcphee11-tui/googleBotMigrate/googleBotMigrateDeployModel"
	"github.com/mcphee11/mcphee11-tui/utils"
)

var (
	bannerStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).Background(lipgloss.Color("#655ad5")).Padding(0, 1).Render
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(5)
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(5)
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle().PaddingLeft(5)
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(5)

	focusedSubmitButton = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("[ Submit ]")
	blurredSubmitButton = fmt.Sprintf("[ %s ]", lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Submit"))
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 5),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.PromptStyle = noStyle
		t.TextStyle = blurredStyle
		t.Width = 200
		t.CharLimit = 320

		switch i {
		case 0:
			t.Placeholder = "Type of output eg: digitalBot or knowledgeBase"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Google Cloud Project ID"
		case 2:
			t.Placeholder = "Dialogflow Language Code. eg: en-AU"
		case 3:
			t.Placeholder = "Genesys Cloud Bot Flow Name"
		case 4:
			t.Placeholder = "Google Cloud Key Path eg: /path/to/key.json"
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {

				for i := 0; i < len(m.inputs); i++ {
					if m.inputs[i].Value() == "" {
						m.inputs[i].Focus()
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
						return m, nil
					}
				}

				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	fmt.Fprintf(&b, "\n     %s\n\n", bannerStyle("Google Bot Migration"))

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	buttonSubmit := &blurredSubmitButton
	if m.focusIndex == len(m.inputs) {
		buttonSubmit = &focusedSubmitButton
	}
	fmt.Fprintf(&b, "\n\n     %s\n\n", *buttonSubmit)

	b.WriteString(helpStyle.Render("ctrl+c to quit"))

	return b.String()
}

func MainInputs() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	returnedModel, err := p.Run()
	if err != nil {
		utils.TuiLogger("Fatal", "(googleBotMigrateModel) Could not start program")
	}
	// A space i getting added at the start.. not sure why right now??
	input0 := strings.TrimPrefix(returnedModel.(model).inputs[0].Value(), " ")
	input1 := strings.TrimPrefix(returnedModel.(model).inputs[1].Value(), " ")
	input2 := strings.TrimPrefix(returnedModel.(model).inputs[2].Value(), " ")
	input3 := strings.TrimPrefix(returnedModel.(model).inputs[3].Value(), " ")
	input4 := strings.TrimPrefix(returnedModel.(model).inputs[4].Value(), " ")

	if input0 != "knowledgeBase" && input0 != "digitalBot" {
		utils.TuiLogger("Info", "Please enter a valid type of output either: digitalBot or knowledgeBase")
	}

	if input0 == "digitalBot" {
		googleBotMigrateDeploy.MigrateLoadingMain(input1, input2, input3, input4, "buildDigitalBot")
	}
	if input0 == "knowledgeBase" {
		googleBotMigrateDeploy.MigrateLoadingMain(input1, input2, input3, input4, "buildKnowledgeBase")
	}

}
