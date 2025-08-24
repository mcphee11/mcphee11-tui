// Package webtemplates: used to build neovim templates for coding examples
package webtemplates

import (
	"embed"
	"fmt"
	"os"

	"github.com/mcphee11/mcphee11-tui/utils"
)

//go:embed _webTemplateFiles/*/*
var webTemplateFiles embed.FS

//go:embed _webTemplateFiles/build.sh
var webTemplateBuild embed.FS

func WebTemplatesList() []map[string]string {
	utils.TuiLogger("Info", "building web template items")
	list := []map[string]string{
		{"id": "webOne", "title": "Basic HTML page", "desc": "HTML page with basic JS OAuth"},
		{"id": "webTwo", "title": "HTML page with WSS", "desc": "HTML page with JS OAuth & WSS"},
		{"id": "webThree", "title": "Custom Report HTML page", "desc": "HTML page with JS OAuth & Report template"},
	}
	return list
}

func BuildWebTemplate(template string) []map[string]string {
	utils.TuiLogger("Info", fmt.Sprintf("Building web template: %s", template))
	projectName := "project"

	// ------------ Create project folder -----------------
	err := os.Mkdir(projectName, 0o777)
	if err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("Error creating directory %s, exiting build.", template))
	}

	// ------------ Copy template folder -----------------
	files, err := webTemplateFiles.ReadDir(fmt.Sprintf("_webTemplateFiles/%s", template))
	if err != nil {
		_ = os.RemoveAll(projectName)
		utils.TuiLogger("Fatal", fmt.Sprintf("(build web template) Error reading dir: %s", err))
	}

	for _, file := range files {
		err := createFile(file.Name(), projectName, fmt.Sprintf("_webTemplateFiles/%s/%s", template, file.Name()))
		if err != nil {
			_ = os.RemoveAll(projectName)
			utils.TuiLogger("Fatal", fmt.Sprintf("(build web template) Error reading dir: %s", err))
		}
		utils.TuiLogger("Info", fmt.Sprintf("Generated File: %s", file.Name()))
	}

	// ------------ Copy build script -----------------
	build, err := webTemplateBuild.ReadFile("_webTemplateFiles/build.sh")
	if err != nil {
		_ = os.RemoveAll(projectName)
		utils.TuiLogger("Fatal", fmt.Sprintf("(build web template) reading file: %s", err))
	}
	err = os.WriteFile(fmt.Sprintf("%s/build.sh", projectName), []byte(build), 0o777)
	if err != nil {
		_ = os.RemoveAll(projectName)
		utils.TuiLogger("Fatal", fmt.Sprintf("(build web template) writting file: %s", err))
	}
	utils.TuiLogger("Info", "Generated File: build.sh")

	list := []map[string]string{}
	return list
}

func createFile(file, directory, embeddedLocation string) error {
	data, err := webTemplateFiles.ReadFile(embeddedLocation)
	if err != nil {
		_ = os.RemoveAll(directory)
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", directory, file), []byte(data), 0o777)
	if err != nil {
		_ = os.RemoveAll(directory)
		return err
	}
	return nil
}
