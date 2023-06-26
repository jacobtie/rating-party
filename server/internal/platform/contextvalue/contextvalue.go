package contextvalue

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
	JWT          *jwt.Token
	RequestID    string
	RequestStart time.Time
	StatusCode   int
	Method       string
	Path         string
	IP           string
	Host         string
	Referrer     string
	Error        bool
}
