package genesysLogin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/mcphee11/mcphee11-tui/utils"
	"github.com/mypurecloud/platform-client-sdk-go/platformclientv2"
)

func GenesysLogin() (configReturned *platformclientv2.Configuration, err error) {
	region := os.Getenv("MCPHEE11_TUI_REGION")
	clientId := os.Getenv("MCPHEE11_TUI_CLIENT_ID")
	secret := os.Getenv("MCPHEE11_TUI_SECRET")

	if region != "" && clientId != "" && secret != "" {
		utils.TuiLogger("Info", "Using environment variables for Genesys Cloud configuration")
		//Do Genesys Cloud OAuth
		config := platformclientv2.GetDefaultConfiguration()
		config.BasePath = "https://api." + region
		if err := config.AuthorizeClientCredentials(clientId, secret); err != nil {
			return nil, fmt.Errorf("failed to authorize client credentials: %w", err)
		}
		return config, nil
	} else {
		// Check if GC CLI is installed
		gcPath, err := exec.LookPath("gc")
		if err != nil {
			utils.TuiLogger("Info", "Genesys Cloud CLI (gc) not installed")
		}

		// If GC installed use that to login
		if gcPath != "" {
			utils.TuiLogger("Info", "Genesys Cloud CLI (gc) found")
			profile := os.Getenv("MCPHEE11_TUI_PROFILE")
			if profile == "" {
				profile = "default"
				utils.TuiLogger("Info", "Using default profile: "+profile)
			} else {
				utils.TuiLogger("Info", "Using profile: "+profile)
			}
			configFile := os.ExpandEnv("$HOME/.gc/config.toml")
			tomlData, err := os.ReadFile(configFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}

			type profileConfig struct {
				Environment string `toml:"environment"`
				ClientId    string `toml:"client_credentials"`
				Secret      string `toml:"client_secret"`
			}
			var configMap map[string]profileConfig
			if _, err := toml.Decode(string(tomlData), &configMap); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}

			profileSection, ok := configMap[profile]
			if !ok {
				return nil, fmt.Errorf("profile %s not found in config file", profile)
			}
			region = profileSection.Environment
			clientId = profileSection.ClientId
			secret = profileSection.Secret

			//Do Genesys Cloud OAuth
			config := platformclientv2.GetDefaultConfiguration()
			config.BasePath = "https://api." + region
			if err := config.AuthorizeClientCredentials(clientId, secret); err != nil {
				return nil, fmt.Errorf("failed to authorize client credentials: %w", err)
			}
			return config, nil
		}
	}
	return nil, fmt.Errorf("Genesys Cloud CLI (gc) not found and environment variables are not set cant login to Genesys Cloud")
}

func GenesysCreds() (region string, clientId string, secret string, err error) {
	region = os.Getenv("MCPHEE11_TUI_REGION")
	clientId = os.Getenv("MCPHEE11_TUI_CLIENT_ID")
	secret = os.Getenv("MCPHEE11_TUI_SECRET")

	if region != "" && clientId != "" && secret != "" {
		return region, clientId, secret, nil
	} else {
		// Check if GC CLI is installed
		gcPath, err := exec.LookPath("gc")
		if err != nil {
			utils.TuiLogger("Error", "No Genesys Cloud Creds")
		}

		// If GC installed use that to login
		if gcPath != "" {
			profile := os.Getenv("MCPHEE11_TUI_PROFILE")
			if profile == "" {
				profile = "default"
			}
			configFile := os.ExpandEnv("$HOME/.gc/config.toml")
			tomlData, err := os.ReadFile(configFile)
			if err != nil {
				return "", "", "", fmt.Errorf("failed to read config file: %w", err)
			}

			type profileConfig struct {
				Environment string `toml:"environment"`
				ClientId    string `toml:"client_credentials"`
				Secret      string `toml:"client_secret"`
			}
			var configMap map[string]profileConfig
			if _, err := toml.Decode(string(tomlData), &configMap); err != nil {
				return "", "", "", fmt.Errorf("failed to parse config file: %w", err)
			}

			profileSection, ok := configMap[profile]
			if !ok {
				return "", "", "", fmt.Errorf("profile %s not found in config file", profile)
			}
			region = profileSection.Environment
			clientId = profileSection.ClientId
			secret = profileSection.Secret

			return region, clientId, secret, nil
		}
	}
	return "", "", "", fmt.Errorf("Genesys Cloud CLI (gc) not found and environment variables are not set cant login to Genesys Cloud")
}
