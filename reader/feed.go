package reader

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
)

const FEED_FILE string = "./reader/feed.xml"

type Feed struct {
	XmlName xml.Name `xml:"feed"`
	Blogs   []Blog   `xml:"blog"`
}

type Blog struct {
	XmlName xml.Name `xml:"blog"`
	Author  string   `xml:"author"`
	Url     string   `xml:"url"`
}

type Rss struct {
	XmlName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XmlName       xml.Name `xml:"channel"`
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Image         Image    `xml:"image"`
	Generator     string   `xml:"generator"`
	Language      string   `xml:"language"`
	LastBuildDate string   `xml:"lastBuildDate"`
	AtomLink      string   `xml:"atom:link"`
	Items         []Item   `xml:"item"`
}

type Image struct {
	XmlName xml.Name `xml:"image"`
	Url     string   `xml:"url"`
	Link    string   `xml:"link"`
}

type Item struct {
	XmlName     xml.Name `xml:"image"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Guid        string   `xml:"guid"`
	Description string   `xml:"description"`
}

func GetFeedUpdates() map[string]string {
	// Get the URL for the RSS XML file
	feed := getFeedLinks()

	// Fetch data from the blogUrls
	latestAuthorItem := getUpdates(feed)

	return latestAuthorItem
}

func getFeedLinks() Feed {
	feedLinks, err := os.Open(FEED_FILE)
	if err != nil {
		log.Panicf("Something went wrong while opening %s", FEED_FILE)
	}

	linksInBytes, _ := io.ReadAll(feedLinks)

	var feed Feed
	xml.Unmarshal(linksInBytes, &feed)

	return feed
}

func printFeed(blogs []Blog) {
	for _, blog := range blogs {
		log.Printf("Author: %s\nURL: %s\n\n", blog.Author, blog.Url)
	}
}

func getUpdates(feed Feed) map[string]string {
	latestAuthorItem := make(map[string]string, 0)

	for i := 0; i < 2; i++ {
		url := feed.Blogs[i].Url
		author := feed.Blogs[i].Author

		log.Printf("Fetching URL: %s\n", url)
		rssXml, err := http.Get(url)

		if err != nil {
			log.Panicf("Something went wrong while fetching %s", url)
		}

		defer rssXml.Body.Close()

		rssFileBytes, err := io.ReadAll(rssXml.Body)
		if err != nil {
			log.Panicln("Something went wrong while reading XML body")
		}

		var rss Rss
		xml.Unmarshal(rssFileBytes, &rss)

		latestAuthorItem[author] = rss.Channel.Items[0].Title
	}

	return latestAuthorItem
}

func printRssItems(items []Item) {
	for _, item := range items {
		log.Printf("Title: %s\nPublished Date: %s\nGUID: %s\nLink: %s\n", item.Title, item.PubDate, item.Guid, item.Link)
	}
}

func printMap(latestAuthorItem map[string]Item) {
	for k, v := range latestAuthorItem {
		log.Printf("\nAuthor: %s\nTitle: %s\n\n", k, v.Title)
	}
}
