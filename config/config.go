// Package config provides a simple configuration system.
package config

import (
	"os"

	"github.com/charmbracelet/log"

	"gopkg.in/yaml.v2"
)

type config struct {
	Bot    bot    `json:"bot" yaml:"bot"`
	User   user   `json:"user" yaml:"user"`
	Poller poller `json:"poller" yaml:"poller"`
}

type bot struct {
	Token     string  `json:"token" yaml:"token"`
	Whitelist []int64 `json:"whitelist" yaml:"whitelist"`
}

type user struct {
	SessionPath string `json:"session_path" yaml:"session_path"`
	AppID       int    `json:"app_id" yaml:"app_id"`
	AppHash     string `json:"app_hash" yaml:"app_hash"`
}

type poller struct {
	StoragePath      string   `json:"storage_path" yaml:"storage_path"`
	Interval         int      `json:"interval" yaml:"interval"`
	ChannelsInterval int      `json:"channels_interval" yaml:"channels_interval"`
	MarkAsUnread     bool     `json:"mark_as_unread" yaml:"mark_as_unread"`
	Filters          []filter `json:"filters" yaml:"filters"`
}

type filter struct {
	Channels []int64 `json:"channels" yaml:"channels"`
	Text     string  `json:"text" yaml:"text"`
}

var (
	Config = defaultConfig()
	Bot    = &Config.Bot
	User   = &Config.User
	Poller = &Config.Poller
)

func defaultConfig() config {
	return config{
		User: user{
			SessionPath: "session.json",
		},
		Poller: poller{
			StoragePath:      "storage.json",
			Interval:         60,
			ChannelsInterval: 10,
		},
	}
}

func Load(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return yaml.NewDecoder(file).Decode(&Config)
}

func Save(path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(file).Encode(Config)
}

func Validate() {
	if Bot.Token == "" {
		log.Fatal("Empty token")
	}
	if User.AppID <= 0 {
		log.Fatal("Invalid app_id", "app_id", User.AppID)
	}
	if User.AppHash == "" {
		log.Fatal("Empty app_hash")
	}
	if User.SessionPath == "" {
		log.Fatal("Empty session_path")
	}
	if Poller.StoragePath == "" {
		log.Fatal("Empty storage_path")
	}
	if Poller.Interval <= 0 {
		log.Fatal("Invalid polling interval", "interval", Poller.Interval)
	}
	if Poller.ChannelsInterval <= 0 {
		log.Fatal("Invalid channels polling interval", "interval", Poller.ChannelsInterval)
	}
}
