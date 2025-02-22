package session

import (
	"os"
	"time"

	"github.com/dchest/uniuri"
)

type SessionBag struct {
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
		for key, sess := range sBag.sessions {
			if time.Now().After(sess.expiresAt) {
				Invalidate(key)
			}
		}
		time.Sleep(1 * time.Minute)
	}()

	return sBag
}

func Validate(key string) (bool, time.Time) {
	appMode := os.Getenv("APP_ENV")
	if appMode == "development" {
		return true, time.Now()
	}

	sess, ok := sBag.sessions[key]
	if !ok {
		return false, time.Now()
	}

	if time.Now().After(sess.expiresAt) {
		return false, sess.expiresAt
	}

	return true, sess.expiresAt
}

func Invalidate(key string) {
	_, ok := sBag.sessions[key]
	if !ok {
		return
	}
	delete(sBag.sessions, key)
}

func Create() string {
	newSessionKey := uniuri.NewLen(64)
	sBag.sessions[newSessionKey] = session{
		expiresAt: time.Now().Add(time.Minute * 30),
	}
	return newSessionKey
}
