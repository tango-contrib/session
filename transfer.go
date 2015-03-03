package session

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Transfer provide and set sessionid
type Transfer interface {
	SetMaxAge(maxAge time.Duration)
	Get(req *http.Request) (Id, error)
	Set(req *http.Request, rw http.ResponseWriter, id Id)
	Clear(rw http.ResponseWriter)
}

// CookieRetriever provide sessionid from cookie
type CookieTransfer struct {
	Name     string
	MaxAge   time.Duration
	Lock     sync.Mutex
	Secure   bool
	RootPath string
	Domain   string
}

func NewCookieTransfer(name string, maxAge time.Duration, secure bool, rootPath string) *CookieTransfer {
	return &CookieTransfer{
		Name:     name,
		MaxAge:   maxAge,
		Secure:   secure,
		RootPath: rootPath,
	}
}

func (transfer *CookieTransfer) SetMaxAge(maxAge time.Duration) {
	transfer.MaxAge = maxAge
}

func (transfer *CookieTransfer) Get(req *http.Request) (Id, error) {
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

func (transfer *CookieTransfer) Set(req *http.Request, rw http.ResponseWriter, id Id) {
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
		if transfer.MaxAge > 0 {
			cookie.MaxAge = int(transfer.MaxAge / time.Second)
			//cookie.Expires = time.Now().Add(transfer.maxAge).UTC()
		}

		req.AddCookie(cookie)
	} else {
		cookie.Value = sid
		if transfer.MaxAge > 0 {
			cookie.MaxAge = int(transfer.MaxAge / time.Second)
			//cookie.Expires = time.Now().Add(transfer.maxAge)
		}
	}
	http.SetCookie(rw, cookie)
}

func (transfer *CookieTransfer) Clear(rw http.ResponseWriter) {
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

var _ Transfer = NewCookieTransfer("test", 0, false, "/")

// CookieRetriever provide sessionid from url
/*type UrlTransfer struct {
}

func NewUrlTransfer() *UrlTransfer {
	return &UrlTransfer{}
}

func (transfer *UrlTransfer) Get(req *http.Request) (string, error) {
	return "", nil
}

func (transfer *UrlTransfer) Set(rw http.ResponseWriter, id Id) {

}

var (
	_ Transfer = NewUrlTransfer()
)
*/

//for SWFUpload ...
func NewCookieUrlTransfer(name string, maxAge time.Duration, secure bool, rootPath string) *CookieUrlTransfer {
	return &CookieUrlTransfer{
		CookieTransfer: CookieTransfer{
			Name:     name,
			MaxAge:   maxAge,
			Secure:   secure,
			RootPath: rootPath,
		},
	}
}

type CookieUrlTransfer struct {
	CookieTransfer
}

func (transfer *CookieUrlTransfer) Get(req *http.Request) (Id, error) {
	sessionId := req.URL.Query().Get(transfer.Name)
	if sessionId != "" {
		sessionId, _ = url.QueryUnescape(sessionId)
		return Id(sessionId), nil
	}

	return transfer.CookieTransfer.Get(req)
}
