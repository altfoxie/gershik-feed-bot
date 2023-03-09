package storage

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"github.com/gotd/td/telegram/peers"
)

type peersStorage struct {
	Phones       map[string]peers.Key `json:"phones"`
	Data         peersData            `json:"data"`
	ContactsHash atomic.Int64         `json:"contacts_hash"`
}

type peersData map[peers.Key]peers.Value

func (d peersData) MarshalJSON() ([]byte, error) {
	data := make(map[string]peers.Value, len(d))
	for k, v := range d {
		b, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		data[string(b)] = v
	}
	return json.Marshal(data)
}

func (d *peersData) UnmarshalJSON(b []byte) error {
	data := make(map[string]peers.Value)
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	*d = make(peersData, len(data))
	for k, v := range data {
		var key peers.Key
		if err := json.Unmarshal([]byte(k), &key); err != nil {
			return err
		}
		(*d)[key] = v
	}
	return nil
}

func newPeersStorage() *peersStorage {
	return &peersStorage{
		Phones: make(map[string]peers.Key),
		Data:   make(peersData),
	}
}

func PeersStorage() peers.Storage {
	return gStorage.PeersStorage
}

func (s *peersStorage) Save(ctx context.Context, key peers.Key, value peers.Value) error {
	gStorage.Lock()
	defer gStorage.Unlock()

	s.Data[key] = value
	save()
	return nil
}

func (s *peersStorage) Find(ctx context.Context, key peers.Key) (value peers.Value, found bool, _ error) {
	gStorage.Lock()
	defer gStorage.Unlock()

	value, found = s.Data[key]
	return
}

func (s *peersStorage) SavePhone(ctx context.Context, phone string, key peers.Key) error {
	gStorage.Lock()
	defer gStorage.Unlock()

	s.Phones[phone] = key
	save()
	return nil
}

func (s *peersStorage) FindPhone(ctx context.Context, phone string) (key peers.Key, value peers.Value, found bool, err error) {
	gStorage.Lock()
	defer gStorage.Unlock()

	key, found = s.Phones[phone]
	if !found {
		return
	}
	value, found = s.Data[key]
	return
}

func (s *peersStorage) GetContactsHash(ctx context.Context) (int64, error) {
	return s.ContactsHash.Load(), nil
}

func (s *peersStorage) SaveContactsHash(ctx context.Context, hash int64) error {
	s.ContactsHash.Store(hash)
	return nil
}
