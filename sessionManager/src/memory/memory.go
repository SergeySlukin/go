package memory

import (
	"time"
	"container/list"
	"sync"
	"sessionManager/src/session"
	"sessionManager/src/provides"
)

var memoryProvider = Provider{list: list.New()}

func init() {
	memoryProvider.sessions = make(map[string]*list.Element, 0)
	provides.Register("memory", memoryProvider)
}

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

func NewSessionStore(sid string) *SessionStore {
	return &SessionStore{
		sid:          sid,
		timeAccessed: time.Now(),
		value:        make(map[interface{}]interface{}),
	}
}

func (s *SessionStore) Set(k, v interface{}) {
	memoryProvider.SessionUpdate(s.sid)
	s.value[k] = v
}

func (s *SessionStore) Get(k interface{}) interface{} {
	memoryProvider.SessionUpdate(s.sid)
	if v, ok := s.value[k]; ok {
		return v
	}
	return nil
}

func (s *SessionStore) Delete(k interface{}) {
	memoryProvider.SessionUpdate(s.sid)
	if v, ok := s.value[k]; ok {
		delete(s.value, v)
	}
}

func (s *SessionStore) SessionID() string {
	return s.sid
}

type Provider struct {
	mu       sync.Mutex
	list     *list.List
	sessions map[string]*list.Element
}

func (p Provider) SessionInit(sid string) session.Session {
	p.mu.Lock()
	defer p.mu.Unlock()
	newSession := NewSessionStore(sid)
	element := p.list.PushBack(newSession)
	p.sessions[sid] = element
	return newSession
}

func (p Provider) SessionRead(sid string) session.Session {
	element, ok := p.sessions[sid]
	if ok {
		return element.Value.(*SessionStore)
	} else {
		newSession := p.SessionInit(sid)
		return newSession
	}
}

func (p Provider) SessionDestroy(sid string) {
	if element, ok := p.sessions[sid]; ok {
		delete(p.sessions, sid)
		p.list.Remove(element)
	}
}

func (p Provider) SessionGC(maxLifeTime uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for {
		element := p.list.Back()
		if element == nil {
			break
		}

		if element.Value.(*SessionStore).timeAccessed.Unix() + int64(maxLifeTime) < time.Now().Unix() {
			p.list.Remove(element)
			delete(p.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (p Provider) SessionUpdate(sid string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if element, ok := p.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		p.list.MoveToFront(element)
	}
	return
}
