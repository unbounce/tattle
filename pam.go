package main

import (
	"os"
)

type PAMEvent struct {
	Username string
	RemoteHost string
	SessionType string
}

func NewPAMEvent() PAMEvent {
	var p PAMEvent
	p.Username = os.Getenv("PAM_USER")
	p.RemoteHost = os.Getenv("PAM_RHOST")
	p.SessionType = os.Getenv("PAM_TYPE")
	return p
}
