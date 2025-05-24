package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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
var bannerStyleLoading = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).Background(lipgloss.Color("#655ad5")).Padding(0, 1).Render
var bannerWarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#f7d720")).Padding(0, 1)

type item struct {
	title, desc, typeSelected, id string
}

func (i item) Title() string        { return i.title }
func (i item) Description() string  { return i.desc }
func (i item) Id() string           { return i.id }
func (i item) TypeSelected() string { return i.typeSelected }
func (i item) FilterValue() string  { return i.title + " " + i.desc }

type model struct {
	list             list.Model
	lastSelectedItem item
	spinner          spinner.Model
	spinning         bool
}

type ttsSelection struct {
	ttsEngineGet     string
	ttsEngineGetName string
	ttsEngineSet     string
	ttsEngineSetName string
	ttsGet           string
	ttsAPIGet        string
	ttsSet           string
	flows            []map[string]string
	downloadedDir    string
	updatedDir       string
}

// --- Message Types for Async Operations ---
type searchResultsMsg struct {
	results []map[string]string
	err     error
}
type ttsEnginesMsg struct {
	engines []map[string]string
	err     error
}
type ttsVoicesMsg struct {
	voices []map[string]string
	err    error
}
type flowsResultMsg struct { // Renamed to avoid conflict with flows package
	flowsData []map[string]string
	err       error
}

// --- Cmds for Async Operations ---
func fetchSearchResultsCmd(query string) tea.Cmd {
	return func() tea.Msg {
		results := searchReleaseNotes.SearchReleaseNotes(query)
		return searchResultsMsg{results: results, err: nil}
	}
}
func fetchTTSEnginesCmd(config *platformclientv2.Configuration) tea.Cmd {
	return func() tea.Msg {
		engines := ttsChanger.CurrentTTSEngines(config)
		return ttsEnginesMsg{engines: engines, err: nil}
	}
}
func fetchTTSVoicesCmd(config *platformclientv2.Configuration, engineID string) tea.Cmd {
	return func() tea.Msg {
		voices := ttsChanger.CurrentTTSVoices(config, engineID)
		return ttsVoicesMsg{voices: voices, err: nil}
	}
}
func fetchFlowsCmd(config *platformclientv2.Configuration, searchType string, searchValue string) tea.Cmd {
	return func() tea.Msg {
		flowsResult := flows.GetFlows(config, searchType, searchValue)
		return flowsResultMsg{flowsData: flowsResult, err: nil}
	}
}

func fetchFlowsCUSTOMCmd(config *platformclientv2.Configuration, flowTypes []string) tea.Cmd {
	return func() tea.Msg {
		flowsResult, err := flows.GetFlowsCUSTOM(config, flowTypes)
		return flowsResultMsg{flowsData: flowsResult, err: err}
	}
}

