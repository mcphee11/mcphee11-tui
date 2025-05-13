package ttsChanger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/mcphee11/mcphee11-tui/utils"
	"github.com/mypurecloud/platform-client-sdk-go/platformclientv2"
)

func CheckForArchy() (err error) {
	cmd := exec.Command("archy", "-version")
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("Ensure you have archy installed: https://developer.genesys.cloud/devapps/archy/")
	}

	return nil
}

func CurrentTTSVoices(config *platformclientv2.Configuration) []map[string]string {

	var voices []map[string]string

	apiIntegrationsInstance := platformclientv2.NewIntegrationsApiWithConfig(config)
	ttsVoices, err := getTTSVoicesEnhancedPages(apiIntegrationsInstance, 1)

	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) Error getting architect flows: %v", err))
		os.Exit(1)
	}

	pageNumber := 2
	for pageNumber <= *ttsVoices.PageCount {
		nextPage, err := getTTSVoicesEnhancedPages(apiIntegrationsInstance, pageNumber)
		if err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) %s", err))
			os.Exit(1)
		}
		*ttsVoices.Entities = append(*ttsVoices.Entities, *nextPage.Entities...)
		pageNumber++
	}

	for _, entity := range *ttsVoices.Entities {
		voices = append(voices, map[string]string{
			"title": *entity.Name,
			"id":    *entity.Id,
			"desc":  *entity.Language + " | " + *entity.Gender,
		})
	}

	return voices
}

func GetFlows(config *platformclientv2.Configuration, searchId string) []map[string]string {
	var flowIds []string
	var final []map[string]string
	dependencies := getDependencyTracking(config, searchId)

	for i := range dependencies {
		flowIds = append(flowIds, dependencies[i]["id"])
	}
	latestPublished, err := GetFlowsCUSTOM(config, flowIds)

	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) %s", err))
		return nil
	}

	for _, ver := range latestPublished {
		for _, dep := range dependencies {
			if ver["version"] == dep["version"] && ver["id"] == dep["id"] {
				final = append(final, ver)
			}
		}
	}

	return final
}

func getDependencyTracking(config *platformclientv2.Configuration, searchId string) []map[string]string {

	var flows []map[string]string

	apiInstance := platformclientv2.NewArchitectApiWithConfig(config)

	id := searchId                     // Object ID
	var version string                 // Object version
	objectType := "TTSVOICE"           // Object type
	consumedResources := true          // Include resources this item consumes
	consumingResources := true         // Include resources that consume this item
	var consumedResourceType []string  // Types of consumed resources to return, if consumed resources are requested
	var consumingResourceType []string // Types of consuming resources to return, if consuming resources are requested
	consumedResourceRequest := true    // Indicate that this is going to look up a consumed resource object

	data, _, err := apiInstance.GetArchitectDependencytrackingObject(id, version, objectType, consumedResources, consumingResources, consumedResourceType, consumingResourceType, consumedResourceRequest)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) Error calling GetArchitectDependencytrackingObject: %v", err))
		os.Exit(1)
	}
	for _, entity := range *data.ConsumingResources {
		flows = append(flows, map[string]string{
			"title":   *entity.Name,
			"id":      *entity.Id,
			"desc":    *entity.VarType + " | " + *entity.Version,
			"version": *entity.Version,
		})
	}
	return flows
}

func getTTSVoicesEnhancedPages(apiInstance *platformclientv2.IntegrationsApi, page int) (ttsVoices *platformclientv2.Ttsvoiceentitylisting, err error) {
	engineId := "genesys_enhanced" // The engine ID
	//var pageNumber int // Page number
	pageSize := 500 // Page size
	// Get a list of voices for a TTS engine
	data, _, err := apiInstance.GetIntegrationsSpeechTtsEngineVoices(engineId, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("Error calling GetIntegrationsSpeechTtsEngineVoices: %v\n", err)
	} else {
		return data, nil
	}
}

func GetFlowsCUSTOM(config *platformclientv2.Configuration, flows []string) (flowReturned []map[string]string, err error) {
	var allFlows []map[string]string

	flowPage, err := getFlowsPageCUSTOM(config, flows, 1)

	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) %s", err))
		return nil, err
	}

	pageNumber := 2
	for pageNumber <= *flowPage.PageCount {
		nextPage, err := getFlowsPageCUSTOM(config, flows, pageNumber)
		if err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) %s", err))
			os.Exit(1)
		}
		*flowPage.Entities = append(*flowPage.Entities, *nextPage.Entities...)
		pageNumber++
	}

	for _, entity := range *flowPage.Entities {
		if entity.PublishedVersion == nil || entity.PublishedVersion.Id == nil {
			utils.TuiLogger("Warning", fmt.Sprintf("(ttsChanger) Skipping entity with missing PublishedVersion or Id: %s", *entity.Id))
			continue
		}
		allFlows = append(allFlows, map[string]string{
			"title":   *entity.Name,
			"id":      *entity.Id,
			"desc":    *entity.PublishedVersion.Id,
			"version": *entity.PublishedVersion.Id,
		})
	}
	return allFlows, nil
}

func getFlowsPageCUSTOM(config *platformclientv2.Configuration, flows []string, pageNumber int) (response *FlowentitylistingCUSTOM, err error) {

	url := fmt.Sprintf("https://api.mypurecloud.com.au/api/v2/flows?pageNumber=%d&pageSize=500&sortBy=asc&id=%s", pageNumber, strings.Join(flows, ","))
	token := config.AccessToken

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) request creation: %s", err))
		return nil, err
	}

	// Add the Authorization header with the bearer token
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) request sending: %s", err))
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) Request failed with status code: %d", resp.StatusCode))
		bodyBytes, _ := io.ReadAll(resp.Body)
		utils.TuiLogger("Info", fmt.Sprintf("(ttsChanger) Response body: %s", string(bodyBytes)))
		return nil, err
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) Reading body: %s", err))
		return nil, err
	}

	// Parse the JSON response into the custom struct
	var flowResponse FlowentitylistingCUSTOM
	err = json.Unmarshal(bodyBytes, &flowResponse)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) Unmarshalling JSON: %s", err))
		return nil, err
	}

	return &flowResponse, nil
}

// Once SDK is fixed will move to this func
func getFlowPage(apiInstance *platformclientv2.ArchitectApi, page int, flowIds []string) (flowPage *platformclientv2.Flowentitylisting, err error) {
	// Get the flow page
	var varType []string // Type
	//var pageNumber int   // Page number
	pageSize := 200      // Page size
	var sortBy string    // Sort by
	var sortOrder string // Sort order
	//var id []string              // ID
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

	data, _, er := apiInstance.GetFlows(varType, page, pageSize, sortBy, sortOrder, flowIds, name, description, nameOrDescription, publishVersionId, editableBy, lockedBy, lockedByClientId, secure, deleted, includeSchemas, publishedAfter, publishedBefore, divisionId)
	if er != nil {
		return nil, fmt.Errorf("failed to get flow page: %w", er)
	}
	return data, nil
}
