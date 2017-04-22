// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"time"

	"github.com/lunny/tango"
)

const (
	DefaultMaxAge        = 30 * time.Minute
	DefaultSessionIdName = "SESSIONID"
	DefaultCookiePath    = "/"
)

type Sessions struct {
	Options
}

type Sessionser interface {
	SetSessions(*Sessions)
}

type Options struct {
	MaxAge           time.Duration
	SessionIdName    string
	Store            Store
	Generator        IdGenerator
	Tracker          Tracker
	OnSessionNew     func(*Session)
	OnSessionRelease func(*Session)
}

func preOptions(opts []Options) Options {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.MaxAge == 0 {
		opt.MaxAge = DefaultMaxAge
	}
	if opt.Store == nil {
		opt.Store = NewMemoryStore(opt.MaxAge)
	}
	if opt.SessionIdName == "" {
		opt.SessionIdName = DefaultSessionIdName
	}
	if opt.Generator == nil {
		opt.Generator = NewSha1Generator(string(GenRandKey(16)))
	}
	if opt.Tracker == nil {
		opt.Tracker = NewCookieTracker(opt.SessionIdName, 0, false, DefaultCookiePath)
	}
	return opt
}

func New(opts ...Options) *Sessions {
	opt := preOptions(opts)
	sessions := &Sessions{
		Options: opt,
	}

	sessions.Run()

	return sessions
}

func Default() *Sessions {
	return NewWithTimeout(DefaultMaxAge)
}

func NewWithTimeout(maxAge time.Duration) *Sessions {
	return New(Options{MaxAge: maxAge})
}

func (itor *Sessions) Handle(ctx *tango.Context) {
	if action := ctx.Action(); action != nil {
		if s, ok := action.(Sessionser); ok {
			s.SetSessions(itor)
		}

		if s, ok := action.(Sessioner); ok {
			err := s.InitSession(itor, ctx.Req(), ctx.ResponseWriter)
			if err != nil {
				ctx.Result = err
				return
			}
		}
	}

	ctx.Next()
}

func (manager *Sessions) SessionFromID(id Id) *Session {
	return &Session{
		id:      id,
		maxAge:  manager.MaxAge,
		manager: manager,
	}
}

func (manager *Sessions) SetMaxAge(maxAge time.Duration) {
	manager.MaxAge = maxAge
	manager.Tracker.SetMaxAge(maxAge)
	manager.Store.SetMaxAge(maxAge)
}

func (manager *Sessions) SetIdMaxAge(id Id, maxAge time.Duration) {
	manager.Store.SetIdMaxAge(id, maxAge)
}

func (manager *Sessions) Exist(id Id) bool {
	return manager.Store.Exist(id)
}

func (manager *Sessions) Session(req *http.Request, rw http.ResponseWriter) *Session {
	var sess Session
	sess.InitSession(manager, req, rw)
	return &sess
}

func (manager *Sessions) Invalidate(rw http.ResponseWriter, session *Session) {
	if manager.OnSessionRelease != nil {
		manager.OnSessionRelease(session)
	}
	manager.Store.Clear(session.id)
	manager.Tracker.Clear(rw)
}

func (manager *Sessions) Run() error {
	return manager.Store.Run()
}
