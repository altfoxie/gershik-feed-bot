package storage

type Feed struct {
	ID        int64  `json:"id"`
	ChannelID int64  `json:"channel_id"`
	Name      string `json:"name"`
}

func Feeds() []Feed {
	return gStorage.Feeds
}

func FeedByID(id int64) *Feed {
	for _, feed := range Feeds() {
		if feed.ID == id {
			return &feed
		}
	}
	return nil
}

func FeedByChannelID(channelID int64) *Feed {
	for _, feed := range Feeds() {
		if feed.ChannelID == channelID {
			return &feed
		}
	}
	return nil
}

func CreateFeed(channelID int64, name string) {
	gStorage.Lock()
	defer gStorage.Unlock()

	feed := Feed{
		ID:        gStorage.FeedsID,
		ChannelID: channelID,
		Name:      name,
	}
	gStorage.FeedsID++
	gStorage.Feeds = append(gStorage.Feeds, feed)
	save()
}

func DeleteFeed(id int64) bool {
	gStorage.Lock()
	defer gStorage.Unlock()

	for i, feed := range gStorage.Feeds {
		if feed.ID == id {
			gStorage.Feeds = append(gStorage.Feeds[:i], gStorage.Feeds[i+1:]...)
			save()
			return true
		}
	}
	return false
}
