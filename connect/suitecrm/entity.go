package suitecrm

import (
	"sync"
	"time"

	"github.com/asm2me/asterlink/connect"
	log "github.com/sirupsen/logrus"
)

type entity struct {
	ID string `json:"id"`

	Dir           connect.Direction   `json:"dir"`
	DID           string              `json:"did"`
	CID           string              `json:"cid"`
	IsAnswered    bool                `json:"answered"`
	TimeStamp     time.Time           `json:"time"`
	Relationships map[string]relation `json:"relations"`
	exts          sync.Map

	log *log.Entry
	mux sync.Mutex
}

func (e *entity) isRegistred() bool {
	if e.ID != "" {
		return true
	}

	e.mux.Lock()
	e.mux.Unlock()

	if e.ID == "" {
		return false
	}

	return true
}
