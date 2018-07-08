package fuzzy

// Middleware represent a functional middleware
type Middleware func(CommandHandler) CommandHandler

// MiddlewareChain represents a collection of middleware
type MiddlewareChain struct {
	middleware []Middleware
}

// NewMiddlewareChain creates a new middleware chain
func NewMiddlewareChain(ms ...Middleware) MiddlewareChain {
	return MiddlewareChain{append(([]Middleware)(nil), ms...)}
}

// Then uses all the middleware in the chain on h and then calls h
func (c MiddlewareChain) Then(h CommandHandler) CommandHandler {
	if h == nil {
		return nil
	}

	for i := range c.middleware {
		h = c.middleware[len(c.middleware)-1-i](h)
	}

	return h
}

// Add adds middleware to the chain
func (c *MiddlewareChain) Add(ms ...Middleware) {
	c.middleware = append(c.middleware, ms...)
}
