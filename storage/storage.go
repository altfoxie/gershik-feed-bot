// Package storage provides a persistent storage for feed data.
package storage

import (
	"encoding/json"
	"os"
	"sync"
)

type storage struct {
	sync.Mutex `json:"-"`
	path       string `json:"-"`

	FeedsID int64  `json:"feed_id"`
	Feeds   []Feed `json:"feeds"`

	ChannelsID int64     `json:"channel_id"`
	Channels   []Channel `json:"channels"`

	PeersStorage *peersStorage `json:"peers_storage"`
}

var gStorage = &storage{
	PeersStorage: newPeersStorage(),
}

func Load(path string) (err error) {
	gStorage.path = path
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return json.NewDecoder(file).Decode(&gStorage)
}

func save() (err error) {
	file, err := os.Create(gStorage.path)
	if err != nil {
		return err
	}
	return json.NewEncoder(file).Encode(gStorage)
}

func Save() (err error) {
	gStorage.Lock()
	defer gStorage.Unlock()
	return save()
}
