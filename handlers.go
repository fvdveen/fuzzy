package fuzzy

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) initHandlers() {
	b.RegisterHandler(b.messageHandler())
}

func (b *Bot) messageHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if s.State.User.ID == m.Author.ID {
			return
		}
		if !strings.HasPrefix(m.Content, b.conf.Prefix) {
			return
		}
		msg := strings.TrimPrefix(m.Content, b.conf.Prefix)
		for _, com := range b.commands {
			if msg == com.Name() || strings.HasPrefix(msg, com.Name()+" ") {
				msg = strings.TrimSpace(strings.TrimPrefix(msg, com.Name()))
				ctx := b.generator.contextGenerator(context.Background(), msg, m, b, s, com)
				b.middleware.Then(com).Handle(ctx)
				return
			}
		}
	}
}
