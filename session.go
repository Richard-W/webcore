package webcore

import (
	http	"net/http"
	rand	"crypto/rand"
	time	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const cookieLifetime = 3600

var sessions = map[string]*Session {}

type Session struct {
	sessionId	string
	vars		map[string]string
}

func (s *Session) Set (key, val string) {
	if val != "" {
		s.vars[key] = val
	} else {
		delete (s.vars, key)
	}
}

func (s *Session) Get (key string) (string, bool) {
	result, ok := s.vars[key]
	return result, ok
}

func getSession (w http.ResponseWriter, r *http.Request) *Session {
	cookie, err := r.Cookie ("webcore-session-id")
	if err != nil {
		return newSession (w)
	}
	session, ok := sessions[cookie.Value]
	if !ok {
		return newSession (w)
	}
	return session
}

func newSession (w http.ResponseWriter) *Session {
	session := new (Session)
	session.sessionId = generateSessionId ()
	session.vars = map[string]string {}
	go func () {
		time.Sleep (cookieLifetime * time.Second)
	} ()
	http.SetCookie (w, &http.Cookie {
		Name:	"webcore-session-id",
		Value:	session.sessionId,
		MaxAge:	cookieLifetime,
	})
	sessions[session.sessionId] = session
	return session
}


func generateSessionId () string {
	sessionId := ""
	var b [64]byte
	n, err := rand.Read (b[0:64])
	if err != nil {
		panic (err.Error ())
	}
	if n != 64 {
		panic ("sessionId generation failed")
	}
	for i := 0; i < 64; i++ {
		sessionId = sessionId + string (chars[int (b[i]) % len (chars)])
	}
	return sessionId
}
