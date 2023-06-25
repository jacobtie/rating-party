package web

type Middleware func(Handler) Handler

func wrapMiddleware(handler Handler, mw []Middleware) Handler {
	currHandler := handler
	for _, mwFunc := range mw {
		if mwFunc != nil {
			currHandler = mwFunc(currHandler)
		}
	}
	return currHandler
}