var ttsData ttsSelection
var genesysLoginConfig *platformclientv2.Configuration
var orgName string

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "ctrl+e" {
			var csvData [][]string
			csvData = append(csvData, []string{"Title", "Description", "Type", "ID"})
			for _, row := range m.list.VisibleItems() {
				csvData = append(csvData, []string{row.(item).Title(), row.(item).Description(), row.(item).TypeSelected(), row.(item).Id()})
			}
			err := utils.ExportToCSV(csvData, fmt.Sprintf("export_%s.csv", fmt.Sprint(time.Now().Unix())))
			if err != nil {
				utils.TuiLogger("Error", fmt.Sprintf("Export Error: %s", err))
			}
			m.list.Title = "Exported to CSV"
		}
		if msg.String() == "enter" {
			if m.list.FilterState() == list.Filtering {
				// Let the list handle the filter input first
				break
			}
			selected := m.list.SelectedItem().(item)
			m.lastSelectedItem = selected
			m.list.Styles.Title = bannerStyle
			switch selected.typeSelected {
			// TTS section
			case "ttsEngineGet":
				if orgName == "not provided" {
					m.list.Title = "WARNING: you need to connect to an ORG first"
					m.list.Styles.Title = bannerWarningStyle
				} else {
					m.spinning = true
					m.list.Title = "Loading TTS Engines..."
					cmds = append(cmds, m.spinner.Tick, fetchTTSEnginesCmd(genesysLoginConfig))
				}
			case "ttsSelectVoiceGet":
				ttsData.ttsEngineGet = selected.id
				ttsData.ttsEngineGetName = selected.title
				m.spinning = true
				m.list.Title = "Loading TTS Voices..."
				cmds = append(cmds, m.spinner.Tick, fetchTTSVoicesCmd(genesysLoginConfig, ttsData.ttsEngineGet))
			case "ttsEngineSet":
				ttsData.ttsGet = selected.Title()
				ttsData.ttsAPIGet = selected.id
				m.spinning = true
				m.list.Title = "Loading TTS Engines for Set..."
				cmds = append(cmds, m.spinner.Tick, fetchTTSEnginesCmd(genesysLoginConfig))
			case "ttsSelectVoiceSet":
				ttsData.ttsEngineSet = selected.id
				ttsData.ttsEngineSetName = selected.title
				m.spinning = true
				m.list.Title = "Loading TTS Voices for Set..."
				cmds = append(cmds, m.spinner.Tick, fetchTTSVoicesCmd(genesysLoginConfig, ttsData.ttsEngineSet))
			case "ttsSet":
				ttsData.ttsSet = selected.Title()
				m.spinning = true
				m.list.Title = "Loading flows with " + ttsData.ttsGet + "..."
				cmds = append(cmds, m.spinner.Tick, fetchFlowsCmd(genesysLoginConfig, "TTSVOICE", ttsData.ttsEngineGet+"/"+ttsData.ttsAPIGet))
			case "flowUpdate":
				m.list.Title = "Updating..."
				return m, tea.Quit
			// PWA section
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
			// flow backup section
			case "flowBackupSelect":
				if orgName == "not provided" {
					m.list.Title = "WARNING: you need to connect to an ORG first"
					m.list.Styles.Title = bannerWarningStyle
				} else {
					m.spinning = true
					m.list.Title = "Loading flows for backup..."
					cmds = append(cmds, m.spinner.Tick, fetchFlowsCUSTOMCmd(genesysLoginConfig, []string{}))
				}
			case "flowBackup":
				m.list.Title = "Backing up..."
				return m, tea.Quit
			// common module section
			case "commonModule":
				if orgName == "not provided" {
					m.list.Title = "WARNING: you need to connect to an ORG first"
					m.list.Styles.Title = bannerWarningStyle
				} else {
					m.spinning = true
					m.list.Title = "Loading Common Modules..."
					cmds = append(cmds, m.spinner.Tick, fetchFlowsCUSTOMCmd(genesysLoginConfig, []string{}))
				}
			case "commonDependency":
				ttsData.ttsSet = selected.title
				m.spinning = true
				m.list.Title = "Loading dependencies for " + selected.title + "..."
				cmds = append(cmds, m.spinner.Tick, fetchFlowsCmd(genesysLoginConfig, "commonModuleFlow", selected.id))
			case "commonRefresh":
				m.list.Title = "RePublishing Common Modules"
				return m, tea.Quit
			// search release notes section
			case "searchReleases":
				m.spinning = true
				m.list.Title = "Searching Release Notes..."
				cmds = append(cmds, m.spinner.Tick, fetchSearchResultsCmd(" "))
			case "backMain":
				m.list.Title = "McPhee11 TUI - Genesys Cloud ORG: " + orgName
				m.list.SetItems(menuMain())
				m.list.ResetFilter()
			case "help":
				m.list.Title = "Help Menu"
				m.list.SetItems(menuHelp())
				m.list.ResetFilter()
			case "version":
				currentVersion := utils.GetVersion()
				laterVersion, newerVersion, err := utils.CheckForNewerVersion(currentVersion)
				if err != nil {
					utils.TuiLogger("Info", fmt.Sprintf("Unable to check for new version: %s", err))
					m.list.Title = currentVersion + " (update check failed)"
				} else {
					if laterVersion {
						m.list.Title = fmt.Sprintf("%s - Newer version: %s. Update with: go install github.com/mcphee11/mcphee11-tui@latest", currentVersion, newerVersion)
					} else {
						m.list.Title = currentVersion + " (up to date)"
					}
				}
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

		// --- Handle Async Operation Results ---
	case searchResultsMsg:
		m.spinning = false
		if msg.err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("Search Error: %v", msg.err))
			m.list.Title = "Error searching release notes!"
		} else {
			m.list.Title = "Release notes"
			m.list.SetItems(menuSearchReleaseNotes(msg.results))
			m.list.Filter = utils.CustomSubstringFilter
			m.list.ResetFilter()
		}

	case ttsEnginesMsg:
		m.spinning = false
		if msg.err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("TTS Engines Load Error: %v", msg.err))
			m.list.Title = "Error Loading TTS Engines!"
		} else {
			switch m.lastSelectedItem.typeSelected {
			case "ttsEngineGet":
				m.list.Title = "Select the TTS Engine you want to replace"
				m.list.SetItems(menuCurrentTTSVoices(msg.engines, "ttsSelectVoiceGet"))
			case "ttsEngineSet":
				m.list.Title = "Select the TTS Engine you want to SET"
				m.list.SetItems(menuCurrentTTSVoices(msg.engines, "ttsSelectVoiceSet"))
			}
			m.list.ResetFilter()
		}

	case ttsVoicesMsg:
		m.spinning = false
		if msg.err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("TTS Voices Load Error: %v", msg.err))
			m.list.Title = "Error Loading TTS Voices!"
		} else {
			switch m.lastSelectedItem.typeSelected {
			case "ttsSelectVoiceGet":
				m.list.Title = "Select the Voice you want to replace"
				m.list.SetItems(menuCurrentTTSVoices(msg.voices, "ttsEngineSet"))
			case "ttsSelectVoiceSet":
				m.list.Title = "Select the Voice you want to SET"
				m.list.SetItems(menuCurrentTTSVoices(msg.voices, "ttsSet"))
			}
			m.list.ResetFilter()
		}
	case flowsResultMsg:
		m.spinning = false
		if msg.err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("Flows Load Error: %v", msg.err))
			m.list.Title = "Error Loading Flows!"
		} else {
			utils.TuiLogger("Info", "flowsResultMsg!!!!!")
			ttsData.flows = msg.flowsData
			switch m.lastSelectedItem.typeSelected {
			case "ttsSet":
				m.list.Title = "These flows include " + ttsData.ttsGet + ". Update one or ALL"
				m.list.SetItems(menuCurrentFlows(msg.flowsData, "flowUpdate"))
			case "flowBackupSelect":
				m.list.Title = "Select the flow you want to backup"
				m.list.SetItems(menuCurrentFlows(msg.flowsData, "flowBackup"))
			case "commonModule":
				var commonModules []map[string]string
				for _, flow := range msg.flowsData {
					if flow["flowType"] == "COMMONMODULE" {
						commonModules = append(commonModules, flow)
					}
				}
				ttsData.flows = commonModules
				m.list.Title = "Select Common Module flow that you have updated"
				m.list.SetItems(menuCurrentFlows(commonModules, "commonDependency"))
			case "commonDependency":
				m.list.Title = "Select from one or ALL of the flows that depend on " + ttsData.ttsSet
				m.list.SetItems(menuCurrentFlows(msg.flowsData, "commonRefresh"))
			}
			m.list.ResetFilter()
		}

	case spinner.TickMsg:
		var var_cmd tea.Cmd // var_cmd to avoid conflict with package
		m.spinner, var_cmd = m.spinner.Update(msg)
		cmds = append(cmds, var_cmd)
	}

	if !m.spinning {
		updatedList, cmd := m.list.Update(msg)
		m.list = updatedList
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.spinning {
		title := fmt.Sprintf("%s%s", m.spinner.View(), bannerStyleLoading(" "+m.list.Title))
		m.list.SetShowTitle(false)
		return docStyle.Render("\n\n\n" + title + m.list.View())
	}
	return docStyle.Render(m.list.View())
}

