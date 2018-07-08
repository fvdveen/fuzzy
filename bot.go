package fuzzy

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Bot represents a discord bot
type Bot struct {
	conf *Config
	sess *discordgo.Session

	commands      []Command
	voiceHandlers map[string]VoiceHandler
	voiceMu       sync.RWMutex
	generator     *Generator
	middleware    MiddlewareChain
}

// New creates a new bot
func New(opts ...OptionFunc) (*Bot, error) {
	b := &Bot{
		conf:          &Config{},
		voiceHandlers: make(map[string]VoiceHandler),
		commands:      []Command{},
		generator:     DefaultGenerator(),
		middleware:    NewMiddlewareChain(),
	}

	for _, opt := range opts {
		opt(b)
	}

	sess, err := discordgo.New("Bot " + b.conf.Token)
	if err != nil {
		return nil, err
	}
	b.sess = sess

	b.initHandlers()

	return b, nil
}

// Open opens the discord session
func (b *Bot) Open() error {
	return b.sess.Open()
}

// Close closes the discord session
func (b *Bot) Close() error {
	for _, vh := range b.voiceHandlers {
		vh.Stop()
	}
	return b.sess.Close()
}

// Session returns the discordgo session of the bot
func (b *Bot) Session() *discordgo.Session {
	return b.sess
}

// Config gives the bot's config
func (b *Bot) Config() *Config {
	return b.conf
}

// RegisterCommand registers a command with the bot
func (b *Bot) RegisterCommand(cs ...Command) error {
	for _, c := range cs {
		for _, com := range b.commands {
			if c.Name() == com.Name() {
				return ErrDuplicateCommand
			}
		}
		b.commands = append(b.commands, c)
	}

	return nil
}

// RegisterHandler adds a discordgo handler to the bot
func (b *Bot) RegisterHandler(hs ...interface{}) {
	for _, h := range hs {
		b.sess.AddHandler(h)
	}
}

// UseMiddleware will call the middleware with the context before command handlers are called
func (b *Bot) UseMiddleware(ms ...Middleware) {
	b.middleware.Add(ms...)
}

// Commands returns a copy of all the bot's commands
func (b *Bot) Commands() []Command {
	cs := append(([]Command)(nil), b.commands...)
	return cs
}

// VoiceHandler returns the voicehandler for the given guildID
// if it does not exist it will be created
func (b *Bot) VoiceHandler(gid string) (VoiceHandler, error) {
	b.voiceMu.RLock()
	defer b.voiceMu.RUnlock()

	vh, ok := b.voiceHandlers[gid]
	if !ok {
		return nil, ErrVoiceHandlerNotExists
	}

	return vh, nil
}

// PlaySound puts the VoiceItem in the queue of the voicehandler of the given guild
func (b *Bot) PlaySound(gid, voiceChanID, textChanID string, vi VoiceItem) {
	b.voiceMu.Lock()
	defer b.voiceMu.Unlock()

	vh, ok := b.voiceHandlers[gid]
	if !ok {
		b.voiceHandlers[gid] = b.generator.voiceHandlerGenerator(b.sess, b, gid, voiceChanID, textChanID, vi)
		return
	}
	vh.Play(vi)
}

// GetVoiceState returns the voicestate of the given userID
func (b *Bot) GetVoiceState(gid, userID string) (*discordgo.VoiceState, error) {
	g, err := b.sess.State.Guild(gid)
	if err != nil {
		return nil, err
	}
	for _, vs := range g.VoiceStates {
		if vs.UserID == userID {
			return vs, nil
		}
	}
	return nil, ErrUnknownVoiceState
}

// DeleteVoiceHandler will remove the voice handler from the bot
func (b *Bot) DeleteVoiceHandler(gid string) {
	b.voiceMu.Lock()
	defer b.voiceMu.Unlock()

	delete(b.voiceHandlers, gid)
}

// Generator returns the bot's generator
func (b *Bot) Generator() *Generator {
	return b.generator
}
