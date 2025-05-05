package searchReleaseNotes

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
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
		fmt.Println("Something went wrong: ", err)
		errCount++
	})

	c.OnHTML("ul", func(e *colly.HTMLElement) {

		if strings.Contains(strings.ToLower(e.Text), strings.ToLower(searchString)) {
			var release = map[string]string{
				"section": e.DOM.Prev().Text(),
				"link":    e.DOM.NextAllFiltered("p.view").First().Children().AttrOr("href", ""),
				"notes":   strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(e.Text, "\n", " | "), " | "), " | "),
			}
			releases = append(releases, release)
			matchedCount++
		}
	})

	c.Visit("https://help.mypurecloud.com/monthly-archive/")

	c.Wait()

	fmt.Printf("Found String: %s on %d pages. With %d errors\n", searchString, matchedCount, errCount)
	return releases
}
