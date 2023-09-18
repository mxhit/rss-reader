package main

import (
	"rss-reader/database"
	"rss-reader/reader"
	"rss-reader/utils"
)

func main() {
	latestItems := reader.GetFeedUpdates()
	feedEntities := database.GetExistingData(latestItems)

	toUpdate := compare(latestItems, feedEntities)

	if len(toUpdate) > 0 {
		database.UpdateTable(toUpdate)
		utils.SendUpdateMail(toUpdate)
	}
}

func compare(latestItems, lastEntries map[string]string) map[string]string {
	toUpdate := make(map[string]string, 0)

	for author, title := range latestItems {
		if lastEntry, ok := lastEntries[author]; ok {
			if lastEntry != title {
				toUpdate[author] = title
			} else {
				continue
			}
		}
	}

	return toUpdate
}
