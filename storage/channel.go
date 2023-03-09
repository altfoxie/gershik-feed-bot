package storage

type Channel struct {
	ID            int64  `json:"id"`
	ChannelID     int64  `json:"channel_id"`
	Name          string `json:"name"`
	FeedID        int64  `json:"feed_id"`
	LastMessageID int    `json:"last_message_id"`
}

func Channels() []Channel {
	return gStorage.Channels
}

func ChannelByID(id int64) *Channel {
	for _, channel := range Channels() {
		if channel.ID == id {
			return &channel
		}
	}
	return nil
}

func ChannelByChannelID(channelID int64) *Channel {
	for _, channel := range Channels() {
		if channel.ChannelID == channelID {
			return &channel
		}
	}
	return nil
}

func CreateChannel(channelID, feedID int64, name string) {
	gStorage.Lock()
	defer gStorage.Unlock()

	channel := Channel{
		ID:        gStorage.ChannelsID,
		ChannelID: channelID,
		Name:      name,
		FeedID:    feedID,
	}
	gStorage.ChannelsID++
	gStorage.Channels = append(gStorage.Channels, channel)
	save()
}

func DeleteChannel(id int64) bool {
	gStorage.Lock()
	defer gStorage.Unlock()

	for i, channel := range gStorage.Channels {
		if channel.ID == id {
			gStorage.Channels = append(gStorage.Channels[:i], gStorage.Channels[i+1:]...)
			save()
			return true
		}
	}
	return false
}

func UpdateChannelLastMessageID(id int64, lastMessageID int) bool {
	gStorage.Lock()
	defer gStorage.Unlock()

	for i, channel := range gStorage.Channels {
		if channel.ID == id {
			gStorage.Channels[i].LastMessageID = lastMessageID
			save()
			return true
		}
	}
	return false
}
