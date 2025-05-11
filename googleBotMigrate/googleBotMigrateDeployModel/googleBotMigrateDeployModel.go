package googleBotMigrateDeploy

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mcphee11/mcphee11-tui/utils"
)

const (
	padding  = 5
	maxWidth = 80
)

var bannerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#655ad5")).Padding(0, 1).Render
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

// Global variables
var (
	status        string
	stepsInStage  int
	GlobalProgram *tea.Program // To send messages from goroutines
	folderStage1  string
	folderStage2  string

	projectId, lang, flowName, keyPath, migrateType string
)

// --- New Message Types ---
type flowProcessedMsg struct{}
type stage1CompleteMsg struct{}
type stage2CompleteMsg struct{}
type internalUpdateStatusMsg struct {
	newStatus string
}

// --- Model ---
type model struct {
	progress         progress.Model
	stage1InProgress bool
	stage2InProgress bool
	processedCount   int
}

func MigrateLoadingMain(projectIdIn, langIn, flowNameIn, keyPathIn, migrateTypeIn string) {
	stepsInStage = 3
	projectId = projectIdIn
	lang = langIn
	flowName = flowNameIn
	keyPath = keyPathIn
	migrateType = migrateTypeIn
	status = fmt.Sprintf("Ready. Press 's' to start migration of Google Bot Locally. Once completed you will be able to Deploy it")

	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	GlobalProgram = p // Store the program instance globally

	if _, err := p.Run(); err != nil {
		utils.TuiLogger("Fatal", "(googleBotMigrateDeployModel) Could not start program")
	}
}

func (m model) Init() tea.Cmd {
	return nil // No initial command needed for progress
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "s":
			if m.stage1InProgress {
				return m, nil // Already syncing up
			}
			m.stage1InProgress = true
			m.processedCount = 0
			status = fmt.Sprintf("Starting migration of %d stages...", stepsInStage)

			m.progress.SetPercent(0) // Reset progress bar

			// Command to start the build process in a goroutine
			cmd := func() tea.Msg {
				// Pass necessary data to the build process
				// The GlobalProgram allows the goroutine to send messages back
				if migrateType == "buildDigitalBot" {
					go BuildDigitalBot(projectId, lang, flowName, keyPath, GlobalProgram)
				}
				if migrateType == "buildKnowledgeBase" {
					go BuildKnowledgeBaseCSV(projectId, lang, flowName, keyPath, GlobalProgram)
				}
				return nil // The goroutine will send messages asynchronously
			}
			return m, cmd
		case "u":
			if m.stage2InProgress {
				return m, nil // Already syncing up
			}
			m.stage2InProgress = true
			m.processedCount = -1
			m.progress.SetPercent(0) // Reset progress bar
			if !m.stage1InProgress {
				// Command to start the update process in a goroutine
				cmd := func() tea.Msg {
					// Pass necessary data to the backup process
					// The GlobalProgram allows the goroutine to send messages back
					// TODO: go RunUpdateProcess(flagName, shortName, color, icon, banner, region, environment, deploymentId, bucketName string GlobalProgram)
					return nil // The goroutine will send messages asynchronously
				}
				// Now start the Update
				return m, cmd
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case flowProcessedMsg:
		m.processedCount++
		var currentProgress float64
		if stepsInStage > 0 {
			currentProgress = float64(m.processedCount) / float64(stepsInStage)
			utils.TuiLogger("Info", "(googleBotMigrateDeployModel) flowProcessedMsg increase")
		} else {
			currentProgress = 1.0 // Or 0.0 if no flows
			utils.TuiLogger("Info", "(googleBotMigrateDeployModel) flowProcessedMsg Set % 1.0")
		}
		// Status updated by internalUpdateStatusMsg from the goroutine for more detail
		// status = fmt.Sprintf("Processing... %d/%d complete.", m.processedCount, stepsInStage)
		return m, m.progress.SetPercent(currentProgress)

	case stage1CompleteMsg:
		m.stage1InProgress = false
		// Status will be set by the runStage1Process or a final internalUpdateStatusMsg
		if !strings.Contains(status, "Build COMPLETED.") && !strings.HasPrefix(status, "ERROR:") { // Avoid overriding error messages
			status = "ERROR..."
			utils.TuiLogger("Error", "(googleBotMigrateDeployModel) stage1CompleteMsg ERROR...")
		}
		return m, m.progress.SetPercent(1.0)

	case stage2CompleteMsg:
		m.stage2InProgress = false
		if status != "Publish COMPLETED." && !strings.HasPrefix(status, "ERROR:") {
			status = "PUBLISH ERROR"
			utils.TuiLogger("Error", "(googleBotMigrateDeployModel) stage2CompleteMsg PUBLISH ERROR")
		}
		utils.TuiLogger("Info", "(googleBotMigrateDeployModel) stage2CompleteMsg Set % 1.0")
		return m, m.progress.SetPercent(1.0)
	case internalUpdateStatusMsg:
		status = msg.newStatus
		return m, nil

	case progress.FrameMsg: // This handles the animation rendering
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model) // Corrected type assertion
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	// Display current progress percentage directly in the status or help text for clarity
	return "\n" +
		pad + bannerStyle("Migrate Google BOT Progress") + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle(status) + "\n\n" +
		pad + helpStyle("Press 'q' to quit.")
}
