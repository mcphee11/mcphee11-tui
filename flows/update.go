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

// This function runs in a goroutine and performs the update.
// It sends messages back to the Tea program via the 'p' instance.
func RunUpdateProcess(totalFlows int, p *tea.Program) {
	// Helper to send status messages to the UI thread
	sendMsgToUI := func(msg tea.Msg) {
		if p != nil {
			p.Send(msg)
		}
	}
	sendStatusUpdate := func(s string) {
		sendMsgToUI(internalUpdateStatusMsg{newStatus: s})
	}

	// Sent to zero out progress bar
	sendMsgToUI(flowProcessedMsg{})

	region := os.Getenv("MCPHEE11_TUI_REGION")
	if region == "" {
		sendStatusUpdate("ERROR: environment variable MCPHEE11_TUI_REGION is not set")
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	clientId := os.Getenv("MCPHEE11_TUI_CLIENT_ID")
	if clientId == "" {
		sendStatusUpdate("ERROR: environment variable MCPHEE11_TUI_CLIENT_ID is not set")
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	secret := os.Getenv("MCPHEE11_TUI_SECRET")
	if secret == "" {
		sendStatusUpdate("ERROR: environment variable MCPHEE11_TUI_SECRET is not set")
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	folderUpdate := fmt.Sprintf("flowsUpdate_%s", strconv.FormatInt(time.Now().Unix(), 10))
	sendStatusUpdate(fmt.Sprintf("Creating update directory: %s", folderUpdate))
	err := os.Mkdir(folderUpdate, 0777)
	if err != nil {
		sendStatusUpdate(fmt.Sprintf("ERROR: creating directory %s: %v", folderUpdate, err))
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	// Getting BackedUp files
	filesBackedUp, err := os.ReadDir(folderBackup)
	if err != nil {
		sendStatusUpdate(fmt.Sprintf("ERROR: reading directory %s: %v", folderBackup, err))
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	var backedUpYamlFiles []string
	for _, file := range filesBackedUp {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			backedUpYamlFiles = append(backedUpYamlFiles, file.Name())
		}
	}

	// Save modified backed-up files as new files in the folderUpdate directory
	for _, fileName := range backedUpYamlFiles {
		oldFilePath := fmt.Sprintf("%s/%s", folderBackup, fileName)
		newFilePath := fmt.Sprintf("%s/%s", folderUpdate, fileName)

		content, err := os.ReadFile(oldFilePath)
		if err != nil {
			sendStatusUpdate(fmt.Sprintf("ERROR: reading file %s: %v", oldFilePath, err))
			continue
		}

		updatedContent := strings.ReplaceAll(string(content), fmt.Sprintf("voice: %s", ttsGetting), fmt.Sprintf("voice: %s", ttsSetting))

		err = os.WriteFile(newFilePath, []byte(updatedContent), 0644)
		if err != nil {
			sendStatusUpdate(fmt.Sprintf("ERROR: writing file %s: %v", newFilePath, err))
			continue
		}

		sendStatusUpdate(fmt.Sprintf("Saved updated voice file: %s", fileName))
	}

	// Getting Updated files
	files, err := os.ReadDir(folderUpdate)
	if err != nil {
		sendStatusUpdate(fmt.Sprintf("ERROR: reading directory %s: %v", folderUpdate, err))
		fmt.Printf("ERROR: reading directory %s: %v", folderUpdate, err)
		os.Exit(1)
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	var updatedYamlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			updatedYamlFiles = append(updatedYamlFiles, file.Name())
		}
	}

	// Running archy update
	if len(updatedYamlFiles) > 1 {
		if len(updatedYamlFiles) == 0 {
			sendStatusUpdate("No flows to update up.")
			fmt.Printf("Line: 115")
			os.Exit(1)
			sendMsgToUI(updateCompleteMsg{})
			return
		}

		for i, currentFlow := range updatedYamlFiles {
			sendStatusUpdate(fmt.Sprintf("Publishing up flow %d/%d: %s...", i+1, totalFlows, currentFlow))
			archyCmd := exec.Command("archy", "publish", "--forceUnlock", "--clientId", clientId, "--clientSecret", secret, "--location", region, "--file", fmt.Sprintf("%s/%s", folderUpdate, currentFlow))

			if err := archyCmd.Run(); err != nil {
				fmt.Printf("ERROR 126 up %s: %v", currentFlow, err)
				os.Exit(1)
				sendStatusUpdate(fmt.Sprintf("ERROR backing up %s: %v", currentFlow, err))
				// Optionally, decide to stop or continue on error
				// For now, it continues and reports error for this flow.
				// TODO add better error logging
			}
			sendMsgToUI(flowProcessedMsg{}) // Signal one flow is processed
		}
	} else { // Single flow update
		sendStatusUpdate(fmt.Sprintf("Publishing up single flow: %s...", updatedYamlFiles[0]))
		archyCmd := exec.Command("archy", "publish", "--forceUnlock", "--clientId", clientId, "--clientSecret", secret, "--location", region, "--file", fmt.Sprintf("%s/%s", folderUpdate, updatedYamlFiles[0]))

		if err := archyCmd.Run(); err != nil {
			fmt.Printf("ERROR 140 %s", err)
			os.Exit(1)
			sendStatusUpdate(fmt.Sprintf("ERROR publishing up %s: %v", updatedYamlFiles[0], err))
		}
		sendMsgToUI(flowProcessedMsg{})
	}

	finalStatus := "Publish COMPLETED."
	// TODO make error logging better
	if strings.Contains(status, "ERROR:") {
		finalStatus = "Publish process finished with errors. Check logs."
	}
	sendStatusUpdate(finalStatus)
	sendMsgToUI(updateCompleteMsg{})
}
