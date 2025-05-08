package pwaDeploy

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
	sendStatusUpdate := func(s string) {
		sendMsgToUI(internalUpdateStatusMsg{newStatus: s})
	}
	// ------------ Create project folder -----------------
	err := os.Mkdir(flagShortName, 0777)
	if err != nil {
		fmt.Printf("Error creating directory %s, exiting build.", flagShortName)
		sendStatusUpdate(fmt.Sprintf("Error creating directory %s, exiting build.", flagShortName))
		sendMsgToUI(stage1CompleteMsg{})
		return
	}

	err = os.Mkdir(fmt.Sprintf("%s/svgs", flagShortName), 0777)
	if err != nil {
		fmt.Printf("Error creating directory %s/svgs, exiting build.", flagShortName)
		sendStatusUpdate(fmt.Sprintf("Error creating directory %s/svgs, exiting build.", flagShortName))
		sendMsgToUI(stage1CompleteMsg{})
		return
	}

	// ------------------ create svgs ------------------
	svgs, err := pwaTemplates.ReadDir("_pwaTemplates/svgs")
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	for i := 0; i < len(svgs); i++ {
		err := createFile(svgs[i].Name(), fmt.Sprintf("%s/svgs", flagShortName), fmt.Sprintf("_pwaTemplates/svgs/%s", svgs[i].Name()))
		if err != nil {
			fmt.Println(err.Error())
			_ = os.RemoveAll(flagShortName)
			return
		}
		sendStatusUpdate(fmt.Sprintf("Generated Svg: %s\n", svgs[i].Name()))
	}
	sendMsgToUI(flowProcessedMsg{})

	// ------------------ create icons ------------------
	sendStatusUpdate("Generating App Icons this can take a min so please wait...")
	icons, err := pwaTemplates.ReadFile("_pwaTemplates/icons.sh")
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedIcons := strings.ReplaceAll(string(icons), "$icon", flagIcon)
	err = os.WriteFile(fmt.Sprintf("%s/icons.sh", flagShortName), []byte(formattedIcons), 0777)
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting working dir: %s", err)
		_ = os.RemoveAll(flagShortName)
		return
	}
	cmdIcon := exec.Command("./icons.sh")
	cmdIcon.Dir = fmt.Sprintf("%s/%s", currentDir, flagShortName)

	if err := cmdIcon.Run(); err != nil {
		fmt.Printf("icons.sh error: %s", err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	os.Remove(fmt.Sprintf("%s/icons.sh", flagShortName))
	sendStatusUpdate("Generating icons completed... starting build additional files...")
	sendMsgToUI(flowProcessedMsg{})
	// ------------------ move local image files ------------------
	// TODO add windows support for "/"
	fileNameIcon := lastString(strings.Split(flagIcon, "/"))
	pasteIcon := flagShortName + "/" + fileNameIcon
	cmdCpIcon := exec.Command("cp", flagIcon, pasteIcon)

	if err := cmdCpIcon.Run(); err != nil {
		fmt.Printf("pasteIcon: %s", pasteIcon)
		fmt.Printf("copy icon error: %s", err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	// TODO add windows support for "/"
	fileNameBanner := lastString(strings.Split(flagBanner, "/"))
	pasteBanner := flagShortName + "/" + fileNameBanner
	cmdCpBanner := exec.Command("cp", flagBanner, pasteBanner)
	if err := cmdCpBanner.Run(); err != nil {
		fmt.Printf("pasteBanner: %s flagBanner: %s", pasteBanner, flagBanner)
		fmt.Printf("copy banner error: %s", err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build home.html file ------------------
	sendStatusUpdate("Generating home.html file")
	home, err := pwaTemplates.ReadFile("_pwaTemplates/home.html")
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build index.html file ------------------
	sendStatusUpdate("Generating index.html file")
	index, err := pwaTemplates.ReadFile("_pwaTemplates/index.html")
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedIndex := strings.ReplaceAll(string(index), "LOGO", fileNameIcon)
	formattedIndex = strings.ReplaceAll(string(formattedIndex), "THEME_COLOR", flagColor)
	formattedIndex = strings.ReplaceAll(string(formattedIndex), "BANNER", fileNameBanner)
	err = os.WriteFile(fmt.Sprintf("%s/index.html", flagShortName), []byte(formattedIndex), 0777)
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build index.css file ------------------
	sendStatusUpdate("Generating index.css file")
	css, err := pwaTemplates.ReadFile("_pwaTemplates/index.css")
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedCss := strings.ReplaceAll(string(css), "THEME_COLOR", flagColor)
	err = os.WriteFile(fmt.Sprintf("%s/index.css", flagShortName), []byte(formattedCss), 0777)
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build manifest.json file ------------------
	sendStatusUpdate("Generating manifest.json file")
	manifest, err := pwaTemplates.ReadFile("_pwaTemplates/manifest.json")
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedManifest := strings.ReplaceAll(string(manifest), "DEMO_NAME", flagName)
	formattedManifest = strings.ReplaceAll(string(formattedManifest), "THEME_COLOR", flagColor)
	formattedManifest = strings.ReplaceAll(string(formattedManifest), "DEMO_SHORT_NAME", flagShortName)
	err = os.WriteFile(fmt.Sprintf("%s/manifest.json", flagShortName), []byte(formattedManifest), 0777)
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build deploy.sh file ------------------
	sendStatusUpdate("Generating deploy.sh file")
	deploy, err := pwaTemplates.ReadFile("_pwaTemplates/deploy.sh")
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}
	formattedDeploy := strings.ReplaceAll(string(deploy), "$bucketName", flagBucketName)
	formattedDeploy = strings.ReplaceAll(string(formattedDeploy), "$shortName", flagShortName)
	err = os.WriteFile(fmt.Sprintf("%s/deploy.sh", flagShortName), []byte(formattedDeploy), 0777)
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(flagShortName)
		return
	}

	// ------------------ build script.js file ------------------
	err = createFile("script.js", flagShortName, "_pwaTemplates/script.js")
	if err != nil {
		return
	}
	sendStatusUpdate("Generating script.js file")
	// ------------------ build genesys.js file ------------------
	err = createFile("genesys.js", flagShortName, "_pwaTemplates/genesys.js")
	if err != nil {
		return
	}
	sendStatusUpdate("Generating genesys.js file")
	// ------------------ build service-worker.js file ------------------
	err = createFile("service-worker.js", flagShortName, "_pwaTemplates/service-worker.js")
	if err != nil {
		return
	}
	sendStatusUpdate("Build COMPLETED")
	sendMsgToUI(flowProcessedMsg{})

	//fmt.Println("PS. Don't forget there is a deploy.sh file for you to deploy it to GCP")
}

func lastString(ss []string) string {
	return ss[len(ss)-1]
}

func createFile(file, directory, embeddedLocation string) error {
	data, err := pwaTemplates.ReadFile(embeddedLocation)
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(directory)
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", directory, file), []byte(data), 0777)
	if err != nil {
		fmt.Println(err.Error())
		_ = os.RemoveAll(directory)
		return err
	}
	return nil
}
