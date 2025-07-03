package pwaDeploy

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mcphee11/mcphee11-tui/utils"
)

const API_BASE_URL = "https://appimagegenerator-prod-dev.azurewebsites.net"

type APIResponse struct {
	Uri string `json:"Uri"`
}

func GenerateIcons(iconPath, flagShortName string) {
	utils.TuiLogger("Info", fmt.Sprintf("Creating app images from %s...", iconPath))

	// --- Step 1: Upload the icon file and get the response URL ---
	responseURL, downloadedFileName, err := uploadIcon(iconPath)
	if err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("Error uploading icon: %v", err))
	}
	utils.TuiLogger("Info", fmt.Sprintf("Downloading from: %s%s...", API_BASE_URL, responseURL))

	// --- Step 2: Download the generated zip file ---
	zipFilePath, err := downloadFile(responseURL, downloadedFileName)
	if err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("Error downloading file: %v", err))
	}
	utils.TuiLogger("Info", fmt.Sprintf("Downloaded zip file to: %s", zipFilePath))

	// --- Step 3: Unzip the downloaded file into the AppImages directory ---
	appImagesDir := fmt.Sprintf("%s/AppImages", flagShortName)
	utils.TuiLogger("Info", fmt.Sprintf("Unzipping %s into %s directory...", zipFilePath, appImagesDir))

	if err := unzipFile(zipFilePath, appImagesDir); err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("Error unzipping file: %v", err))
	}
	utils.TuiLogger("Info", fmt.Sprintf("Successfully unzipped app images into %s.", appImagesDir))

	// --- Step 4: Clean up by removing the downloaded zip file ---
	if err := os.Remove(zipFilePath); err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("Error removing zip file: %v", err))
	}
	utils.TuiLogger("Info", "Cleanup complete. App images are ready!")
}

// uploadIcon uploads the specified icon file to the API and returns the URI from the response.
func uploadIcon(iconPath string) (string, string, error) {
	file, err := os.Open(iconPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open icon file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("fileName", filepath.Base(iconPath))
	if err != nil {
		return "", "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", "", fmt.Errorf("failed to copy file content to form: %w", err)
	}

	_ = writer.WriteField("padding", "0.3")
	_ = writer.WriteField("color", "transparent")
	_ = writer.WriteField("platform", "windows11")
	_ = writer.WriteField("platform", "android")
	_ = writer.WriteField("platform", "ios")

	writer.Close()

	req, err := http.NewRequest("POST", API_BASE_URL+"/api/image", body)
	if err != nil {
		return "", "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBodyBytes, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(respBodyBytes))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse APIResponse
	err = json.Unmarshal(respBody, &apiResponse)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	downloadedFileName := apiResponse.Uri
	if strings.HasPrefix(apiResponse.Uri, "/api/") {
		downloadedFileName = apiResponse.Uri[5:]
	} else {
		downloadedFileName = apiResponse.Uri
	}
	// needed to remove "?" from file name to support windows I also did "=" for better formatting
	downloadedFileName = strings.ReplaceAll(downloadedFileName, "?", "_")
	downloadedFileName = strings.ReplaceAll(downloadedFileName, "=", "_")
	utils.TuiLogger("Info", fmt.Sprintf("Using '%s' as the file name.", downloadedFileName))

	return apiResponse.Uri, downloadedFileName, nil
}

// downloadFile downloads a file from the given URL and saves it to the current directory.
func downloadFile(uri string, fileName string) (string, error) {
	fullURL := API_BASE_URL + uri

	out, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", fileName, err)
	}
	defer out.Close()

	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to download file from %s: %w", fullURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write downloaded content to file: %w", err)
	}

	return fileName, nil
}

// unzipFile unzips a zip archive into a specified destination directory.
func unzipFile(zipFilePath, destDir string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fpath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip %s: %w", f.Name, err)
		}
		defer rc.Close()

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", fpath, err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return fmt.Errorf("failed to copy content for %s: %w", f.Name, err)
		}
	}
	return nil
}
