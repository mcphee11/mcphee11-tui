package pwaDeploy

import (
	"archive/zip"
	"bytes"
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

const API_BASE_URL = "https://www.pwabuilder.com/api/images/generateStoreImages"

type APIResponse struct {
	Uri string `json:"Uri"`
}

func GenerateIcons(iconPath, flagShortName string) {
	utils.TuiLogger("Info", fmt.Sprintf("Creating app images from %s...", iconPath))

	// --- Step 1: Upload the icon file and get the response URL ---
	zipFilePath, err := uploadIcon(iconPath)
	if err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("Error uploading icon: %v", err))
	}
	utils.TuiLogger("Info", "Downloaded Images")

	// --- Step 2: unzip files ---
	appImagesDir := fmt.Sprintf("%s/AppImages", flagShortName)
	utils.TuiLogger("Info", fmt.Sprintf("Unzipping %s into %s directory...", zipFilePath, appImagesDir))

	if err := unzipFile(zipFilePath, appImagesDir); err != nil {
		utils.TuiLogger("Fatal", fmt.Sprintf("Error unzipping file: %v", err))
	}
	utils.TuiLogger("Info", fmt.Sprintf("Successfully unzipped app images into %s.", appImagesDir))

	// --- Step 3: Clean up by removing the downloaded zip file ---
	if err := os.Remove(zipFilePath); err != nil {
		utils.TuiLogger("Error", fmt.Sprintf("Error removing zip file: %v", err))
	}
	utils.TuiLogger("Info", "Cleanup complete. App images are ready!")
}

// uploadIcon uploads the specified icon file to the API and returns the URI from the response.
func uploadIcon(iconPath string) (string, error) { //string, string, error) {
	file, err := os.Open(iconPath)
	if err != nil {
		return "", fmt.Errorf("failed to open icon file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("baseImage", filepath.Base(iconPath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content to form: %w", err)
	}

	_ = writer.WriteField("padding", "0.3")
	_ = writer.WriteField("backgroundColor", "transparent")
	_ = writer.WriteField("platforms", "windows11")
	_ = writer.WriteField("platforms", "android")
	_ = writer.WriteField("platforms", "ios")

	writer.Close()

	req, err := http.NewRequest("POST", API_BASE_URL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(respBodyBytes))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	err = os.WriteFile("icons.zip", respBody, 0777)
	if err != nil {
		return "", fmt.Errorf("failed saving zip file: %w", err)
	}

	return "icons.zip", nil
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
