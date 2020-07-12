package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/vvksh/amigo"
	"github.com/vvksh/configurt"
	"github.com/vvksh/nightswatch"
)

var clConfigurtClient *configurt.Client
var fp = gofeed.NewParser()

var checkedPositings = make(map[string]bool)

// CraigslistWatcher implements nightswatch.Watcher interface
type CraigslistWatcher func()

func init() {
	githubAccessToken, exists := os.LookupEnv("GITHUB_TOKEN")
	if !exists {
		log.Panicf("GITHUB_TOKEN not found\n")
	}
	clConfigurtClient = configurt.NewClient("vvksh", githubAccessToken, "configs", "craigslist_watcher.json", 30*time.Minute)
	var cw CraigslistWatcher
	nightswatch.Register(&cw)
}

func (cw *CraigslistWatcher) Check() []string {
	searchQueries := clConfigurtClient.GetAsStringArray("search_queries")
	updates := []string{}

	for _, searchQuery := range searchQueries {
		feed, err := fp.ParseURL(searchQuery)
		if err != nil {
			log.Printf("error reading feed: %s error: %s\n", searchQuery, err.Error())
			continue
		}
		log.Printf("Got %d feed items", len(feed.Items))
		for _, item := range feed.Items {
			if _, ok := checkedPositings[item.Title]; !ok {
				// log.Printf("%v", item)
				update := getFormattedPosting(item)
				updates = append(updates, update)

				checkedPositings[item.Title] = true
			}
		}

	}

	if len(checkedPositings) > 1000 {
		checkedPositings = make(map[string]bool)
	}

	return updates

}

func (cw *CraigslistWatcher) Interval() time.Duration {
	intervalMin := clConfigurtClient.GetAsInt("interval_min")
	return time.Minute * time.Duration(intervalMin)
}

func (cw *CraigslistWatcher) SlackChannel() string {
	return clConfigurtClient.GetAsString("slack_channel")
}

func getFormattedPosting(posting *gofeed.Item) string {
	return fmt.Sprintf("%s \n <%s|%s> \n %s", posting.Published, posting.Link, amigo.Sanitize(posting.Title), amigo.Sanitize(posting.Description))
}

func main() {
	nightswatch.Start()
}
