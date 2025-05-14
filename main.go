package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mcphee11/mcphee11-tui/flows"
	"github.com/mcphee11/mcphee11-tui/genesysLogin"
	"github.com/mcphee11/mcphee11-tui/googleBotMigrate"
	"github.com/mcphee11/mcphee11-tui/pwaBanking"
	"github.com/mcphee11/mcphee11-tui/searchReleaseNotes"
	"github.com/mcphee11/mcphee11-tui/ttsChanger"
	"github.com/mcphee11/mcphee11-tui/utils"
	"github.com/mypurecloud/platform-client-sdk-go/platformclientv2"
)

var Debug bool
var docStyle = lipgloss.NewStyle().Margin(1, 2)
var bannerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).Background(lipgloss.Color("#655ad5")).Padding(0, 1)
var bannerWarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#f7d720")).Padding(0, 1)

type item struct {
	title, desc, typeSelected, id string
}

func (i item) Title() string        { return i.title }
func (i item) Description() string  { return i.desc }
func (i item) Id() string           { return i.id }
func (i item) TypeSelected() string { return i.typeSelected }
func (i item) FilterValue() string  { return i.title }

type model struct {
	list             list.Model
	lastSelectedItem item
}

type ttsSelection struct {
	ttsGet        string
	ttsAPIGet     string
	ttsSet        string
	flows         []map[string]string
	downloadedDir string
	updatedDir    string
}

var ttsData ttsSelection
var genesysLoginConfig *platformclientv2.Configuration
var orgName string

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			selected := m.list.SelectedItem().(item)
			m.lastSelectedItem = selected
			m.list.Styles.Title = bannerStyle
			switch selected.typeSelected {
			// TTS
			case "ttsChanger":
				if orgName == "not provided" {
					m.list.Title = "WARNING: you need to connect to an ORG first"
					m.list.Styles.Title = bannerWarningStyle
				} else {
					m.list.Title = "Select the Voice you want to replace"
					voices := ttsChanger.CurrentTTSVoices(genesysLoginConfig)
					m.list.SetItems(menuCurrentTTSVoices(voices, "ttsGet"))
					m.list.Cursor()
				}
			case "ttsGet":
				ttsData.ttsGet = selected.Title()
				ttsData.ttsAPIGet = selected.id
				m.list.Title = "Select the Voice you want to SET"
				voices := ttsChanger.CurrentTTSVoices(genesysLoginConfig)
				m.list.SetItems(menuCurrentTTSVoices(voices, "ttsSet"))
				m.list.ResetFilter()
				m.list.Cursor()
			case "ttsSet":
				ttsData.ttsSet = selected.Title()
				m.list.Title = "These flows include " + ttsData.ttsGet + ". Update one or ALL"
				flows := ttsChanger.GetFlows(genesysLoginConfig, "genesys_enhanced/"+ttsData.ttsAPIGet)
				ttsData.flows = flows
				m.list.SetItems(menuCurrentFlows(flows, "flowUpdate"))
				m.list.ResetFilter()
				m.list.Cursor()
			case "flowUpdate":
				m.list.Title = "Updating..."
				return m, tea.Quit
			// PWA
			case "pwaBanking":
				m.list.Title = "Building Banking PWA"
				return m, tea.Quit
			case "botMigrate":
				if orgName == "not provided" {
					m.list.Title = "WARNING: you need to connect to an ORG first"
					m.list.Styles.Title = bannerWarningStyle
				} else {
					m.list.Title = "Migrating Google Bots"
					return m, tea.Quit
				}
			case "flowBackupSelect":
				if orgName == "not provided" {
					m.list.Title = "WARNING: you need to connect to an ORG first"
					m.list.Styles.Title = bannerWarningStyle
				} else {
					m.list.Title = "Select the flow you want to backup"
					justFlows, err := ttsChanger.GetFlowsCUSTOM(genesysLoginConfig, []string{})
					if err != nil {
						utils.TuiLogger("Error", fmt.Sprintf("%s", err))
					}
					ttsData.flows = justFlows
					m.list.SetItems(menuCurrentFlows(justFlows, "flowBackup"))
				}
			case "flowBackup":
				m.list.Title = "Backing up..."
				return m, tea.Quit
			case "searchReleases":
				m.list.StartSpinner()
				search := searchReleaseNotes.SearchReleaseNotes(" ")
				m.list.Title = "Release notes"
				m.list.SetItems(menuSearchReleaseNotes(search))
				m.list.Cursor()
				// Back to main menu
			case "backMain":
				m.list.Title = "McPhee11 TUI - Genesys Cloud ORG: " + orgName
				m.list.SetItems(menuMain())
				m.list.Cursor()
			case "help":
				m.list.Title = "Help Menu"
				m.list.SetItems(menuHelp())
				m.list.Cursor()
			case "version":
				m.list.Title = utils.GetVersion()
			case "link":
				err := openURL(selected.id)
				if err != nil {
					utils.TuiLogger("Error", fmt.Sprintf("Open URL Error: %s", err))
				}
			default:
				m.list.Title = "DEFAULT"
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	// initialize logger
	err := utils.TuiLoggerStart()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Check for genesys cloud environment
	utils.TuiLogger("Info", "(main) Checking for Genesys Cloud environment variables...")
	config, err := genesysLogin.GenesysLogin()
	if err != nil {
		orgName = "not provided"
		utils.TuiLogger("Info", "(main) ORG not provided")
	} else {
		genesysLoginConfig = config
		apiInstance := platformclientv2.NewOrganizationApiWithConfig(config)
		data, _, err := apiInstance.GetOrganizationsMe()
		if err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("(main) Failed calling GetOrganizationsMe: %s", err))
		} else {
			orgName = *data.Name
		}
	}

	m := model{list: list.New(menuMain(), list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "McPhee11 TUI - Genesys Cloud ORG: " + orgName

	p := tea.NewProgram(m, tea.WithAltScreen())

	returnedModel, err := p.Run()
	if err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("(main) Error running program: %s", err))
	}

	switch returnedModel.(model).lastSelectedItem.typeSelected {
	case "pwaBanking":
		pwaBanking.MainInputs()
	case "botMigrate":
		googleBotMigrate.MainInputs()
	case "flowUpdate":
		flows.FlowsLoadingMainBackup(returnedModel.(model).lastSelectedItem.id, ttsData.flows, ttsData.ttsGet, ttsData.ttsSet, true)
	case "flowBackup":
		flows.FlowsLoadingMainBackup(returnedModel.(model).lastSelectedItem.id, ttsData.flows, "", "", false)
	}
}

