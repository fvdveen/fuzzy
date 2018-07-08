package middleware

import (
	"github.com/fvdveen/fuzzy"
)

func Logging() fuzzy.Middleware {
	return func(next fuzzy.CommandHandler) fuzzy.CommandHandler {
		return fuzzy.CommandHandlerFunc(func(ctx fuzzy.Context) {
			ctx.Logger().WithFields(map[string]interface{}{
				"user":    ctx.MessageEvent().Author.Username,
				"user-ID": ctx.MessageEvent().Author.ID,
			}).Infof("[%s] %s", ctx.Command().Name(), ctx.Message())
			next.Handle(ctx)
		})
	}
}
