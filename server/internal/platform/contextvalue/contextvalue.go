package contextvalue

import "time"

type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
	RequestID    string
	RequestStart time.Time
	StatusCode   int
	RequestBody  string
	Method       string
	Path         string
	IP           string
	Host         string
	Referrer     string
	Error        bool
}
