package provides

import (
	"sessionManager/src/session"
)

type Provider interface {
	SessionInit(sid string) (session.Session, error)
	SessionRead(sid string) (session.Session, error)
	SessionDestroy(sid string) (session.Session, error)
	SessionGC(maxLifeTime uint64)
}

var providesMap = make(map[string]Provider)

func Register(name string, provider Provider)  {
	if provider == nil {
		panic("session: Register provider is nil")
	}

	if _, dup := providesMap[name]; dup {
		panic("session: Register called twice for provider " + name)
	}

	providesMap[name] = provider
}

func GetProvider(name string) (Provider, bool) {
	provider, ok := providesMap[name]
	if !ok {
		return nil, false
	}
	return provider, true
}

