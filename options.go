package fuzzy

// OptionFunc sets an option in the bot
type OptionFunc func(*Bot)

// WithConfig sets the bot's config
func WithConfig(c *Config) OptionFunc {
	return func(b *Bot) {
		b.conf = c
	}
}

// WithToken sets the bot's token
func WithToken(t string) OptionFunc {
	return func(b *Bot) {
		b.conf.Token = t
	}
}

// WithPrefix sets the bots prefix
func WithPrefix(p string) OptionFunc {
	return func(b *Bot) {
		b.conf.Prefix = p
	}
}

// WithGenerator sets the bot's generator
func WithGenerator(g *Generator) OptionFunc {
	return func(b *Bot) {
		b.generator = g
	}
}