func customKeys() []key.Binding {
	return []key.Binding{key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "export to csv"))}
}

func main() {
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

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#3c71a8"))

	m := model{list: list.New(menuMain(), list.NewDefaultDelegate(), 0, 0), spinner: s, spinning: false}
	m.list.Title = "McPhee11 TUI - Genesys Cloud ORG: " + orgName
	m.list.AdditionalFullHelpKeys = customKeys

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
		flows.FlowsLoadingMainBackup(returnedModel.(model).lastSelectedItem.id, ttsData.flows, ttsData.ttsGet, ttsData.ttsSet, ttsData.ttsEngineGetName, ttsData.ttsEngineSetName, "tts", true)
	case "flowBackup":
		flows.FlowsLoadingMainBackup(returnedModel.(model).lastSelectedItem.id, ttsData.flows, "", "", "", "", "tts", false)
	case "commonRefresh":
		flows.FlowsLoadingMainBackup(returnedModel.(model).lastSelectedItem.id, ttsData.flows, "", ttsData.ttsSet, "", "", "rePublish", true)
	}
}

func menuMain() []list.Item {
	return []list.Item{
		item{typeSelected: "searchReleases", title: "Search Release Notes", desc: "Search the Genesys Cloud Release Notes"},
		item{typeSelected: "pwaBanking", title: "Build Banking PWA", desc: "Building a PWA mobile app for demos based on banking"},
		item{typeSelected: "ttsEngineGet", title: "Update TTS", desc: "Update the TTS engine used in your Flows"},
		item{typeSelected: "commonModule", title: "Common Modules", desc: "Update the flows that have a specifc common module set"},
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
