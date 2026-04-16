package searchReleaseNotes

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/mcphee11/mcphee11-tui/utils"
)

func SearchReleaseNotes(searchString string) []map[string]string {
	var releases []map[string]string
	errCount := 0
	matchedCount := 0

	c := colly.NewCollector(colly.Async(false)) // removed async to stop 429

	c.Limit(&colly.LimitRule{
		// limit the parallel requests to 1 request at a time to stop 429
		Parallelism: 1,
	})

	c.OnError(func(_ *colly.Response, err error) {
		utils.TuiLogger("Error", fmt.Sprintf("(searchReleaseNotes) Something went wrong: %s", err))
		errCount++
	})

	c.OnHTML("ul", func(e *colly.HTMLElement) {
		if strings.Contains(strings.ToLower(e.Text), strings.ToLower(searchString)) {
			release := map[string]string{
				"section": e.DOM.Prev().Text() + " | " + strings.TrimPrefix(strings.TrimSuffix(e.DOM.Parent().NextAllFiltered("p").First().Children().AttrOr("href", ""), "/"), "/release-notes/genesys-cloud/"),
				"link":    "https://help.genesys.cloud" + e.DOM.Parent().NextAllFiltered("p").First().Children().AttrOr("href", ""),
				"notes":   strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(e.Text, "\n", " | "), " | "), " | "),
			}
			releases = append(releases, release)
			matchedCount++
		}
	})

	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2026")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2025")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2024")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2023")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2022")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2021")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2020")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2019")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2018")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2017")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2016")
	c.Visit("https://help.genesys.cloud/release-notes/genesys-cloud/archive/2015")

	c.Wait()

	utils.TuiLogger("Info", fmt.Sprintf("(searchReleaseNotes) Found String: %s on %d pages. With %d errors", searchString, matchedCount, errCount))
	return releases
}
