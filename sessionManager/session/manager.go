package session

import (
	"sync"
	"fmt"
)

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifeTime uint64
}

func NewManager(providerName, cookieName string, maxLifeTime uint64) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknows provider %q (forgottent import?)", providerName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}
