// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"time"
)

type Session struct {
	id      Id
	maxAge  time.Duration
	manager *Sessions
	rw      http.ResponseWriter
}

func (session *Session) Id() Id {
	return session.id
}

func (session *Session) SetId(id Id) {
	session.id = id
}

func (session *Session) Get(key string) interface{} {
	return session.manager.Store.Get(session.id, key)
}

func (session *Session) Set(key string, value interface{}) {
	session.manager.Store.Set(session.id, key, value)
}

func (session *Session) Del(key string) bool {
	return session.manager.Store.Del(session.id, key)
}

func (session *Session) Release() {
	session.manager.Invalidate(session.rw, session)
}

func (session *Session) IsValid() bool {
	return session.manager.Generator.IsValid(session.id)
}

func (session *Session) SetMaxAge(maxAge time.Duration) {
	session.maxAge = maxAge
}

func (session *Session) Sessions() *Sessions {
	return session.manager
}

func (session *Session) GetSession() *Session {
	return session
}

func (session *Session) InitSession(manager *Sessions, req *http.Request, rw http.ResponseWriter) error {
	id, err := manager.Tracker.Get(req)
	if err != nil {
		return err
	}

	var renew bool

	if !manager.Generator.IsValid(id) {
		id = manager.Generator.Gen(req)
		manager.Tracker.Set(req, rw, id)
		manager.Store.Add(id)
		renew = true
	}

	session.id = id
	session.maxAge = manager.MaxAge
	session.manager = manager
	session.rw = rw

	if renew && manager.OnSessionNew != nil {
		manager.OnSessionNew(session)
	}

	return nil
}

// Sessioner Session interface
type Sessioner interface {
	GetSession() *Session
	InitSession(*Sessions, *http.Request, http.ResponseWriter) error
}

type proxySession struct {
	Session
}

var _ Sessioner = &proxySession{}
