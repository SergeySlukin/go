package main

import "sessionManager/session"

var globalSession *session.Manager

func init()  {
	globalSession, _ = session.NewManager("memory", "gosessionid", 3600)
}


