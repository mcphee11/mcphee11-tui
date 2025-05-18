package flows

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mcphee11/mcphee11-tui/utils"
)

// This function runs in a goroutine and performs the update.
// It sends messages back to the Tea program via the 'p' instance.
func RunUpdateProcess(totalFlows int, p *tea.Program, updateType string) {
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

	// Sent to zero out progress bar
	sendMsgToUI(flowProcessedMsg{})

	region := os.Getenv("MCPHEE11_TUI_REGION")
	if region == "" {
		sendStatusUpdate("Error", "ERROR: environment variable MCPHEE11_TUI_REGION is not set")
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	clientId := os.Getenv("MCPHEE11_TUI_CLIENT_ID")
	if clientId == "" {
		sendStatusUpdate("Error", "ERROR: environment variable MCPHEE11_TUI_CLIENT_ID is not set")
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	secret := os.Getenv("MCPHEE11_TUI_SECRET")
	if secret == "" {
		sendStatusUpdate("Error", "ERROR: environment variable MCPHEE11_TUI_SECRET is not set")
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	folderUpdate := fmt.Sprintf("flowsUpdate_%s", strconv.FormatInt(time.Now().Unix(), 10))
	sendStatusUpdate("Info", fmt.Sprintf("Creating update directory: %s", folderUpdate))
	err := os.Mkdir(folderUpdate, 0777)
	if err != nil {
		sendStatusUpdate("Info", fmt.Sprintf("ERROR: creating directory %s: %v", folderUpdate, err))
		sendMsgToUI(updateCompleteMsg{})
		return
	}

	// Getting BackedUp files
	filesBackedUp, err := os.ReadDir(folderBackup)
	if err != nil {
		sendStatusUpdate("Error", fmt.Sprintf("ERROR: reading directory %s: %v", folderBackup, err))
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
			sendStatusUpdate("Error", fmt.Sprintf("ERROR: reading file %s: %v", oldFilePath, err))
			continue
		}

		var updatedContent string
		if updateType == "tts" {
			utils.TuiLogger("Info", fmt.Sprintf("ttsSetting == %s", ttsSetting))
			if ttsSetting == "Default" {
				updatedContent = strings.ReplaceAll(string(content), fmt.Sprintf("%s:", ttsGettingEngine), "defaultEngine:")
				updatedContent = strings.ReplaceAll(string(updatedContent), fmt.Sprintf("voice: %s", ttsGetting), "defaultVoice: true")
			} else {
				// check for existing
				updatedContent = strings.ReplaceAll(string(content), fmt.Sprintf("%s:", ttsGettingEngine), fmt.Sprintf("%s:", ttsSettingEngine))
				updatedContent = strings.ReplaceAll(string(updatedContent), fmt.Sprintf("voice: %s", ttsGetting), fmt.Sprintf("voice: %s", ttsSetting))
				// check if default
				updatedContent = strings.ReplaceAll(string(updatedContent), "defaultEngine:", fmt.Sprintf("%s:", ttsSettingEngine))
				updatedContent = strings.ReplaceAll(string(updatedContent), "defaultVoice: true", fmt.Sprintf("voice: %s", ttsSetting))
			}
		}
		if updateType == "rePublish" {
			// no changes to yaml files just rePublishing them
			updatedContent = string(content)
		}

		err = os.WriteFile(newFilePath, []byte(updatedContent), 0644)
		if err != nil {
			sendStatusUpdate("Error", fmt.Sprintf("ERROR: writing file %s: %v", newFilePath, err))
			continue
		}
		sendStatusUpdate("Info", fmt.Sprintf("Saved updated voice file: %s", fileName))
	}

	// Getting Updated files
	files, err := os.ReadDir(folderUpdate)
	if err != nil {
		sendStatusUpdate("Fatal", fmt.Sprintf("ERROR: reading directory %s: %v", folderUpdate, err))
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
			sendStatusUpdate("Info", "No flows to update up.")
			sendMsgToUI(updateCompleteMsg{})
			return
		}

		for i, currentFlow := range updatedYamlFiles {
			sendStatusUpdate("Info", fmt.Sprintf("Publishing up flow %d/%d: %s...", i+1, totalFlows, currentFlow))
			archyCmd := exec.Command("archy", "publish", "--forceUnlock", "--clientId", clientId, "--clientSecret", secret, "--location", region, "--file", fmt.Sprintf("%s/%s", folderUpdate, currentFlow))

			if err := archyCmd.Run(); err != nil {
				sendStatusUpdate("Fatal", fmt.Sprintf("ERROR backing up %s: %v", currentFlow, err))
				// os.Exit triggered on Fatal as want to ensure ALL are backed up first
			}
			sendMsgToUI(flowProcessedMsg{}) // Signal one flow is processed
		}
	} else { // Single flow update
		sendStatusUpdate("Info", fmt.Sprintf("Publishing up single flow: %s...", updatedYamlFiles[0]))
		archyCmd := exec.Command("archy", "publish", "--forceUnlock", "--clientId", clientId, "--clientSecret", secret, "--location", region, "--file", fmt.Sprintf("%s/%s", folderUpdate, updatedYamlFiles[0]))

		if err := archyCmd.Run(); err != nil {
			sendStatusUpdate("Error", fmt.Sprintf("ERROR publishing up %s: %v", updatedYamlFiles[0], err))
			// decided not to Fatal so others would try to publish
		}
		sendMsgToUI(flowProcessedMsg{})
	}

	finalStatus := "Publish COMPLETED."
	// TODO make error logging better
	if strings.Contains(status, "ERROR:") {
		finalStatus = "Publish process finished with errors. Check logs."
	}
	sendStatusUpdate("Info", finalStatus)
	sendMsgToUI(updateCompleteMsg{})
}
