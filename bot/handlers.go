package bot

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/altfoxie/gershik-feed-bot/config"
	"github.com/altfoxie/gershik-feed-bot/storage"
	"github.com/charmbracelet/log"
	tbot "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func (b *Bot) registerHandlers() {
	b.bot.Use(middleware.Whitelist(config.Bot.Whitelist...))
	b.bot.Handle("/start", b.onStart)
	b.bot.Handle("/help", b.onHelp)
	b.bot.Handle("/addfeed", b.onAddFeed)
	b.bot.Handle("/add", b.onAdd)
	b.bot.Handle("/ls", b.onList)
	b.bot.Handle(tbot.OnText, b.onText)

	b.bot.Handle(callback("channel"), b.onAddCallback)
}

func callback(unique string) tbot.CallbackEndpoint {
	return &tbot.Btn{Unique: unique}
}

func (b *Bot) onStart(ctx tbot.Context) error {
	return ctx.Send("Hello, world!")
}

const helpText = `<b>Commands:</b>
/addfeed &lt;link/username/tg id&gt; - add feed
/add &lt;link/username/tg id&gt; - add channel to feed
/ls - list feeds and channels
/help - show this message`

func (b *Bot) onHelp(ctx tbot.Context) error {
	return ctx.Send(helpText)
}

var inviteHashRe = regexp.MustCompile(`(?i)(https?:\/\/)?(t(elegram)?\.me|telegram\.org)(\/joinchat)?\/\+?([a-zA-Z0-9\_-]{5,32})`)

func (b *Bot) onAddFeed(ctx tbot.Context) error {
	peer, err := b.resolvePeer(ctx.Data())
	if err != nil {
		log.Error("Failed to resolve feed peer", "err", err)
		return ctx.Send("Failed to resolve peer")
	}

	if storage.FeedByChannelID(peer.ID()) != nil {
		return ctx.Send("This feed is already added")
	}

	storage.CreateFeed(peer.ID(), peer.VisibleName())
	log.Info("Feed added", "peer", peer, "id", peer.ID(), "name", peer.VisibleName())
	return ctx.Send(fmt.Sprintf("Feed \"%s\" added", peer.VisibleName()))
}

func (b *Bot) onDeleteFeed(ctx tbot.Context, idStr string) error {
	id, _ := strconv.ParseInt(idStr, 10, 64)

	feed := storage.FeedByID(id)
	if feed == nil {
		return ctx.Send("Feed not found")
	}

	storage.DeleteFeed(id)
	log.Info("Feed deleted", "id", id)
	return ctx.Send(fmt.Sprintf("Feed \"%s\" deleted", feed.Name))
}

func (b *Bot) onAdd(ctx tbot.Context) error {
	if ctx.Data() == "" {
		return ctx.Send("Send a link to the channel you want to add")
	}

	if len(storage.Feeds()) == 0 {
		return ctx.Send("There are no feeds added yet. Use /addfeed to add one")
	}

	peer, err := b.resolvePeer(ctx.Data())
	if err != nil {
		log.Error("Failed to resolve channel peer", "err", err)
		return ctx.Send("Failed to resolve peer")
	}

	if storage.ChannelByChannelID(peer.ID()) != nil {
		return ctx.Send("This channel is already added")
	}

	markup := &tbot.ReplyMarkup{}
	rows := []tbot.Row{}
	for _, feed := range storage.Feeds() {
		rows = append(rows, markup.Row(
			markup.Data(feed.Name, "channel",
				fmt.Sprintf("%d:%d", peer.ID(), feed.ID)),
		))
	}
	markup.Inline(rows...)

	return ctx.Send("Select a feed to add this channel to", markup)
}

func (b *Bot) onAddCallback(ctx tbot.Context) error {
	parts := strings.SplitN(ctx.Data(), ":", 2)
	if len(parts) != 2 {
		return errors.New("invalid callback data")
	}

	channelID, _ := strconv.ParseInt(parts[0], 10, 64)
	feedID, _ := strconv.ParseInt(parts[1], 10, 64)

	if storage.ChannelByChannelID(channelID) != nil {
		return ctx.Edit("This channel is already added")
	}

	if storage.FeedByID(feedID) == nil {
		log.Warn("Feed not found", "id", feedID)
		return ctx.Edit("This feed does not exist")
	}

	peer, err := b.peers.ResolveChannelID(bctx, channelID)
	if err != nil {
		log.Error("Failed to resolve peer", "err", err)
		return ctx.Edit("Failed to resolve peer")
	}

	storage.CreateChannel(channelID, feedID, peer.VisibleName())
	log.Info(
		"Channel added",
		"peer", peer, "id", peer.ID(), "name", peer.VisibleName(),
		"feed", feedID,
	)

	return ctx.Edit(fmt.Sprintf("Channel \"%s\" added to feed", peer.VisibleName()))
}

func (b *Bot) onDeleteChannel(ctx tbot.Context, idStr string) error {
	id, _ := strconv.ParseInt(idStr, 10, 64)

	channel := storage.ChannelByID(id)
	if channel == nil {
		return ctx.Send("Channel not found")
	}

	storage.DeleteChannel(id)
	log.Info("Channel deleted", "id", id)
	return ctx.Send(fmt.Sprintf("Channel \"%s\" deleted", channel.Name))
}

func (b *Bot) onList(ctx tbot.Context) error {
	feedIDs := []int64{}
	m := make(map[int64][]storage.Channel)
	for _, feed := range storage.Feeds() {
		m[feed.ID] = []storage.Channel{}
		feedIDs = append(feedIDs, feed.ID)
	}

	sort.Slice(feedIDs, func(i, j int) bool {
		return feedIDs[i] < feedIDs[j]
	})

	for _, feed := range storage.Channels() {
		if channels, ok := m[feed.FeedID]; ok {
			m[feed.FeedID] = append(channels, feed)
		}
	}

	var text string
	for _, feedID := range feedIDs {
		channels := m[feedID]
		sort.Slice(channels, func(i, j int) bool {
			return channels[i].ID < channels[j].ID
		})

		feed := storage.FeedByID(feedID)

		text += fmt.Sprintf("\n<b>%s</b> (/delfeed_%d)\n", feed.Name, feed.ID)
		for _, channel := range channels {
			text += fmt.Sprintf(" - %s - <code>%d</code> (/delchannel_%d)\n", channel.Name, channel.ChannelID, channel.ID)
		}
	}

	if text == "" {
		text = "There are no channels added yet"
	}

	return ctx.Send(text)
}

func (b *Bot) onText(ctx tbot.Context) error {
	switch {
	case strings.HasPrefix(ctx.Text(), "/delfeed_"):
		return b.onDeleteFeed(ctx, strings.TrimPrefix(ctx.Text(), "/delfeed_"))
	case strings.HasPrefix(ctx.Text(), "/delchannel_"):
		return b.onDeleteChannel(ctx, strings.TrimPrefix(ctx.Text(), "/delchannel_"))
	}
	return nil
}
