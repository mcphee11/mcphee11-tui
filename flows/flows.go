package flows

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mcphee11/mcphee11-tui/genesysLogin"
	"github.com/mcphee11/mcphee11-tui/utils"
	"github.com/mypurecloud/platform-client-sdk-go/platformclientv2"
)

func GetFlows(config *platformclientv2.Configuration, objType, searchId string) []map[string]string {
	var flowIds []string
	var final []map[string]string
	utils.TuiLogger("Info", fmt.Sprintf("Dependency searching objectType: %s with id: %s", objType, searchId))
	dependencies := getDependencyTracking(config, objType, searchId)

	for i := range dependencies {
		// Check if dependencies[i]["id"] already exists in previous dependencies
		duplicate := false
		for j := 0; j < i; j++ {
			if dependencies[j]["id"] == dependencies[i]["id"] {
				duplicate = true
				break
			}
		}
		if !duplicate {
			flowIds = append(flowIds, dependencies[i]["id"])
		}
	}

	// check for flowIds length to stop Error 413 for to much data
	var latestPublished []map[string]string
	if len(flowIds) > 100 {
		utils.TuiLogger("Info", fmt.Sprintf("%s", flowIds))
		chunkSize := 100
		for i := 0; i < len(flowIds); i += chunkSize {
			end := i + chunkSize
			if end > len(flowIds) {
				end = len(flowIds)
			}
			chunk := flowIds[i:end]
			utils.TuiLogger("Info", fmt.Sprintf("Processing chunk of %d flowIds (from index %d to %d)", len(chunk), i, end-1))
			chunkResponse, err := GetFlowsCUSTOM(config, chunk)
			if err != nil {
				utils.TuiLogger("Error", fmt.Sprintf("Error processing chunk of flowIds: %s", err))
				if latestPublished == nil {
					return nil
				}
				return latestPublished
			}
			latestPublished = append(latestPublished, chunkResponse...)
		}
	} else {
		latestPublishedResponse, err := GetFlowsCUSTOM(config, flowIds)
		if err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("%s", err))
			return nil
		}
		latestPublished = latestPublishedResponse
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

func getDependencyTracking(config *platformclientv2.Configuration, objType, searchId string) []map[string]string {

	var flows []map[string]string

	apiInstance := platformclientv2.NewArchitectApiWithConfig(config)

	id := searchId                     // Object ID
	var version string                 // Object version
	objectType := objType              // Object type
	consumedResources := true          // Include resources this item consumes
	consumingResources := true         // Include resources that consume this item
	var consumedResourceType []string  // Types of consumed resources to return, if consumed resources are requested
	var consumingResourceType []string // Types of consuming resources to return, if consuming resources are requested
	consumedResourceRequest := true    // Indicate that this is going to look up a consumed resource object

	data, _, err := apiInstance.GetArchitectDependencytrackingObject(id, version, objectType, consumedResources, consumingResources, consumedResourceType, consumingResourceType, consumedResourceRequest)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("Error calling GetArchitectDependencytrackingObject: %v", err))
		os.Exit(1)
	}
	for _, entity := range *data.ConsumingResources {
		// Skip deleted entities
		if entity.Deleted != nil && *entity.Deleted {
			continue
		}
		
		flows = append(flows, map[string]string{
			"title":   *entity.Name,
			"id":      *entity.Id,
			"desc":    *entity.VarType + " | " + *entity.Version,
			"version": *entity.Version,
		})
	}
	return flows
}

func GetFlowsCUSTOM(config *platformclientv2.Configuration, flows []string) (flowReturned []map[string]string, err error) {
	var allFlows []map[string]string

	flowPage, err := getFlowsPageCUSTOM(config, flows, 1)

	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("%s", err))
		return nil, err
	}

	pageNumber := 2
	for pageNumber <= *flowPage.PageCount {
		nextPage, err := getFlowsPageCUSTOM(config, flows, pageNumber)
		if err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("%s", err))
			os.Exit(1)
		}
		*flowPage.Entities = append(*flowPage.Entities, *nextPage.Entities...)
		pageNumber++
	}

	for _, entity := range *flowPage.Entities {
		if entity.PublishedVersion == nil || entity.PublishedVersion.Id == nil {
			utils.TuiLogger("Warning", fmt.Sprintf("Skipping entity with missing PublishedVersion or Id: %s", *entity.Id))
			continue
		}
		allFlows = append(allFlows, map[string]string{
			"title":    *entity.Name,
			"id":       *entity.Id,
			"desc":     *entity.VarType + " | " + *entity.PublishedVersion.Id,
			"version":  *entity.PublishedVersion.Id,
			"flowType": *entity.VarType,
		})
	}
	return allFlows, nil
}

func getFlowsPageCUSTOM(config *platformclientv2.Configuration, flows []string, pageNumber int) (response *FlowentitylistingCUSTOM, err error) {

	//used to get the region from the config until SDK is fixed
	region, _, _, err := genesysLogin.GenesysCreds()
	if err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("Failed to get Genesys Cloud credentials: %v", err))
	}

	url := fmt.Sprintf("https://api.%s/api/v2/flows?pageNumber=%d&pageSize=500&sortBy=asc&id=%s", region, pageNumber, strings.Join(flows, ","))
	token := config.AccessToken

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("request creation: %s", err))
		return nil, err
	}

	// Add the Authorization header with the bearer token
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("request sending: %s", err))
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		utils.TuiLogger("Error", fmt.Sprintf("Request failed with status code: %d", resp.StatusCode))
		bodyBytes, _ := io.ReadAll(resp.Body)
		utils.TuiLogger("Info", fmt.Sprintf("Response body: %s", string(bodyBytes)))
		return nil, err
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("Reading body: %s", err))
		return nil, err
	}

	// Parse the JSON response into the custom struct
	var flowResponse FlowentitylistingCUSTOM
	err = json.Unmarshal(bodyBytes, &flowResponse)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("Unmarshalling JSON: %s", err))
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
