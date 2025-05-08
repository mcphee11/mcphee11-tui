package flows

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// This function runs in a goroutine and performs the backup.
// It sends messages back to the Tea program via the 'p' instance.
func RunBackupProcess(flowId string, flows []map[string]string, totalFlows int, p *tea.Program) {
	// Helper to send status messages to the UI thread
	sendMsgToUI := func(msg tea.Msg) {
		if p != nil {
			p.Send(msg)
		}
	}
	sendStatusUpdate := func(s string) {
		sendMsgToUI(internalUpdateStatusMsg{newStatus: s})
	}

	region := os.Getenv("MCPHEE11_TUI_REGION")
	if region == "" {
		sendStatusUpdate("ERROR: environment variable MCPHEE11_TUI_REGION is not set")
		sendMsgToUI(backupCompleteMsg{})
		return
	}

	clientId := os.Getenv("MCPHEE11_TUI_CLIENT_ID")
	if clientId == "" {
		sendStatusUpdate("ERROR: environment variable MCPHEE11_TUI_CLIENT_ID is not set")
		sendMsgToUI(backupCompleteMsg{})
		return
	}

	secret := os.Getenv("MCPHEE11_TUI_SECRET")
	if secret == "" {
		sendStatusUpdate("ERROR: environment variable MCPHEE11_TUI_SECRET is not set")
		sendMsgToUI(backupCompleteMsg{})
		return
	}

	folderBackup = fmt.Sprintf("flowsBackup_%s", strconv.FormatInt(time.Now().Unix(), 10))
	sendStatusUpdate(fmt.Sprintf("Creating backup directory: %s", folderBackup))
	err := os.Mkdir(folderBackup, 0777)
	if err != nil {
		sendStatusUpdate(fmt.Sprintf("ERROR: creating directory %s: %v", folderBackup, err))
		sendMsgToUI(backupCompleteMsg{})
		return
	}

	if flowId == "ALL" {
		if totalFlows == 0 {
			sendStatusUpdate("No flows to back up.")
			sendMsgToUI(backupCompleteMsg{})
			return
		}
		for i, currentFlow := range flows {
			sendStatusUpdate(fmt.Sprintf("Backing up flow %d/%d: %s...", i+1, totalFlows, currentFlow["id"]))
			archyCmd := exec.Command("archy", "export", "--exportType", "yaml", "--flowId", currentFlow["id"], "--clientId", clientId, "--clientSecret", secret, "--location", region, "--outputDir", folderBackup)

			if err := archyCmd.Run(); err != nil {
				sendStatusUpdate(fmt.Sprintf("ERROR backing up %s: %v", currentFlow["id"], err))
				// Optionally, decide to stop or continue on error
				// For now, it continues and reports error for this flow.
				// TODO add better error logging
			}
			sendMsgToUI(flowProcessedMsg{}) // Signal one flow is processed
		}
	} else { // Single flow backup
		sendStatusUpdate(fmt.Sprintf("Backing up single flow: %s...", flowId))
		archyCmd := exec.Command("archy", "export", "--exportType", "yaml", "--flowId", flowId, "--clientId", clientId, "--clientSecret", secret, "--location", region, "--outputDir", folderBackup)

		if err := archyCmd.Run(); err != nil {
			sendStatusUpdate(fmt.Sprintf("ERROR backing up %s: %v", flowId, err))
		}
		sendMsgToUI(flowProcessedMsg{})
	}

	finalStatus := "Backup COMPLETED."
	if updateRequested {
		finalStatus = fmt.Sprintf("Backup COMPLETED... Press u to start the upgrade of these flows to TTS Voice: %s", ttsSetting)
	}
	// TODO make error logging better
	if strings.Contains(status, "ERROR:") {
		finalStatus = "Backup process finished with errors. Check logs."
	}
	sendStatusUpdate(finalStatus)
	sendMsgToUI(backupCompleteMsg{})
}
