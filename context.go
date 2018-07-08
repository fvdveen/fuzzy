package fuzzy

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ContextGenerator creates a Context
type ContextGenerator func(context.Context, string, *discordgo.MessageCreate, *Bot, *discordgo.Session, Command) Context

// Context holds all information required by a Command
type Context interface {
	context.Context
	Message() string
	MessageEvent() *discordgo.MessageCreate
	Bot() *Bot
	Session() *discordgo.Session
	Command() Command
	Logger() Logger
	Guild() (*discordgo.Guild, error)

	WithContext(ctx context.Context) Context
	VoiceHandler() (VoiceHandler, error)
	PlaySound(VoiceItem) error

	SendMessage(string)
	SendEmbed(*discordgo.MessageEmbed)
}

type defaultContext struct {
	ctx context.Context

	msg           string
	messageCreate *discordgo.MessageCreate
	bot           *Bot
	sess          *discordgo.Session
	command       Command
}

// DefaultContext is the default context generator
func DefaultContext(ctx context.Context, msg string, mc *discordgo.MessageCreate, b *Bot, sess *discordgo.Session, com Command) Context {
	return &defaultContext{
		ctx:           ctx,
		msg:           msg,
		messageCreate: mc,
		bot:           b,
		sess:          sess,
		command:       com,
	}
}

func (ctx *defaultContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.ctx.Deadline()
}

func (ctx *defaultContext) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

func (ctx *defaultContext) Err() error {
	return ctx.ctx.Err()
}

func (ctx *defaultContext) Value(key interface{}) interface{} {
	return ctx.ctx.Value(key)
}

func (ctx *defaultContext) Message() string {
	return ctx.msg
}

func (ctx *defaultContext) MessageEvent() *discordgo.MessageCreate {
	return ctx.messageCreate
}

func (ctx *defaultContext) Bot() *Bot {
	return ctx.bot
}

func (ctx *defaultContext) Session() *discordgo.Session {
	return ctx.sess
}

func (ctx *defaultContext) Command() Command {
	return ctx.command
}

func (ctx *defaultContext) Logger() Logger {
	return ctx.bot.Generator().Logger(ctx.Bot().Config().LogLevel)
}

func (ctx *defaultContext) Guild() (*discordgo.Guild, error) {
	c, err := ctx.sess.Channel(ctx.messageCreate.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("could not get channel: %v", err)
	}

	g, err := ctx.sess.Guild(c.GuildID)
	if err != nil {
		return nil, fmt.Errorf("could not get guild: %v", err)
	}
	return g, nil
}

func (ctx *defaultContext) WithContext(c context.Context) Context {
	ctx2 := new(defaultContext)
	*ctx2 = *ctx
	ctx2.ctx = c
	return ctx2
}

func (ctx *defaultContext) SendMessage(msg string) {
	_, _ = ctx.sess.ChannelMessageSend(ctx.messageCreate.ChannelID, msg)
}

func (ctx *defaultContext) SendEmbed(e *discordgo.MessageEmbed) {
	_, _ = ctx.sess.ChannelMessageSendEmbed(ctx.messageCreate.ChannelID, e)
}

func (ctx *defaultContext) VoiceHandler() (VoiceHandler, error) {
	g, err := ctx.sess.Channel(ctx.messageCreate.ChannelID)
	if err != nil {
		return nil, err
	}
	return ctx.bot.VoiceHandler(g.GuildID)
}

func (ctx *defaultContext) PlaySound(vi VoiceItem) error {
	c, err := ctx.sess.Channel(ctx.messageCreate.ChannelID)
	if err != nil {
		return err
	}

	vs, err := ctx.bot.GetVoiceState(c.GuildID, ctx.messageCreate.Author.ID)
	if err != nil {
		return err
	}

	ctx.bot.PlaySound(c.GuildID, vs.ChannelID, ctx.messageCreate.ChannelID, vi)

	return nil
}
