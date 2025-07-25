package flows

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mcphee11/mcphee11-tui/genesysLogin"
	"github.com/mcphee11/mcphee11-tui/utils"
)

// This function runs in a goroutine and performs the backup.
// It sends messages back to the Tea program via the 'p' instance.
func RunBackupProcess(flowId string, flows []map[string]string, totalFlows int, p *tea.Program) {
	// Get ORG details
	region, clientId, secret, err := genesysLogin.GenesysCreds()
	if err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("Failed to get Genesys Cloud credentials: %v", err))
	}

	// Helper to send status messages to the UI thread
	sendMsgToUI := func(msg tea.Msg) {
		if p != nil {
			p.Send(msg)
		}
	}
	sendStatusUpdate := func(t, s string) {
		sendMsgToUI(internalUpdateStatusMsg{newStatus: s})
		utils.TuiLogger(t, s) // logging output if enabled
	}

	if region == "" {
		sendStatusUpdate("Error", "ERROR: environment variable MCPHEE11_TUI_REGION is not set")
		sendMsgToUI(backupCompleteMsg{})
		return
	}

	if clientId == "" {
		sendStatusUpdate("Error", "ERROR: environment variable MCPHEE11_TUI_CLIENT_ID is not set")
		sendMsgToUI(backupCompleteMsg{})
		return
	}

	if secret == "" {
		sendStatusUpdate("Error", "ERROR: environment variable MCPHEE11_TUI_SECRET is not set")
		sendMsgToUI(backupCompleteMsg{})
		return
	}

	folderBackup = fmt.Sprintf("flowsBackup_%s", strconv.FormatInt(time.Now().Unix(), 10))
	sendStatusUpdate("Info", fmt.Sprintf("Creating backup directory: %s", folderBackup))
	err = os.Mkdir(folderBackup, 0777)
	if err != nil {
		sendStatusUpdate("Error", fmt.Sprintf("ERROR: creating directory %s: %v", folderBackup, err))
		sendMsgToUI(backupCompleteMsg{})
		return
	}

	if flowId == "ALL" {
		if totalFlows == 0 {
			sendStatusUpdate("Info", "No flows to back up.")
			sendMsgToUI(backupCompleteMsg{})
			return
		}
		for i, currentFlow := range flows {
			sendStatusUpdate("Info", fmt.Sprintf("Backing up flow %d/%d: %s...", i+1, totalFlows, currentFlow["id"]))
			archyCmd := exec.Command("archy", "export", "--exportType", "yaml", "--flowId", currentFlow["id"], "--clientId", clientId, "--clientSecret", secret, "--location", region, "--outputDir", folderBackup)

			if err := archyCmd.Run(); err != nil {
				sendStatusUpdate("Error", fmt.Sprintf("ERROR backing up %s: %v", currentFlow["id"], err))
				// Optionally, decide to stop or continue on error
				// For now, it continues and reports error for this flow.
				// TODO add better error logging
			}
			sendMsgToUI(flowProcessedMsg{}) // Signal one flow is processed
		}
	} else { // Single flow backup
		sendStatusUpdate("Info", fmt.Sprintf("Backing up single flow: %s...", flowId))
		archyCmd := exec.Command("archy", "export", "--exportType", "yaml", "--flowId", flowId, "--clientId", clientId, "--clientSecret", secret, "--location", region, "--outputDir", folderBackup)

		if err := archyCmd.Run(); err != nil {
			sendStatusUpdate("Error", fmt.Sprintf("ERROR backing up %s: %v", flowId, err))
		}
		sendMsgToUI(flowProcessedMsg{})
	}

	finalStatus := "Backup COMPLETED."
	if updateRequested && updateType == "tts" {
		finalStatus = fmt.Sprintf("Backup COMPLETED... Press u to start the upgrade of these flows to TTS Voice: %s", ttsSetting)
	}
	if updateRequested && updateType == "rePublish" {
		finalStatus = fmt.Sprintf("Backup COMPLETED... Press u to start the upgrade of these flows that relate to common module: %s", ttsSetting)
	}
	// TODO make error logging better
	if strings.Contains(status, "ERROR:") {
		finalStatus = "Backup process finished with errors. Check logs."
	}
	sendStatusUpdate("Info", finalStatus)
	sendMsgToUI(backupCompleteMsg{})
}
