package ttsChanger

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mcphee11/mcphee11-tui/genesysLogin"
	"github.com/mypurecloud/platform-client-sdk-go/platformclientv2"
)

func TTSChanger(title string) {
	cmd := exec.Command("archy", "-version")
	err := cmd.Run()

	if err != nil {
		fmt.Println("Ensure you have archy installed: https://developer.genesys.cloud/devapps/archy/")
		fmt.Println("Error: ", err)
		return
	}

	config, err := genesysLogin.GenesysLogin()
	if err != nil {
		fmt.Println(err)
		return
	}

	apiInstance := platformclientv2.NewArchitectApiWithConfig(config)
	page, err := getFlowPage(apiInstance, 1)

	if err != nil {
		fmt.Println("Error getting architect flows", err)
		os.Exit(1)
	}

	fmt.Println(page)

}

func getFlowPage(apiInstance *platformclientv2.ArchitectApi, page int) (flowPage *platformclientv2.Flowentitylisting, err error) {
	// Get the flow page
	var varType []string         // Type
	var pageNumber int           // Page number
	var pageSize int             // Page size
	var sortBy string            // Sort by
	var sortOrder string         // Sort order
	var id []string              // ID
	var name string              // Name
	var description string       // Description
	var nameOrDescription string // Name or description
	var publishVersionId string  // Publish version ID
	var editableBy string        // Editable by
	var lockedBy string          // Locked by
	var lockedByClientId string  // Locked by client ID
	var secure string            // Secure
	var deleted bool             // Include deleted
	var includeSchemas bool      // Include variable schemas
	var publishedAfter string    // Published after
	var publishedBefore string   // Published before
	var divisionId []string      // division ID(s)

	data, response, er := apiInstance.GetFlows(varType, pageNumber, pageSize, sortBy, sortOrder, id, name, description, nameOrDescription, publishVersionId, editableBy, lockedBy, lockedByClientId, secure, deleted, includeSchemas, publishedAfter, publishedBefore, divisionId)
	if er != nil {
		return nil, fmt.Errorf("failed to get flow page: %w | Response: %v", er, response.StatusCode)
	}
	return data, nil
}
