package main

import (
	"os"

	"github.com/altfoxie/gershik-feed-bot/bot"
	"github.com/altfoxie/gershik-feed-bot/config"
	"github.com/altfoxie/gershik-feed-bot/storage"
	"github.com/charmbracelet/log"
)

const (
	configPath = "config.yml"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}

	if err := config.Load(configPath); err != nil {
		if os.IsNotExist(err) {
			log.Info("Config file not found, creating new one", "path", configPath)
			if err := config.Save(configPath); err != nil {
				log.Error("Failed to save config file", "err", err)
			}
			log.Print("Please fill it and restart the bot")
			return
		} else {
			log.Fatal("Failed to load config file", "err", err)
		}
	}
	config.Validate()

	if err := storage.Load(config.Poller.StoragePath); err != nil {
		log.Debug("Failed to load storage", "err", err)
	}

	b := bot.New()
	if err := b.Run(); err != nil {
		log.Fatal("Failed to run bot", "err", err)
	}
	log.Info("Bot started")

	select {}
}
