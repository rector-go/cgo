package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provier     Provider
	maxLifetime int64
}

func NewManager(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provier: provider, cookieName: cookieName, maxLifetime: maxLifeTime}, nil
}

func (p *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (p *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	p.lock.Lock()
	defer p.lock.Unlock()
	cookie, err := r.Cookie(p.cookieName)
	if err != nil || cookie.Value == "" {
		sid := p.sessionId()
		session, _ = p.provier.SessionInit(sid)
		cookie := http.Cookie{Name: p.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(p.maxLifetime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = p.provier.SessionRead(sid)
	}
	return
}

func (p *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(p.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		p.lock.Lock()
		defer p.lock.Unlock()
		_ = p.provier.SessionDestroy(cookie.Value)
		cookie := http.Cookie{Name: p.cookieName, Path: "/", HttpOnly: true, Expires: time.Now(), MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (p *Manager) GC() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.provier.SessionGC(p.maxLifetime)
	time.AfterFunc(time.Duration(p.maxLifetime), func() {
		p.GC()
	})
}
