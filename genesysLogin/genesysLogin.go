package genesysLogin

import (
	"fmt"
	"os"

	"github.com/mypurecloud/platform-client-sdk-go/platformclientv2"
)

func GenesysLogin() (configReturned *platformclientv2.Configuration, err error) {

	region := os.Getenv("MCPHEE11_TUI_REGION")
	if region == "" {
		fmt.Println("Environment variable MCPHEE11_TUI_REGION is not set")
		return
	}

	clientId := os.Getenv("MCPHEE11_TUI_CLIENT_ID")
	if region == "" {
		fmt.Println("Environment variable MCPHEE11_TUI_CLIENT_ID is not set")
		return
	}

	secret := os.Getenv("MCPHEE11_TUI_SECRET")
	if region == "" {
		fmt.Println("Environment variable MCPHEE11_TUI_SECRET is not set")
		return
	}

	//Do Genesys Cloud OAuth
	config := platformclientv2.GetDefaultConfiguration()
	config.BasePath = "https://api." + region
	if err := config.AuthorizeClientCredentials(clientId, secret); err != nil {
		return nil, fmt.Errorf("failed to authorize client credentials: %w", err)
	}
	return config, nil
}
