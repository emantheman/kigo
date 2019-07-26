package memory

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"kigo/session"
)

var pvder = &Provider{list: list.New()}

// SessionStore stores sessions.
type SessionStore struct {
	sid          string                      // unique session id
	timeAccessed time.Time                   // last access time
	value        map[interface{}]interface{} // session value stored inside
}

// Provider stores sessions in memory.
type Provider struct {
	lock     sync.Mutex               // lock
	sessions map[string]*list.Element // save in memory
	list     *list.List               // gc
}

// Set adds an entry to the session store and returns an error.
func (st *SessionStore) Set(key, value interface{}) error {
	// Associates key:value pair in the st
	st.value[key] = value
	// Updates the session
	if err := pvder.SessionUpdate(st.sid); err != nil {
		return err
	}
	return nil
}

// Get returns a value from the session store by key.
func (st *SessionStore) Get(key interface{}) interface{} {
	// Updates the session
	if err := pvder.SessionUpdate(st.sid); err != nil {
		return err
	}
	// Gets the value from the st
	if v, ok := st.value[key]; ok {
		return v
	}
	return fmt.Errorf("sessionstore contains no value associated with that key")
}

// Delete deletes an entry from the session store by key.
func (st *SessionStore) Delete(key interface{}) error {
	// Deletes entry from st
	delete(st.value, key)
	// Updates provider
	err := pvder.SessionUpdate(st.sid)
	return err
}

// SessionID returns the id of the session.
func (st *SessionStore) SessionID() string {
	return st.sid
}

// SessionInit initializes a session and stores it in memory.
func (pvder *Provider) SessionInit(sid string) (session.Session, error) {
	// Blocks provider
	pvder.lock.Lock()
	defer pvder.lock.Unlock()
	// Creates new session
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	// Adds new session to the back of the GC queue
	element := pvder.list.PushBack(newsess)
	// Associates the sid with the session
	pvder.sessions[sid] = element
	// Returns the new session
	return newsess, nil
}

// SessionRead returns a session by sid.
func (pvder *Provider) SessionRead(sid string) (session.Session, error) {
	// Gets session by sid
	if element, ok := pvder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	}
	// If no session is associated w/ that id, returns a new session
	sess, err := pvder.SessionInit(sid)
	return sess, err
}

// SessionDestroy destroys a session by id.
func (pvder *Provider) SessionDestroy(sid string) error {
	// Checks session exists
	if element, ok := pvder.sessions[sid]; ok {
		// Removes session from memory
		delete(pvder.sessions, sid)
		// Removes session from GC queue
		pvder.list.Remove(element)
	}
	return fmt.Errorf("could not destroy: no session with that id exists")
}

// SessionGC removes expired sessions from memory.
func (pvder *Provider) SessionGC(maxlifetime int64) {
	// Blocks provider
	pvder.lock.Lock()
	defer pvder.lock.Unlock()
	// Traverses list from back, removing expired sessions
	for {
		element := pvder.list.Back()
		// Handles empty list
		if element == nil {
			break
		}
		// Deletes expired sessions
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pvder.list.Remove(element)
			delete(pvder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

// SessionUpdate moves session to front of GC queue and updates accesstimestamp.
func (pvder *Provider) SessionUpdate(sid string) error {
	// Blocks provider
	pvder.lock.Lock()
	defer pvder.lock.Unlock()
	// Gets session by sid
	if element, ok := pvder.sessions[sid]; ok {
		// Updates time-accessed
		element.Value.(*SessionStore).timeAccessed = time.Now()
		// Moves element to front of list
		pvder.list.MoveToFront(element)
		return nil
	}
	return fmt.Errorf("could not update session: no session by that sid")
}

// IsAuthenticated checks whether a session is associated with this sid.
func (pvder *Provider) IsAuthenticated(sid string) bool {
	_, ok := pvder.sessions[sid]
	if ok {
		return true
	}
	return false
}

func init() {
	// Initializes provider with an empty sessions list
	pvder.sessions = make(map[string]*list.Element, 0)
	// Registers the provider
	session.Register("memory", pvder)
}
