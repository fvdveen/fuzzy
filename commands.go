package fuzzy

import "fmt"

// Command is a action performed by the bot triggered by the string returned by Name
type Command interface {
	CommandHandler

	// Name is not only the command name but also the trigger
	Name() string
	Description() string
}

// CommandHandler handles a command by the client
type CommandHandler interface {
	Handle(Context)
}

// CommandHandlerFunc implements CommandHandler on a function
type CommandHandlerFunc func(Context)

// Handle implements CommandHandler
func (com CommandHandlerFunc) Handle(c Context) {
	com(c)
}

// textCommand is a implementation of command
type textCommand struct {
	name        string
	description string

	run func(Context)
}

// NewCommand creates a new Command
func NewCommand(n, d string, h func(Context)) Command {
	return &textCommand{
		name:        n,
		description: d,
		run:         h,
	}
}

// Name gives the name of the command
func (c textCommand) Name() string {
	return c.name
}

// Description gives the description of the command
func (c textCommand) Description() string {
	return c.description
}

// Run runs the command
func (c textCommand) Handle(ctx Context) {
	c.run(ctx)
}

// HelpCommand is a standard help command
func HelpCommand(t, d string) Command {
	return NewCommand("help", "Shows all commands", func(ctx Context) {
		msg := fmt.Sprintf("%s - Commands:", t)
		for _, com := range ctx.Bot().Commands() {
			msg = fmt.Sprintf("%s\n`%s%s` %s", msg, ctx.Bot().Config().Prefix, com.Name(), com.Description())
		}

		ctx.SendMessage(msg)
	})
}
