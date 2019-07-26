package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"text/template"
	"time"
)

var (
	// MegaManager to be used throughout application to manage sessions.
	MegaManager *Manager
)

// Manager struct defines a global session manager.
type Manager struct {
	cookieName  string     // private cookiename
	lock        sync.Mutex // protects session
	provider    Provider   // session creation/deletion
	maxlifetime int64      // shelf-life of a manager
}

// Provider implements session session storage management. Sessions should be persisted either in memory, in the file system, or in the database.
type Provider interface {
	SessionInit(sid string) (Session, error) // inits and returns a new session
	SessionRead(sid string) (Session, error) // creates or returns existing session represented by an sid
	SessionDestroy(sid string) error         // deletes session by sid
	SessionGC(maxLifeTime int64)             // deletes expired session variables according to maxLifeTime
}

// Session implements session operations: setValue, getValue, deleteValue, and getCurrentSessionID.
type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // get current sessionID
}

// <"memory"|"filesys"|"database"> : <Provider>
var providers = make(map[string]Provider)

// Register makes a session provider available by the provided key. If a Register is called twice with the same name or if the driver is nil,
// it panics.
func Register(name string, provider Provider) {
	// Handles nil provider
	if provider == nil {
		panic("session: Register provider is nil")
	}
	// Handles provider already registered -- WHY 2 RETURNS FROM provides[name]??
	if _, dup := providers[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	// Registers provider
	providers[name] = provider
}

// NewManager returns a new session manager.
func NewManager(providerName, cookieName string, maxlifetime int64) (*Manager, error) {
	// provideName indicates in what way sessions will be persisted
	provider, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", providerName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

// Creates a unique sessionID.
func (manager *Manager) sessionID() string {
	// Makes byte-slice of length: 32
	b := make([]byte, 32)
	// Copies 32 random bytes into b
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	// Converts to string
	return base64.URLEncoding.EncodeToString(b)
}

// SessionStart checks the existence of any sessions related to the current user, creating a new session if none is found.
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	// Protects session during start
	manager.lock.Lock()
	defer manager.lock.Unlock()
	// Retrieves cookie by name
	cookie, err := r.Cookie(manager.cookieName)
	// In case of error or empty cookie,
	if err != nil || cookie.Value == "" {
		// Gets a new sessionID
		sid := manager.sessionID()
		// Creates a new session using sid
		session, _ = manager.provider.SessionInit(sid)
		// Responds to client w/ cookie containing sid
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
	} else { // If the sid is empty,
		// Gets sid from cookie
		sid, err := url.QueryUnescape(cookie.Value)
		if err != nil {
			http.Error(w, "error unescaping cookie-value", http.StatusInternalServerError)
		}
		// Returns an existing session
		session, err = manager.provider.SessionRead(sid)
		if err != nil {
			http.Error(w, "error reading session", http.StatusInternalServerError)
		}
	}
	return
}

// SessionDestroy deletes session by ID.
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	// Handles empty cookie
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	// Protects session during destroy
	manager.lock.Lock()
	defer manager.lock.Unlock()
	// Deletes the sid
	manager.provider.SessionDestroy(cookie.Value)
	// Responds to client w/ cookie, sans sid
	cookie = &http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: time.Now(), MaxAge: -1}
	http.SetCookie(w, cookie)
}

// // ???
// func Count(w http.ResponseWriter, r *http.Request) {
// 	// Starts a new session
// 	sess := MegaManager.SessionStart(w, r)
// 	// Gets session time-stamp
// 	createtime := sess.Get("createtime")
// 	if createtime == nil {
// 		// Sets time-stamp if nil
// 		sess.Set("createtime", time.Now().Unix())
// 	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
// 		// GCs and starts new session if expired
// 		MegaManager.SessionDestroy(w, r)
// 		sess = MegaManager.SessionStart(w, r)
// 	}
// 	// Gets count
// 	ct := sess.Get("countnum")
// 	if ct == nil {
// 		// Inits to "1" if nil
// 		sess.Set("countnum", 1)
// 	} else {
// 		// Else, increments count
// 		sess.Set("countnum", (ct.(int) + 1))
// 	}
// 	// Get count-template
// 	t, _ := template.ParseFiles("count.gtpl")
// 	// Injects value associated w/ "countnum" into the template and renders
// 	w.Header().Set("Content-Type", "text/html")
// 	t.Execute(w, sess.Get("countnum"))
// }

// GC garbage collects old session.
func (manager *Manager) GC() {
	// Protects session during gc
	manager.lock.Lock()
	defer manager.lock.Unlock()
	// GCs sessions that have exceeded maxlifetime
	manager.provider.SessionGC(manager.maxlifetime)
	// Resets GC cycle
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}

// Example that uses sessions for a login operation.
func login(w http.ResponseWriter, r *http.Request) {
	sess := MegaManager.SessionStart(w, r)
	r.ParseForm()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, sess.Get("username"))
	} else {
		sess.Set("username", r.Form["username"])
		http.Redirect(w, r, "/", 302)
	}
}
