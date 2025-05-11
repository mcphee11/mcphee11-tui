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

	c := colly.NewCollector(colly.Async(false)) //removed async to stop 429

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
			var release = map[string]string{
				"section": e.DOM.Prev().Text() + " | " + strings.TrimPrefix(strings.TrimSuffix(e.DOM.NextAllFiltered("p.view").First().Children().AttrOr("href", ""), "/"), "https://help.mypurecloud.com/releasenote/"),
				"link":    e.DOM.NextAllFiltered("p.view").First().Children().AttrOr("href", ""),
				"notes":   strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(e.Text, "\n", " | "), " | "), " | "),
			}
			releases = append(releases, release)
			matchedCount++
		}
	})

	c.Visit("https://help.mypurecloud.com/monthly-archive/")

	c.Wait()

	utils.TuiLogger("Info", fmt.Sprintf("(searchReleaseNotes) Found String: %s on %d pages. With %d errors", searchString, matchedCount, errCount))
	return releases
}
