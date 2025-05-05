package googleBotMigrate

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func BuildDigitalBot(projectId, lang, flowName, keyPath string) {
	intents, err := ListIntents(projectId, lang, keyPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Received intents: ", len(intents))

	var allVariables = ""
	var allTasks = ""
	var allIntents = ""
	var entityNameReferences = ""
	var allUtterances = ""
	var allEntities = ""
	var allEntityTypes = ""
	var allSlots = []string{}

	for _, intent := range intents {
		var displayName = intent.DisplayName

		// clean up display name with invalid characters
		if strings.Contains(displayName, " ") || strings.Contains(displayName, "-") || strings.Contains(displayName, ".") || strings.Contains(displayName, "@") {
			displayName = strings.ReplaceAll(displayName, " ", "_")
			displayName = strings.ReplaceAll(displayName, "-", "_")
			displayName = strings.ReplaceAll(displayName, ".", "_")
			displayName = strings.ReplaceAll(displayName, "@", "_")
		}
		if !strings.Contains(displayName, "Knowledge_KnowledgeBase") {
			createTask := createTask(displayName)
			allTasks += fmt.Sprintf("\n%s", createTask)
			createIntent := createIntent(displayName)
			allIntents += fmt.Sprintf("\n%s", createIntent)
			var allSegments = ""
			if len(intent.TrainingPhrases) > 0 {
				entityNameReferences = ""
				for _, trainingPhrase := range intent.TrainingPhrases {
					var segment = "            - segments:\n"
					fmt.Println("Training Phrase: ", trainingPhrase.Parts)
					for _, part := range trainingPhrase.Parts {
						// escape quotes if they are in the text
						if strings.Contains(part.Text, "\"") {
							part.Text = strings.ReplaceAll(part.Text, "\"", "\\\"")
						}
						// add text to segment
						segment += fmt.Sprintf("                - text: \"%s\"\n", part.Text)
						// if text contains entity, add entity to segment
						if part.EntityType != "" {
							// clean up entity name with invalid characters
							if strings.Contains(part.EntityType, "@") || strings.Contains(part.EntityType, ".") || strings.Contains(part.EntityType, " ") || strings.Contains(part.EntityType, "-") {
								part.EntityType = strings.ReplaceAll(part.EntityType, "@", "")
								part.EntityType = strings.ReplaceAll(part.EntityType, ".", "_")
								part.EntityType = strings.ReplaceAll(part.EntityType, " ", "_")
								part.EntityType = strings.ReplaceAll(part.EntityType, "-", "_")
							}
							segment += fmt.Sprintf(`                  entity:
                    name: %s`+"\n", part.EntityType)
							// add entity to entities if not already there as well as variables
							if !contains(allSlots, part.EntityType) {
								allSlots = append(allSlots, part.EntityType)
								allVariables += fmt.Sprintf("\n%s", createVariable(part.EntityType))
								allEntities += fmt.Sprintf("\n%s", createEntity(part.EntityType))
								allEntityTypes += fmt.Sprintf("\n%s", createEntityType(part.EntityType))
								entityNameReferences += fmt.Sprintf("\n            - %s", part.EntityType)
							}
						}
					}
					allSegments += fmt.Sprintf("\n%s", segment)
				}

			} else {
				createSegment := createSegment("No utterance")
				allSegments += fmt.Sprintf("\n%s", createSegment)
			}
			if entityNameReferences == "" {
				entityNameReferences = " []"
			}
			createUtterances := createUtterances(allSegments, entityNameReferences, displayName)
			allUtterances += fmt.Sprintf("\n%s", createUtterances)
			fmt.Println("Completed Intent: ", displayName)
		}
	}

	for _, slot := range allSlots {
		fmt.Println("Slot created: ", slot)
	}

	// Get all entities and blank out if none exist with []
	if allVariables == "" {
		allVariables = " []"
	}
	if allTasks == "" {
		allTasks = " []"
	}
	if allIntents == "" {
		allIntents = " []"
	}
	if allUtterances == "" {
		allUtterances = " []"
	}
	if allEntities == "" {
		allEntities = " []"
	}
	if allEntityTypes == "" {
		allEntityTypes = " []"
	}

	createYaml := createYaml(flowName, allVariables, allTasks, allIntents, allUtterances, allEntities, allEntityTypes)
	os.WriteFile(fmt.Sprintf("%s.yaml", flowName), []byte(createYaml), 0777)
	fmt.Println("Flow created: ", flowName)
}

func BuildKnowledgeBaseCSV(projectId, lang, fileName, keyPath string) {
	intents, err := ListIntents(projectId, lang, keyPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Received intents: ", len(intents))

	// CSV file format
	// Published, my article, this is the body, training phrase 1, training phrase 2, ...,
	allArticles := [][]string{
		{"state", "title", "textContent", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing", "phrasing"},
	}

	for _, intent := range intents {
		var displayName = intent.DisplayName
		row := []string{}

		// clean up display name with invalid characters
		if strings.Contains(displayName, " ") || strings.Contains(displayName, "-") || strings.Contains(displayName, ".") || strings.Contains(displayName, "@") {
			displayName = strings.ReplaceAll(displayName, " ", "_")
			displayName = strings.ReplaceAll(displayName, "-", "_")
			displayName = strings.ReplaceAll(displayName, ".", "_")
			displayName = strings.ReplaceAll(displayName, "@", "_")
		}
		if !strings.Contains(displayName, "Knowledge_KnowledgeBase") {
			row = append(row, "Published")
			row = append(row, displayName)

			if intent.GetMessages() != nil {
				if intent.GetMessages()[0].GetText().GetText() != nil {
					row = append(row, intent.GetMessages()[0].GetText().Text[0])
				} else {
					row = append(row, "ENTER_BODY_HERE")
				}
			} else {
				row = append(row, "ENTER_BODY_HERE")
			}

			if len(intent.TrainingPhrases) > 0 {

				for _, trainingPhrase := range intent.TrainingPhrases {
					fmt.Println("Training Phrase: ", trainingPhrase.Parts)
					var phrase = ""
					for _, part := range trainingPhrase.Parts {
						// escape quotes if they are in the text
						if strings.Contains(part.Text, ",") {
							part.Text = strings.ReplaceAll(part.Text, ",", "")
						}
						if strings.Contains(part.Text, "\"") {
							part.Text = strings.ReplaceAll(part.Text, "\"", "'")
						}
						phrase += part.Text
					}
					row = append(row, phrase)
				}

			} else {
				fmt.Println("No utterance")
			}
			allArticles = append(allArticles, row)
			fmt.Println("Completed Intent: ", displayName)
		}
	}

	fmt.Println("building csv")
	//os.WriteFile(fmt.Sprintf("%s.csv", fileName), []byte(allArticles), 0777)
	csvExport(allArticles, fileName)
	fmt.Println("CSV File created: ", fileName)
}

func csvExport(data [][]string, name string) error {
	file, err := os.Create(fmt.Sprintf("%s.csv", name))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		if err := writer.Write(value); err != nil {
			return err
		}
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ListIntents(projectID, lang, keyPath string) ([]*dialogflowpb.Intent, error) {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(keyPath))
	if clientErr != nil {
		return nil, clientErr
	}
	defer intentsClient.Close()

	if projectID == "" {
		return nil, fmt.Errorf("received empty project (%s)", projectID)
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)

	request := dialogflowpb.ListIntentsRequest{Parent: parent, IntentView: dialogflowpb.IntentView_INTENT_VIEW_FULL, LanguageCode: lang}

	intentIterator := intentsClient.ListIntents(ctx, &request)
	var intents []*dialogflowpb.Intent

	for intent, status := intentIterator.Next(); status != iterator.Done; {
		if len(intents) > 1000 {
			fmt.Println("Error: Does your api key have API admin access??")
			os.Exit(1)
		}
		intents = append(intents, intent)
		intent, status = intentIterator.Next()
	}

	return intents, nil
}

func createVariable(name string) string {
	// Create a new variable
	var variable = fmt.Sprintf(`    - stringVariable:
        name: Slot.%s
        initialValue:
          noValue: true
        isInput: true
        isOutput: true`, name)
	return variable
}

func createTask(name string) string {
	// Create a new task
	var task = fmt.Sprintf(`    - task:
        name: %s
        actions:
          - exitBotFlow:
              name: Exit Bot Flow`, name)
	return task
}

func createIntent(name string) string {
	// Create a new intent
	var intent = fmt.Sprintf(`      - intent:
          confirmation:
            exp: "MakeCommunication(\n  ToCommunication(ToCommunication(\"I think you want to %s, is that correct?\")))"
          name: %s
          task:
            name: %s`, name, name, name)
	return intent
}

func createSegment(text string) string {
	var segment = fmt.Sprintf(`            - segments:
                - text: %s`, text)
	return segment
}

func createEntity(name string) string {
	var entity = fmt.Sprintf(`        - name: %s
          type: %sType`, name, name)
	return entity
}

func createEntityType(name string) string {
	var entityType = fmt.Sprintf(`        - name: %sType
          description: The description of the Entity Type "%sType"
          mechanism:
            type: Regex
            restricted: true
            items: []`, name, name)
	return entityType
}

func createUtterances(segments, entityNameReferences, name string) string {
	// Create a new utterance
	var utterance = fmt.Sprintf(`        - utterances:%s
              source: User
          entityNameReferences:%s
          name: %s`, segments, entityNameReferences, name)
	return utterance
}

func createYaml(flowName, variables, task, intent, utterance, entities, entityTypes string) string {
	// Create the yaml file
	var yaml = fmt.Sprintf(`digitalBot:
  name: %s
  defaultLanguage: en-au
  startUpRef: "/digitalBot/bots/bot[Initial Greeting_10]"
  bots:
    - bot:
        name: Initial Greeting
        refId: Initial Greeting_10
        actions:
          - waitForInput:
              name: Wait for Input
              question:
                exp: "MakeCommunication(\n  ToCommunication(ToCommunication(\"What would you like to do?\")))"
              knowledgeSearchResult:
                noValue: true
              noMatch:
                exp: "MakeCommunication(\n  ToCommunication(ToCommunication(\"Tell me again what you would like to do.\")))"
  variables:%s
  tasks:%s
  settingsBotFlow:
    intentSettings:%s
  settingsNaturalLanguageUnderstanding:
    nluDomainVersion:
      intents:%s
      entities:%s
      entityTypes:%s
      language: en-au
      languageVersions: {}
    mutedUtterances: []
`, flowName, variables, task, intent, utterance, entities, entityTypes)
	return yaml
}
