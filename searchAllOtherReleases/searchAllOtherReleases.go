package searchAllOtherReleases

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/mcphee11/mcphee11-tui/utils"
)

type pageInfo struct {
	Type     string
	DateType string
	Link     string
}

var pages = []pageInfo{
	// embedded clients https://help.mypurecloud.com/articles/release-notes-for-the-genesys-cloud-embedded-clients/
	{Type: "cxCloud", DateType: "h2", Link: "https://help.mypurecloud.com/articles/release-notes-for-cx-cloud-from-genesys-and-salesforce/"},
	{Type: "teams", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-for-genesys-cloud-for-microsoft-teams/"},
	{Type: "zendesk", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-for-genesys-cloud-for-zendesk/"},
	{Type: "salesforce", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-genesys-cloud-salesforce/"},
	{Type: "browserExtension", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-for-the-genesys-cloud-browser-extensions/"},
	// data actions https://help.mypurecloud.com/articles/release-notes-for-the-data-actions-integrations/
	{Type: "awsLambda", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-for-the-aws-lambda-data-actions-integration/"},
	{Type: "genesysFunction", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-for-the-genesys-cloud-function-data-actions-integration/"},
	{Type: "genesysDataAction", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-for-the-genesys-cloud-data-actions-integration/"},
	{Type: "googleDataAction", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-for-the-google-data-actions-integration/"},
	{Type: "dynamicsDataAction", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-microsoft-dynamics-365-data-actions-integration/"},
	{Type: "salesforceDataAction", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-salesforce-data-actions-integration/"},
	{Type: "webServicesDataAction", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-web-services-data-actions-integration/"},
	{Type: "zendeskDataAction", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-zendesk-data-actions-integration/"},
	// SCIM https://help.mypurecloud.com/articles/release-notes-for-genesys-cloud-scim-identity-management/
	{Type: "scim", DateType: "h3", Link: "https://help.mypurecloud.com/articles/release-notes-for-genesys-cloud-scim-identity-management/"},
	// desktop apps
	{Type: "desktopAppMAC", DateType: "span", Link: "https://help.mypurecloud.com/release-notes-home/genesys-cloud-for-mac-desktop-app-release-notes/"},
	{Type: "desktopAppWindows", DateType: "span", Link: "https://help.mypurecloud.com/release-notes-home/genesys-cloud-for-windows-desktop-app-release-notes/"},
	{Type: "desktopGCBA", DateType: "span", Link: "https://help.mypurecloud.com/articles/genesys-cloud-background-assistant-gcba-release-notes/"},
}

func SearchAllOtherReleases(searchString string) []map[string]string {
	var totalReleases []map[string]string

	for _, page := range pages {
		var embeddedReleases []map[string]string
		if strings.Contains(page.Type, "desktop") {
			embeddedReleases = scrapType2(searchString, page.Type, page.DateType, page.Link)
		} else {
			embeddedReleases = scrap(searchString, page.Type, page.DateType, page.Link)
		}
		totalReleases = append(totalReleases, embeddedReleases...)
		utils.TuiLogger("Info", fmt.Sprintf("Found %d releases in embedded %s", len(embeddedReleases), page.Type))
	}

	return totalReleases
}

func scrap(searchString, sectionType, dateType, link string) []map[string]string {
	var releases []map[string]string
	errCount := 0
	matchedCount := 0

	c := colly.NewCollector(colly.Async(false)) //removed async to stop 429

	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 1 request at a time to stop 429
		Parallelism: 1,
	})

	c.OnError(func(_ *colly.Response, err error) {
		utils.TuiLogger("Error", fmt.Sprintf("Something went wrong: %s", err))
		errCount++
	})

	c.OnHTML("ul", func(e *colly.HTMLElement) {

		if strings.Contains(strings.ToLower(e.Text), strings.ToLower(searchString)) {
			var rawSection string
			var rawNote = strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix(strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", " | "), "	", ""), " | "), " | "), " |  | ", " | ")
			if e.DOM.Prev().Text() != e.DOM.PrevAllFiltered(dateType).First().Text() && !strings.HasPrefix(e.DOM.Prev().Text(), "Build ") {
				rawSection = e.DOM.Prev().Text() + " | " + e.DOM.PrevAllFiltered(dateType).First().Text()
			} else {
				rawSection = sectionType + " | " + e.DOM.PrevAllFiltered(dateType).First().Text()
			}
			var release = map[string]string{
				"section": rawSection,
				"link":    link,
				"notes":   strings.ReplaceAll(strings.TrimPrefix(rawNote, " | "), " |  | ", " | "), //strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(e.Text, "\n", " | "), " | "), " | "),
			}
			// only append if section is not empty
			if e.DOM.PrevAllFiltered(dateType).First().Text() != "" {
				releases = append(releases, release)
				matchedCount++
			}
		}
	})

	c.Visit(link)
	c.Wait()

	utils.TuiLogger("Info", fmt.Sprintf("Found String: %s on %d pages. With %d errors", searchString, matchedCount, errCount))
	return releases
}

func scrapType2(searchString, sectionType, dateType, link string) []map[string]string {
	var releases []map[string]string
	errCount := 0
	matchedCount := 0

	c := colly.NewCollector(colly.Async(false)) //removed async to stop 429

	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 1 request at a time to stop 429
		Parallelism: 1,
	})

	c.OnError(func(_ *colly.Response, err error) {
		utils.TuiLogger("Error", fmt.Sprintf("Something went wrong: %s", err))
		errCount++
	})

	c.OnHTML("ul", func(e *colly.HTMLElement) {

		if strings.Contains(strings.ToLower(e.Text), strings.ToLower(searchString)) {
			var rawSection string
			var rawNote = strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix(strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", " | "), "	", ""), " | "), " | "), " |  | ", " | ")
			rawSection = sectionType + " | " + fmt.Sprintf("%s", strings.ReplaceAll(strings.ReplaceAll(e.DOM.Parent().Prev().Text(), "  ", ""), "\n", ""))

			var release = map[string]string{
				"section": rawSection,
				"link":    link,
				"notes":   strings.ReplaceAll(strings.TrimPrefix(rawNote, " | "), " |  | ", " | "),
			}
			// only append if section is not correct ul
			if e.DOM.Parent().HasClass("accordion__content") {
				releases = append(releases, release)
				matchedCount++
			}
		}
	})

	c.Visit(link)
	c.Wait()

	utils.TuiLogger("Info", fmt.Sprintf("Found String: %s on %d pages. With %d errors", searchString, matchedCount, errCount))
	return releases
}
