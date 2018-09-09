package session

import (
	"sync"
	"fmt"
	"sessionManager/src/provides"
	"io"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"
)

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provides    provides.Provider
	maxLifeTime uint64
}

func NewManager(providerName, cookieName string, maxLifeTime uint64) (*Manager, error) {
	provider, ok := provides.GetProvider(providerName)
	if !ok {
		return nil, fmt.Errorf("session: unknows provider %q (forgottent import?)", providerName)
	}
	return &Manager{provides: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

func (m *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (m *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	m.lock.Lock()
	defer m.lock.Unlock()
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == ""{
		sid := m.sessionId()
		session, _ = m.provides.SessionInit(sid)
		cookie := http.Cookie{Name: m.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(m.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = m.provides.SessionRead(sid)
	}
	return
}

func (m *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request)  {
	cookie, err := r.Cookie(m.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		m.lock.Lock()
		defer m.lock.Unlock()
		m.provides.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: m.cookieName, Path: "/", HttpOnly: true, Expires:expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (m *Manager) SessionGC()  {
	m.lock.Unlock()
	defer m.lock.Unlock()
	m.provides.SessionGC(m.maxLifeTime)
	time.AfterFunc(time.Duration(m.maxLifeTime), func() {
		m.GC()
	})
}
