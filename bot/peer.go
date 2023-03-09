package bot

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/gotd/td/telegram/peers"
)

func (b *Bot) resolvePeer(from string) (peers.Peer, error) {
	id, _ := strconv.ParseInt(from, 10, 64)
	if id != 0 {
		return b.peers.ResolveChannelID(bctx, id)
	}

	if peer, err := b.peers.Resolve(bctx, from); err == nil {
		return peer, err
	}

	matches := inviteHashRe.FindStringSubmatch(from)
	if len(matches) > 0 {
		hash := matches[len(matches)-1]
		log.Debug("Found invite hash", "hash", hash)
		return b.peers.ImportInvite(bctx, hash)
	}

	return nil, errors.New("peer not resolved")
}
