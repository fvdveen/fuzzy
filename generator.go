package fuzzy

// Generator generates structs that implement all the interfaces required by the bot
type Generator struct {
	contextGenerator      ContextGenerator
	voiceHandlerGenerator VoiceHandlerGenerator
	loggerGenerator       LoggerGenerator
}

// DefaultGenerator is the default generator for the bot
func DefaultGenerator() *Generator {
	g := &Generator{
		contextGenerator:      ContextGenerator(DefaultContext),
		loggerGenerator:       LoggerGenerator(DefaultLogger),
		voiceHandlerGenerator: VoiceHandlerGenerator(DefaultVoiceHandler),
	}

	return g
}

// SetContextGenerator sets the Context generator
func (g *Generator) SetContextGenerator(c ContextGenerator) {
	g.contextGenerator = c
}

// SetVoiceHandlerGenerator sets the VoiceHandler generator
func (g *Generator) SetVoiceHandlerGenerator(v VoiceHandlerGenerator) {
	g.voiceHandlerGenerator = v
}

// SetLoggerGenerator sets the Logger generator
func (g *Generator) SetLoggerGenerator(l LoggerGenerator) {
	g.loggerGenerator = l
}

// Logger creates a new logger
func (g *Generator) Logger(lvl LogLevel) Logger {
	return g.loggerGenerator(lvl)
}
