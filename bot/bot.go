package bot

import (
	"context"

	"github.com/charmbracelet/log"

	"github.com/altfoxie/gershik-feed-bot/config"
	"github.com/altfoxie/gershik-feed-bot/storage"
	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"

	tbot "gopkg.in/telebot.v3"
)

// Workaround to reduce boilerplate.
var bctx = context.Background()

type Bot struct {
	client *telegram.Client
	stop   bg.StopFunc
	self   *tg.User
	peers  *peers.Manager

	bot *tbot.Bot
}

func New() *Bot {
	return &Bot{
		client: telegram.NewClient(
			config.User.AppID, config.User.AppHash, telegram.Options{
				SessionStorage: &session.FileStorage{Path: config.User.SessionPath},
			},
		),
	}
}

func (b *Bot) Run() error {
	stop, err := bg.Connect(b.client)
	if err != nil {
		return err
	}
	b.stop = stop

	if err := b.auth(); err != nil {
		return err
	}

	b.peers = peers.Options{
		Storage: storage.PeersStorage(),
		// TODO: implement cache
	}.Build(b.client.API())

	if b.self, err = b.client.Self(bctx); err != nil {
		return err
	}

	name := "@" + b.self.Username
	if b.self.Username == "" {
		name = b.self.FirstName
		if b.self.LastName != "" {
			name += " " + b.self.LastName
		}
	}
	log.Info("Authorized on user account", "name", name)

	bot, err := tbot.NewBot(tbot.Settings{
		Token:     config.Bot.Token,
		Poller:    &tbot.LongPoller{},
		ParseMode: tbot.ModeHTML,
	})
	if err != nil {
		return err
	}
	b.bot = bot
	log.Info("Authorized on bot account", "name", "@"+bot.Me.Username)

	b.registerHandlers()
	go b.bot.Start()
	go b.poller()

	return nil
}
