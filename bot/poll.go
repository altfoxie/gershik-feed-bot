package bot

import (
	"strings"
	"time"

	"github.com/altfoxie/gershik-feed-bot/config"
	"github.com/altfoxie/gershik-feed-bot/storage"
	"github.com/charmbracelet/log"
	"github.com/gotd/td/tg"
)

func (b *Bot) poller() {
	for {
		for i, channel := range storage.Channels() {
			// Sleep between channels
			if i != 0 {
				time.Sleep(time.Second * time.Duration(config.Poller.ChannelsInterval))
			}

			log.Info("Polling channel", "channel", channel)

			// Get source channel peer
			peer, err := b.peers.ResolveChannelID(bctx, channel.ChannelID)
			if err != nil {
				log.Error("Failed to resolve peer", "err", err)
				continue
			}

			// Get feed
			feed := storage.FeedByID(channel.FeedID)
			if feed == nil {
				log.Error("Failed to get feed", "id", channel.FeedID)
				storage.DeleteChannel(channel.ID)
				continue
			}

			// Get feed peer
			feedPeer, err := b.peers.ResolveChannelID(bctx, feed.ChannelID)
			if err != nil {
				log.Error("Failed to resolve feed peer", "err", err)
				storage.DeleteChannel(channel.ID)
				continue
			}

			min := channel.LastMessageID
			if min == 0 {
				log.Info("New channel, getting last message", "channel_id", channel.ChannelID)
			}

			// Get messages
			msgsClass, err := b.client.API().MessagesGetHistory(bctx, &tg.MessagesGetHistoryRequest{
				Peer:  peer.InputPeer(),
				Limit: 100,
				MinID: min,
			})
			if err != nil {
				log.Error("Failed to get messages", "err", err)
				continue
			}

			// Fuck Go
			msgs, ok := msgsClass.(*tg.MessagesChannelMessages)
			if !ok {
				log.Error("Failed to cast messages")
				continue
			}

			// No new messages
			if len(msgs.Messages) == 0 {
				log.Info("No new messages", "channel_id", channel.ChannelID)
				continue
			}

			// If it's not a new channel, forward new messages
			if min > 0 {
				// Filtering
				blacklistGroups := make(map[int64]struct{})
				blackListMessages := make(map[int]struct{})

				for _, msg := range msgs.Messages {
					if msg, ok := msg.(*tg.Message); ok {
						for _, filter := range config.Poller.Filters {
							// Check if channel is in filter channels
							var matched bool
							for _, id := range filter.Channels {
								if channel.ChannelID == id {
									matched = true
									break
								}
							}

							// If channel is in filter channels and message contains filter text
							if matched && strings.Contains(strings.ToLower(msg.GetMessage()), strings.ToLower(filter.Text)) {
								log.Warn("Blacklisting message", "message_id", msg.GetID())
								blackListMessages[msg.GetID()] = struct{}{}

								if groupID, ok := msg.GetGroupedID(); ok {
									log.Warn("Blacklisting group", "group_id", groupID)
									blacklistGroups[groupID] = struct{}{}
								}
							}
						}
					}
				}

				// NOTE: ids - message ids, rids - random ids
				ids, rids := []int{}, []int64{}
				// Iterate messages in reverse order
				for i := len(msgs.Messages) - 1; i >= 0; i-- {
					msg := msgs.Messages[i]
					if msg, ok := msg.(*tg.Message); ok {
						// Check if message is blacklisted
						if _, ok := blackListMessages[msg.GetID()]; ok {
							continue
						}

						// Check if message group is blacklisted
						if groupID, ok := msg.GetGroupedID(); ok {
							if _, ok := blacklistGroups[groupID]; ok {
								continue
							}
						}
					}

					i, err := b.client.RandInt64()
					if err != nil {
						log.Fatal("Failed to get random int64", "err", err)
					}

					ids = append(ids, msg.GetID())
					rids = append(rids, i)
				}

				// Forward messages if there are any
				if len(ids) > 0 && len(rids) == len(ids) {
					_, err := b.client.API().MessagesForwardMessages(bctx, &tg.MessagesForwardMessagesRequest{
						FromPeer: peer.InputPeer(),
						ToPeer:   feedPeer.InputPeer(),
						ID:       ids,
						RandomID: rids,
					})
					if err != nil {
						log.Error("Failed to forward messages", "err", err)
						continue
					}

					// Mark messages as unread if needed
					if config.Poller.MarkAsUnread {
						log.Info("Marking messages as unread", "feed_channel_id", feed.ChannelID)
						_, err := b.client.API().MessagesMarkDialogUnread(bctx, &tg.MessagesMarkDialogUnreadRequest{
							Peer:   &tg.InputDialogPeer{Peer: feedPeer.InputPeer()},
							Unread: true,
						})
						if err != nil {
							log.Error("Failed to mark messages as unread", "err", err)
							continue
						}
					}
				}
			}

			// Update last message id
			if !storage.UpdateChannelLastMessageID(channel.ID, msgs.Messages[0].GetID()) {
				log.Error("Failed to update channel last message id")
				continue
			}
		}

		time.Sleep(time.Second * time.Duration(config.Poller.Interval))
	}
}
