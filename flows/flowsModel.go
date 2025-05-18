package flows

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
	status           string
	updateRequested  bool
	updateType       string
	flowCount        int
	flowId           string
	savedFlows       []map[string]string
	savedFlowId      string
	GlobalProgram    *tea.Program // To send messages from goroutines
	ttsSetting       string
	ttsGetting       string
	ttsSettingEngine string
	ttsGettingEngine string
	folderBackup     string
	folderUpdate     string
)

// --- New Message Types ---
type flowProcessedMsg struct{}
type backupCompleteMsg struct{}
type updateCompleteMsg struct{}
type internalUpdateStatusMsg struct {
	newStatus string
}

// --- Model ---
type model struct {
	progress       progress.Model
	backingUp      bool
	publishing     bool
	processedCount int
}

func FlowsLoadingMainBackup(flowId string, flows []map[string]string, ttsGet, ttsSet, ttsEngineGet, ttsEngineSet, update string, updateRequired bool) {
	flowCount = len(flows)
	savedFlowId = flowId
	savedFlows = flows
	ttsSetting = ttsSet
	ttsGetting = ttsGet
	ttsSettingEngine = ttsEngineSet
	ttsGettingEngine = ttsEngineGet
	updateRequested = updateRequired
	updateType = update
	if updateRequired {
		if flowId == "ALL" {
			status = fmt.Sprintf("Ready. Press 's' to start backup of %d flows. Once completed you will be able to update them", flowCount)
		} else {
			status = fmt.Sprintf("Ready. Press 's' to start backup of %d flow. Once completed you will be able to update it", 1)
		}
	} else {
		if flowId == "ALL" {
			status = fmt.Sprintf("Ready. Press 's' to start backup of %d flows.", flowCount)
		} else {
			status = fmt.Sprintf("Ready. Press 's' to start backup of %d flow.", 1)
		}
	}

	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	GlobalProgram = p // Store the program instance globally

	if _, err := p.Run(); err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("(flowsModel) Error running program: %s", err))
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
			if m.backingUp {
				return m, nil // Already syncing up
			}
			m.backingUp = true
			m.processedCount = 0
			if flowId == "ALL" {
				status = fmt.Sprintf("Starting backup of %d flows...", flowCount)
			} else {
				status = fmt.Sprintf("Starting backup of %s flow...", flowId)
			}

			m.progress.SetPercent(0) // Reset progress bar

			// Command to start the backup process in a goroutine
			cmd := func() tea.Msg {
				// Pass necessary data to the backup process
				// The GlobalProgram allows the goroutine to send messages back
				go RunBackupProcess(savedFlowId, savedFlows, flowCount, GlobalProgram)
				return nil // The goroutine will send messages asynchronously
			}
			return m, cmd
		case "u":
			if m.publishing {
				return m, nil // Already syncing up
			}
			m.publishing = true
			m.processedCount = -1
			m.progress.SetPercent(0) // Reset progress bar
			if updateRequested && !m.backingUp {
				// Command to start the update process in a goroutine
				cmd := func() tea.Msg {
					// Pass necessary data to the backup process
					// The GlobalProgram allows the goroutine to send messages back
					go RunUpdateProcess(flowCount, GlobalProgram, updateType)
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
		if flowCount > 0 {
			currentProgress = float64(m.processedCount) / float64(flowCount)
			utils.TuiLogger("Info", "(flowsModel) flowProcessedMsg increase")
		} else {
			currentProgress = 1.0 // Or 0.0 if no flows
			utils.TuiLogger("Info", "(flowsModel) flowProcessedMsg Set % 1.0")
		}
		// Status updated by internalUpdateStatusMsg from the goroutine for more detail
		return m, m.progress.SetPercent(currentProgress)

	case backupCompleteMsg:
		m.backingUp = false
		// Status will be set by the runBackupProcess or a final internalUpdateStatusMsg
		if !strings.Contains(status, "Backup COMPLETED.") && !strings.HasPrefix(status, "ERROR:") { // Avoid overriding error messages
			status = "ERROR..."
			utils.TuiLogger("Error", "(flowsModel) stage1CompleteMsg ERROR...")
		}
		return m, m.progress.SetPercent(1.0)

	case updateCompleteMsg:
		m.publishing = false
		if status != "Publish COMPLETED." && !strings.HasPrefix(status, "ERROR:") {
			status = "PUBLISH ERROR"
			utils.TuiLogger("Error", "(flowsModel) stage2CompleteMsg PUBLISH ERROR")
		}
		utils.TuiLogger("Info", "(flowsModel) stage2CompleteMsg Set % 1.0")
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
		pad + bannerStyle("Flows Backup Progress") + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle(status) + "\n\n" +
		pad + helpStyle("Press 'q' to quit.")
}
