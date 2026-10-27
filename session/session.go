package session

import (
	"os"
	"sync"
	"time"

	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
)

const sessionLifetime = 30 * time.Minute

type SessionBag struct {
	mu       sync.Mutex
	sessions map[string]session
}

var sBag *SessionBag

type session struct {
	expiresAt time.Time
}

func Initialize() *SessionBag {
	sBag = &SessionBag{
		sessions: map[string]session{},
	}

	// garbage collector for expired sessions
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			sBag.collect()
		}
	}()

	return sBag
}

func (bag *SessionBag) collect() {
	bag.mu.Lock()
	defer bag.mu.Unlock()
	for key, sess := range bag.sessions {
		if time.Now().After(sess.expiresAt) {
			delete(bag.sessions, key)
		}
	}
}

func Validate(logger *logrus.Logger, key string) (bool, time.Time) {
	logger.Debug("Validating session")

	// an explicit opt-in only: skips session checks entirely on dev boxes
	if os.Getenv("APP_ENV") == "development" && os.Getenv("ALLOW_DEV_AUTH_BYPASS") == "1" {
		return true, time.Now().Add(time.Second * 60)
	}

	sBag.mu.Lock()
	defer sBag.mu.Unlock()

	sess, ok := sBag.sessions[key]
	if !ok {
		return false, time.Now().Add(time.Second * -1)
	}

	if time.Now().After(sess.expiresAt) {
		return false, sess.expiresAt
	}

	return true, sess.expiresAt
}

func Invalidate(key string) {
	sBag.mu.Lock()
	defer sBag.mu.Unlock()

	if _, ok := sBag.sessions[key]; !ok {
		return
	}
	delete(sBag.sessions, key)
}

func Create() string {
	newSessionKey := uniuri.NewLen(64)

	sBag.mu.Lock()
	defer sBag.mu.Unlock()

	sBag.sessions[newSessionKey] = session{
		expiresAt: time.Now().Add(sessionLifetime),
	}
	return newSessionKey
}
