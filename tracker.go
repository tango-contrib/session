package session

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Tracker provide and set sessionid
type Tracker interface {
	SetMaxAge(maxAge time.Duration)
	Get(req *http.Request) (Id, error)
	Set(req *http.Request, rw http.ResponseWriter, id Id)
	Clear(rw http.ResponseWriter)
}

// CookieTracker provide sessionid from cookie
type CookieTracker struct {
	Name     string
	MaxAge   time.Duration
	Lock     sync.Mutex
	Secure   bool
	RootPath string
	Domain   string
}

func NewCookieTracker(name string, maxAge time.Duration, secure bool, rootPath string) *CookieTracker {
	return &CookieTracker{
		Name:     name,
		MaxAge:   maxAge,
		Secure:   secure,
		RootPath: rootPath,
	}
}

func (transfer *CookieTracker) SetMaxAge(maxAge time.Duration) {
	transfer.MaxAge = maxAge
}

func (transfer *CookieTracker) Get(req *http.Request) (Id, error) {
	cookie, err := req.Cookie(transfer.Name)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", err
	}
	if cookie.Value == "" {
		return Id(""), nil
	}
	id, _ := url.QueryUnescape(cookie.Value)
	return Id(id), nil
}

func (transfer *CookieTracker) Set(req *http.Request, rw http.ResponseWriter, id Id) {
	sid := url.QueryEscape(string(id))
	transfer.Lock.Lock()
	defer transfer.Lock.Unlock()
	cookie, _ := req.Cookie(transfer.Name)
	if cookie == nil {
		cookie = &http.Cookie{
			Name:     transfer.Name,
			Value:    sid,
			Path:     transfer.RootPath,
			Domain:   transfer.Domain,
			HttpOnly: true,
			Secure:   transfer.Secure,
		}
		/*if transfer.MaxAge > 0 {
			cookie.MaxAge = int(transfer.MaxAge / time.Second)
		}*/

		req.AddCookie(cookie)
	} else {
		cookie.Value = sid
		/*if transfer.MaxAge > 0 {
			cookie.MaxAge = int(transfer.MaxAge / time.Second)
		}*/
	}
	http.SetCookie(rw, cookie)
}

func (transfer *CookieTracker) Clear(rw http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     transfer.Name,
		Path:     transfer.RootPath,
		Domain:   transfer.Domain,
		HttpOnly: true,
		Secure:   transfer.Secure,
		Expires:  time.Date(0, 1, 1, 0, 0, 0, 0, time.Local),
		MaxAge:   -1,
	}
	http.SetCookie(rw, &cookie)
}

var _ Tracker = NewCookieTracker("test", 0, false, "/")

// UrlTracker provide sessionid from url
type UrlTracker struct {
	Key         string
	ReplaceLink bool
}

func NewUrlTracker(key string, replaceLink bool) *UrlTracker {
	return &UrlTracker{key, replaceLink}
}

func (tracker *UrlTracker) Get(req *http.Request) (Id, error) {
	sessionId := req.URL.Query().Get(tracker.Key)
	if sessionId != "" {
		sessionId, _ = url.QueryUnescape(sessionId)
		return Id(sessionId), nil
	}

	return Id(""), nil
}

func (tracker *UrlTracker) Set(req *http.Request, rw http.ResponseWriter, id Id) {
	if tracker.ReplaceLink {

	}
}

func (tracker *UrlTracker) SetMaxAge(maxAge time.Duration) {

}

func (tracker *UrlTracker) Clear(rw http.ResponseWriter) {
}

var (
	_ Tracker = NewUrlTracker("id", false)
)

//for SWFUpload ...
func NewCookieUrlTracker(name string, maxAge time.Duration, secure bool, rootPath string) *CookieUrlTracker {
	return &CookieUrlTracker{
		CookieTracker: CookieTracker{
			Name:     name,
			MaxAge:   maxAge,
			Secure:   secure,
			RootPath: rootPath,
		},
	}
}

type CookieUrlTracker struct {
	CookieTracker
}

func (tracker *CookieUrlTracker) Get(req *http.Request) (Id, error) {
	sessionId := req.URL.Query().Get(tracker.Name)
	if sessionId != "" {
		sessionId, _ = url.QueryUnescape(sessionId)
		return Id(sessionId), nil
	}

	return tracker.CookieTracker.Get(req)
}
