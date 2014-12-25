package session

import (
	"time"

	"github.com/go-xweb/httpsession"
	"github.com/lunny/tango"
)

type Sessioner interface {
	SetSession(*httpsession.Session)
}

type Session struct {
	*httpsession.Session
}

func (s *Session) SetSession(session *httpsession.Session) {
	s.Session = session
}

type Sessions struct {
	*httpsession.Manager
}

func NewSessions(sessionTimeout time.Duration) *Sessions {
	sessionMgr := httpsession.Default()
	if sessionTimeout > time.Second {
		sessionMgr.SetMaxAge(sessionTimeout)
	}
	sessionMgr.Run()

	return &Sessions{Manager: sessionMgr}
}

func (itor *Sessions) Handle(ctx *tango.Context) {
	if action := ctx.Action(); ctx != nil {
		if s, ok := action.(Sessioner); ok {
			s.SetSession(itor.Session(ctx.Req(), ctx.ResponseWriter))
		}
	}

	ctx.Next()
}
