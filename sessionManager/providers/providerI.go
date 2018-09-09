package providers

import "sessionManager/session"

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) (Session, error)
	SessionGC(maxLifeTime uint64)
}

