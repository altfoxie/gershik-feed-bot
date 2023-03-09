package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

type promptAuthenticator struct{}

func promptFor(s string) (input string) {
	for {
		fmt.Printf("%s: ", s)
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
	}
}

func (p promptAuthenticator) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	return promptFor("Code: "), nil
}

func (p promptAuthenticator) Phone(ctx context.Context) (string, error) {
	return promptFor("Phone number: "), nil
}

func (p promptAuthenticator) Password(ctx context.Context) (string, error) {
	return promptFor("Password: "), nil
}

func (p promptAuthenticator) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return nil
}

func (p promptAuthenticator) SignUp(ctx context.Context) (auth.UserInfo, error) {
	panic("signup not implemented")
}

func (b *Bot) auth() error {
	return b.client.Auth().IfNecessary(
		bctx,
		auth.NewFlow(
			promptAuthenticator{},
			auth.SendCodeOptions{},
		),
	)
}
