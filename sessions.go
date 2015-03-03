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

type Options struct {
	MaxAge           time.Duration
	SessionIdName    string
	Store            Store
	Generator        IdGenerator
	Transfer         Transfer
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
	if opt.Transfer == nil {
		opt.Transfer = NewCookieTransfer(opt.SessionIdName, opt.MaxAge, false, DefaultCookiePath)
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
	if action := ctx.Action(); ctx != nil {
		if s, ok := action.(Sessioner); ok {
			s.SetSession(itor.Session(ctx.Req(), ctx.ResponseWriter))
		}
	}

	ctx.Next()
}

func (manager *Sessions) SetMaxAge(maxAge time.Duration) {
	manager.MaxAge = maxAge
	manager.Transfer.SetMaxAge(maxAge)
	manager.Store.SetMaxAge(maxAge)
}

func (manager *Sessions) SetIdMaxAge(id Id, maxAge time.Duration) {
	manager.Store.SetIdMaxAge(id, maxAge)
}

// TODO:
func (manager *Sessions) Session(req *http.Request, rw http.ResponseWriter) *Session {
	id, err := manager.Transfer.Get(req)
	if err != nil {
		// TODO:
		println("error:", err.Error())
		return nil
	}

	if !manager.Generator.IsValid(id) {
		id = manager.Generator.Gen(req)
		manager.Transfer.Set(req, rw, id)
		manager.Store.Add(id)
	}

	session := &Session{id: id, maxAge: manager.MaxAge, manager: manager}
	// is exist?
	if manager.OnSessionNew != nil {
		manager.OnSessionNew(session)
	}
	return session
}

func (manager *Sessions) Invalidate(rw http.ResponseWriter, session *Session) {
	if manager.OnSessionRelease != nil {
		manager.OnSessionRelease(session)
	}
	manager.Store.Clear(session.id)
	manager.Transfer.Clear(rw)
}

func (manager *Sessions) Run() error {
	return manager.Store.Run()
}
