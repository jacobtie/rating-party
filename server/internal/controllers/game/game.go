package game

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jacobtie/rating-party/server/internal/config"
	"github.com/jacobtie/rating-party/server/internal/platform/db"
)

type Controller struct {
	cfg *config.Config
	db  *db.DB
}

func NewController(cfg *config.Config, db *db.DB) *Controller {
	return &Controller{
		cfg: cfg,
		db:  db,
	}
}

func (c *Controller) GenerateGameCode() string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := [5]int{rand.Intn(36), rand.Intn(36), rand.Intn(36), rand.Intn(36)}
	code := [5]string{}
	for i, num := range nums {
		if num < 10 {
			code[i] = strconv.Itoa(num)
		} else {
			code[i] = string(rune(num + 55)) // convert to capital letter, 10 => 65 (A), 35 => 90 (Z)
		}
	}
	return strings.Join(code[:], "")
}
