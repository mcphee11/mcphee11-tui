package pwaDeploy

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mcphee11/mcphee11-tui/utils"
)

//go:embed _pwaTemplates/*
var pwaTemplates embed.FS

func buildBankingPwa(flagName, flagShortName, flagColor, flagIcon, flagBanner, flagRegion, flagEnvironment, flagDeploymentId, flagBucketName string, p *tea.Program) {

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
	// ------------ Create project folder -----------------
	err := os.Mkdir(flagShortName, 0777)
	if err != nil {
		sendStatusUpdate("Error", fmt.Sprintf("Error creating directory %s, exiting build.", flagShortName))
		sendMsgToUI(stage1CompleteMsg{})
		return
	}

	err = os.Mkdir(fmt.Sprintf("%s/svgs", flagShortName), 0777)
	if err != nil {
		sendStatusUpdate("Error", fmt.Sprintf("Error creating directory %s/svgs, exiting build.", flagShortName))
		sendMsgToUI(stage1CompleteMsg{})
		return
	}
	utils.TuiLogger("Info", fmt.Sprintf("(buildBankingPwa) Created folder: %s", flagShortName))

	// ------------------ create svgs ------------------
	svgs, err := pwaTemplates.ReadDir("_pwaTemplates/svgs")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) Error reading svgs dir: %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	for i := 0; i < len(svgs); i++ {
		err := createFile(svgs[i].Name(), fmt.Sprintf("%s/svgs", flagShortName), fmt.Sprintf("_pwaTemplates/svgs/%s", svgs[i].Name()))
		if err != nil {
			utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) Error reading svgs dir: %s", err))
			_ = os.RemoveAll(flagShortName)
			return
		}
		sendStatusUpdate("Info", fmt.Sprintf("Generated Svg: %s", svgs[i].Name()))
	}
	sendMsgToUI(flowProcessedMsg{})

	// ------------------ create icons ------------------
	sendStatusUpdate("Info", "Generating App Icons this can take a min so please wait...")
	icons, err := pwaTemplates.ReadFile("_pwaTemplates/icons.sh")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedIcons := strings.ReplaceAll(string(icons), "$icon", flagIcon)
	err = os.WriteFile(fmt.Sprintf("%s/icons.sh", flagShortName), []byte(formattedIcons), 0777)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	currentDir, err := os.Getwd()
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) error getting working dir: %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	cmdIcon := exec.Command("./icons.sh")
	cmdIcon.Dir = fmt.Sprintf("%s/%s", currentDir, flagShortName)

	if err := cmdIcon.Run(); err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) icons.sh error: %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	os.Remove(fmt.Sprintf("%s/icons.sh", flagShortName))
	sendStatusUpdate("Info", "Generating icons completed... starting build additional files...")
	sendMsgToUI(flowProcessedMsg{})
	// ------------------ move local image files ------------------
	utils.TuiLogger("Info", "(buildBankingPwa) moving local images")
	// TODO add windows support for "/"
	fileNameIcon := lastString(strings.Split(flagIcon, "/"))
	pasteIcon := flagShortName + "/" + fileNameIcon
	cmdCpIcon := exec.Command("cp", flagIcon, pasteIcon)

	if err := cmdCpIcon.Run(); err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) pasteIcon: %s", pasteIcon))
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) copy icon error: %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	// TODO add windows support for "/"
	fileNameBanner := lastString(strings.Split(flagBanner, "/"))
	pasteBanner := flagShortName + "/" + fileNameBanner
	cmdCpBanner := exec.Command("cp", flagBanner, pasteBanner)
	if err := cmdCpBanner.Run(); err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) pasteBanner: %s flagBanner: %s", pasteBanner, flagBanner))
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) copy banner error: %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	// ------------------ build home.html file ------------------
	sendStatusUpdate("Info", "Generating home.html file")
	home, err := pwaTemplates.ReadFile("_pwaTemplates/home.html")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedHome := strings.ReplaceAll(string(home), "LOGO", fileNameIcon)
	formattedHome = strings.ReplaceAll(string(formattedHome), "THEME_COLOR", flagColor)
	formattedHome = strings.ReplaceAll(string(formattedHome), "BANNER", fileNameBanner)
	formattedHome = strings.ReplaceAll(string(formattedHome), "GC_REGION", flagRegion)
	formattedHome = strings.ReplaceAll(string(formattedHome), "GC_ENVIRONMENT", flagEnvironment)
	formattedHome = strings.ReplaceAll(string(formattedHome), "GC_DEPLOYMENT_ID", flagDeploymentId)
	err = os.WriteFile(fmt.Sprintf("%s/home.html", flagShortName), []byte(formattedHome), 0777)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	// ------------------ build index.html file ------------------
	sendStatusUpdate("Info", "Generating index.html file")
	index, err := pwaTemplates.ReadFile("_pwaTemplates/index.html")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedIndex := strings.ReplaceAll(string(index), "LOGO", fileNameIcon)
	formattedIndex = strings.ReplaceAll(string(formattedIndex), "THEME_COLOR", flagColor)
	formattedIndex = strings.ReplaceAll(string(formattedIndex), "BANNER", fileNameBanner)
	err = os.WriteFile(fmt.Sprintf("%s/index.html", flagShortName), []byte(formattedIndex), 0777)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build index.css file ------------------
	sendStatusUpdate("Info", "Generating index.css file")
	css, err := pwaTemplates.ReadFile("_pwaTemplates/index.css")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedCss := strings.ReplaceAll(string(css), "THEME_COLOR", flagColor)
	err = os.WriteFile(fmt.Sprintf("%s/index.css", flagShortName), []byte(formattedCss), 0777)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build manifest.json file ------------------
	sendStatusUpdate("Info", "Generating manifest.json file")
	manifest, err := pwaTemplates.ReadFile("_pwaTemplates/manifest.json")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedManifest := strings.ReplaceAll(string(manifest), "DEMO_NAME", flagName)
	formattedManifest = strings.ReplaceAll(string(formattedManifest), "THEME_COLOR", flagColor)
	formattedManifest = strings.ReplaceAll(string(formattedManifest), "DEMO_SHORT_NAME", flagShortName)
	err = os.WriteFile(fmt.Sprintf("%s/manifest.json", flagShortName), []byte(formattedManifest), 0777)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build deploy.sh file ------------------
	sendStatusUpdate("Info", "Generating deploy.sh file")
	deploy, err := pwaTemplates.ReadFile("_pwaTemplates/deploy.sh")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedDeploy := strings.ReplaceAll(string(deploy), "$bucketName", flagBucketName)
	formattedDeploy = strings.ReplaceAll(string(formattedDeploy), "$shortName", flagShortName)
	err = os.WriteFile(fmt.Sprintf("%s/deploy.sh", flagShortName), []byte(formattedDeploy), 0777)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build script.js file ------------------
	script, err := pwaTemplates.ReadFile("_pwaTemplates/script.js")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedScript := strings.ReplaceAll(string(script), "LOGO", fileNameIcon)
	err = os.WriteFile(fmt.Sprintf("%s/script.js", flagShortName), []byte(formattedScript), 0777)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	sendStatusUpdate("Info", "Generating script.js file")
	// ------------------ build genesys.js file ------------------
	genesys, err := pwaTemplates.ReadFile("_pwaTemplates/genesys.js")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedGenesys := strings.ReplaceAll(string(genesys), "LOGO", fileNameIcon)
	err = os.WriteFile(fmt.Sprintf("%s/script.js", flagShortName), []byte(formattedGenesys), 0777)
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		_ = os.RemoveAll(flagShortName)
		return
	}
	sendStatusUpdate("Info", "Generating genesys.js file")
	// ------------------ build service-worker.js file ------------------
	utils.TuiLogger("Info", "(buildBankingPwa) Generating service-worker.js")
	err = createFile("service-worker.js", flagShortName, "_pwaTemplates/service-worker.js")
	if err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("(buildBankingPwa) %s", err))
		return
	}
	sendStatusUpdate("Info", "Build COMPLETED")
	sendMsgToUI(flowProcessedMsg{})
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

func createFile(file, directory, embeddedLocation string) error {
	data, err := pwaTemplates.ReadFile(embeddedLocation)
	if err != nil {
		_ = os.RemoveAll(directory)
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", directory, file), []byte(data), 0777)
	if err != nil {
		_ = os.RemoveAll(directory)
		return err
	}
	return nil
}
