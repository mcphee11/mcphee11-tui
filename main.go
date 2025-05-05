package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mcphee11/mcphee11-tui/genesysLogin"
	"github.com/mcphee11/mcphee11-tui/googleBotMigrate"
	"github.com/mcphee11/mcphee11-tui/pwaBanking"
	"github.com/mcphee11/mcphee11-tui/searchReleaseNotes"
	"github.com/mcphee11/mcphee11-tui/ttsChanger"
	"github.com/mcphee11/mcphee11-tui/utils"
	"github.com/mypurecloud/platform-client-sdk-go/platformclientv2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc, id string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) Id() string          { return i.id }
func (i item) FilterValue() string { return i.title }

type model struct {
	list             list.Model
	lastSelectedItem item
}

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
			m.lastSelectedItem.id = selected.id

			switch selected.id {
			// TTS
			case "ttsChanger":
				m.list.Title = "Select the Language you need to use"
				m.list.SetItems(menuTTSLanguage())
				m.list.Cursor()
			case "en-AU":
				m.list.Title = "Select the Amazon Polly Voice you need to use"
				m.list.SetItems(menuTTSenAU())
				m.list.Cursor()
			case "setTTS,en-AU":
				m.list.Title = "Setting the TTS engine"
				return m, tea.Quit
			case "backTTSLanguage":
				m.list.Title = "Select the Language you need to use"
				m.list.SetItems(menuTTSLanguage())
				m.list.Cursor()
			// PWA
			case "pwaBanking":
				m.list.Title = "Building Banking PWA"
				return m, tea.Quit
			case "botMigrate":
				m.list.Title = "Migrating Google Bots"
				return m, tea.Quit
			case "flowBackup":
				fmt.Println("Backing up flows")
			case "searchReleases":
				m.list.StartSpinner()
				search := searchReleaseNotes.SearchReleaseNotes(" ")
				m.list.Title = "Release notes"
				m.list.SetItems(menuSearchReleaseNotes(search))
				m.list.Cursor()
				// Back to main menu
			case "backMain":
				m.list.Title = "McPhee11 TUI for making life easier"
				m.list.SetItems(menuMain())
				m.list.Cursor()
			case "help":
				m.list.Title = "Help Menu"
				m.list.SetItems(menuHelp())
				m.list.Cursor()
			case "version":
				m.list.Title = utils.GetVersion()
			default:
				if strings.Contains(selected.id, "https://") {
					err := openURL(selected.id)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
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
	var org = "unknown"
	// Check for genesys cloud environment
	config, err := genesysLogin.GenesysLogin()
	if err != nil {
		org = "not provided"
	} else {
		apiInstance := platformclientv2.NewOrganizationApiWithConfig(config)
		data, _, err := apiInstance.GetOrganizationsMe()
		if err != nil {
			fmt.Printf("Error calling GetOrganizationsMe: %v\n", err)
		} else {
			org = *data.Name
		}
	}

	m := model{list: list.New(menuMain(), list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "McPhee11 TUI - Genesys Cloud ORG: " + org

	p := tea.NewProgram(m, tea.WithAltScreen())

	returnedModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}

	switch returnedModel.(model).lastSelectedItem.id {
	case "pwaBanking":
		pwaBanking.MainInputs()
	case "setTTS,en-AU":
		ttsChanger.TTSChanger(returnedModel.(model).lastSelectedItem.id)
	case "botMigrate":
		googleBotMigrate.MainInputs()
	}
}

func menuMain() []list.Item {
	return []list.Item{
		item{id: "searchReleases", title: "Search Release Notes", desc: "Search the Genesys Cloud Release Notes"},
		item{id: "pwaBanking", title: "Build Banking PWA", desc: "Building a PWA mobile app for demos based on banking"},
		item{id: "ttsChanger", title: "Update to Genesys Enhanced TTS", desc: "Update the TTS engine used in your Genesys Voice BOTs"},
		item{id: "botMigrate", title: "Google Bot Migration", desc: "Easily migrate Google Bots (ES & CX) to Genesys Digital Bots"},
		item{id: "flowBackup", title: "Backup Flows", desc: "Take a backup of your Genesys Flows"},
		item{id: "help", title: "Help Menu", desc: "Open the help menu"},
		item{id: "version", title: "Version", desc: "Display installed version"},
	}
}

func menuTTSLanguage() []list.Item {
	return []list.Item{
		item{id: "en-AU", title: "en-AU", desc: "English (Australia)"},
		item{id: "en-GB", title: "en-GB", desc: "English (United Kingdom)"},
		item{id: "en-US", title: "en-US", desc: "English (United States)"},
		item{id: "es-ES", title: "es-ES", desc: "Spanish (Spain)"},
		item{id: "backMain", title: "Back", desc: "Back to the previous menu"},
	}
}

func menuTTSenAU() []list.Item {
	return []list.Item{
		item{id: "setTTS,en-AU", title: "Olivia", desc: "Female"},
		item{id: "backTTSLanguage", title: "Back", desc: "Back to the previous menu"},
	}
}

func menuSearchReleaseNotes(search []map[string]string) []list.Item {
	var list []list.Item
	for i := range search {
		if search[i]["link"] == "" {
			continue
		}
		list = append(list, item{id: search[i]["link"], title: search[i]["notes"], desc: search[i]["section"]})
	}
	list = append(list, item{id: "backMain", title: "Back", desc: "Back to the previous menu"})
	return list
}

func menuHelp() []list.Item {
	return []list.Item{
		item{id: "https://github.com/mcphee11/mcphee11-tui", title: "GitHub repo", desc: "Open the GitHub repository"},
		item{id: "https://help.mypurecloud.com", title: "Genesys Cloud Help", desc: "Open the Genesys Cloud Help Center"},
		item{id: "https://developer.genesys.cloud/", title: "Genesys Cloud Developer Center", desc: "Open the Genesys Cloud Developer website"},
		item{id: "backMain", title: "Back", desc: "Back to the previous menu"},
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