func menuMain() []list.Item {
	return []list.Item{
		item{typeSelected: "searchReleases", title: "Search Release Notes", desc: "Search the Genesys Cloud Release Notes"},
		item{typeSelected: "pwaBanking", title: "Build Banking PWA", desc: "Building a PWA mobile app for demos based on banking"},
		item{typeSelected: "ttsChanger", title: "Update to Genesys Enhanced TTS", desc: "Update the TTS engine used in your Genesys Voice BOTs"},
		item{typeSelected: "botMigrate", title: "Google Bot Migration", desc: "Easily migrate Google Bots (ES & CX) to Digital Bots or Knowledge Base for Copilot"},
		item{typeSelected: "flowBackupSelect", title: "Backup Flows", desc: "Take a backup of your Genesys Flows"},
		item{typeSelected: "help", title: "Help Menu", desc: "Open the help menu"},
		item{typeSelected: "version", title: "Version", desc: "Display installed version"},
	}
}

func menuSearchReleaseNotes(search []map[string]string) []list.Item {
	var list []list.Item
	for i := range search {
		if search[i]["link"] == "" {
			continue
		}
		list = append(list, item{typeSelected: "link", id: search[i]["link"], title: search[i]["notes"], desc: search[i]["section"]})
	}
	list = append(list, item{typeSelected: "backMain", id: "backMain", title: "Back", desc: "Back to the previous menu"})
	return list
}

func menuCurrentTTSVoices(voices []map[string]string, ttsType string) []list.Item {
	var list []list.Item
	for i := range voices {
		list = append(list, item{typeSelected: ttsType, id: voices[i]["id"], title: voices[i]["title"], desc: voices[i]["desc"]})
	}
	if ttsType == "ttsSet" {
		list = append(list, item{typeSelected: ttsType, id: "Default", title: "Default", desc: "Default voice thats configured"})
	}
	list = append(list, item{typeSelected: "backMain", id: "backMain", title: "Back", desc: "Back to the main menu"})
	return list
}

func menuCurrentFlows(flows []map[string]string, flowId string) []list.Item {
	var list []list.Item
	if len(flows) == 0 {
		utils.TuiLogger("Info", fmt.Sprintf("No flows found Published with TTS: %s", ttsData.ttsGet))
		list = append(list, item{typeSelected: "backMain", id: "backMain", title: "Back", desc: "Back to the main menu"})
		return list
	}
	list = append(list, item{typeSelected: flowId, id: "ALL", title: "ALL", desc: "Update all the flows with " + ttsData.ttsGet})
	for i := range flows {
		list = append(list, item{typeSelected: flowId, id: flows[i]["id"], title: flows[i]["title"], desc: flows[i]["desc"]})
	}
	list = append(list, item{typeSelected: "backMain", id: "backMain", title: "Back", desc: "Back to the main menu"})
	utils.TuiLogger("Info", fmt.Sprintf("%d flows found Published with TTS %s", len(flows), ttsData.ttsGet))
	return list
}

func menuHelp() []list.Item {
	return []list.Item{
		item{typeSelected: "link", id: "https://github.com/mcphee11/mcphee11-tui", title: "GitHub repo", desc: "Open the GitHub repository"},
		item{typeSelected: "link", id: "https://help.mypurecloud.com", title: "Genesys Cloud Help", desc: "Open the Genesys Cloud Help Center"},
		item{typeSelected: "link", id: "https://developer.genesys.cloud/", title: "Genesys Cloud Developer Center", desc: "Open the Genesys Cloud Developer website"},
		item{typeSelected: "backMain", id: "backMain", title: "Back", desc: "Back to the previous menu"},
	}
}

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	return cmd.Start()
}
