package ttsChanger

import (
	"fmt"
	"os"
	"os/exec"

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

func CurrentTTSVoices(config *platformclientv2.Configuration, engine string) []map[string]string {

	var voices []map[string]string

	apiIntegrationsInstance := platformclientv2.NewIntegrationsApiWithConfig(config)
	ttsVoices, err := getTTSVoicesEnhancedPages(apiIntegrationsInstance, 1, engine)

	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) Error getting tts voices: %v", err))
		os.Exit(1)
	}

	pageNumber := 2
	for pageNumber <= *ttsVoices.PageCount {
		nextPage, err := getTTSVoicesEnhancedPages(apiIntegrationsInstance, pageNumber, engine)
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

func CurrentTTSEngines(config *platformclientv2.Configuration) []map[string]string {

	var engines []map[string]string

	apiIntegrationsInstance := platformclientv2.NewIntegrationsApiWithConfig(config)
	ttsEngines, err := getTTSEnginesPages(apiIntegrationsInstance, 1)

	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) Error getting tts engines: %v", err))
		os.Exit(1)
	}

	pageNumber := 2
	for pageNumber <= *ttsEngines.PageCount {
		nextPage, err := getTTSEnginesPages(apiIntegrationsInstance, pageNumber)
		if err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("(ttsChanger) %s", err))
			os.Exit(1)
		}
		*ttsEngines.Entities = append(*ttsEngines.Entities, *nextPage.Entities...)
		pageNumber++
	}

	for _, entity := range *ttsEngines.Entities {
		engines = append(engines, map[string]string{
			"title": *entity.Name,
			"id":    *entity.Id,
			"desc":  fmt.Sprintf("This TTS Engine supports %d languages", len(*entity.Languages)),
		})
	}

	return engines
}

func getTTSEnginesPages(apiInstance *platformclientv2.IntegrationsApi, page int) (ttsEngines *platformclientv2.Ttsengineentitylisting, err error) {
	pageSize := 100        // Page size
	includeVoices := false // Include voices for the engine
	name := ""             // Filter on engine name
	language := ""         // Filter on supported language. If includeVoices=true then the voices are also filtered.
	// Get a list of TTS engines enabled for org
	data, _, err := apiInstance.GetIntegrationsSpeechTtsEngines(page, pageSize, includeVoices, name, language)
	if err != nil {
		return nil, fmt.Errorf("Error calling GetIntegrationsSpeechTtsEngines: %v\n", err)
	} else {
		return data, nil
	}
}

func getTTSVoicesEnhancedPages(apiInstance *platformclientv2.IntegrationsApi, page int, engine string) (ttsVoices *platformclientv2.Ttsvoiceentitylisting, err error) {
	engineId := engine // The engine ID
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
