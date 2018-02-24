package session

import (
	"net/http"
	"fmt"
	"io"
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"time"
	"log"
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionExist(sid string) bool
	SessionDestroy(sid string) error
	SessionGC(maxLifetime int64)
}

type Manager struct {
	cookieName	string
	provider 	Provider
	maxLifetime	int64
}

var providers = make(map[string]Provider)

func NewManager(providerName, cookieName string, maxLifetime int64) (*Manager, error) {
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", providerName)
	}

	return &Manager{cookieName:cookieName, provider:provider, maxLifetime:maxLifetime}, nil
}

func Register(name string, provider Provider) {
	if provider == nil {
		panic("sesson: Register provider is nil")
	}

	if _, dup := providers[name]; dup {
		panic("session: Register called twice for provider " + name)
	}

	providers[name] = provider
}

func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	sid, _ := manager.getSid(r)
	if sid != "" && manager.provider.SessionExist(sid) {
		log.Println("session is already exist")
		session, _ = manager.provider.SessionRead(sid)
	} else {
		log.Println("create a new session")
		sid = manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{
			Name: 		manager.cookieName,
			Value: 		url.QueryEscape(sid),
			Path: 		"/",
			HttpOnly: 	true,
			MaxAge:		int(manager.maxLifetime)}
		http.SetCookie(w, &cookie)
	}

	return
}

func (manager *Manager) getSid(r *http.Request) (string, error) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return "", err
	}

	// HTTP Request contains cookie for sessionid info.
	return url.QueryUnescape(cookie.Value)
}

func (manager *Manager) SessionExist(r *http.Request) bool {
	sid, err := manager.getSid(r)
	if err != nil {
		return false
	}

	return manager.provider.SessionExist(sid)
}

//Destroy session
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.provider.SessionGC(manager.maxLifetime)
	time.AfterFunc(time.Duration(manager.maxLifetime), func() { manager.GC() })
}