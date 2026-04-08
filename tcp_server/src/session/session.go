package session

import (
	"sync"
	"sync/atomic"
)

type Session struct {
	ID      string
	User    string
	RootDir string
	CurrDir string
}

type Repository struct {
	mu       sync.Mutex
	sessions atomic.Value
}

func NewRepository() *Repository {
	r := &Repository{}
	r.sessions.Store(make(map[string]Session))
	return r
}

func (r *Repository) Create(s Session) Session {
	r.mu.Lock()
	defer r.mu.Unlock()

	next := cloneSessions(r.snapshot())
	next[s.ID] = s
	r.sessions.Store(next)
	return s
}

func (r *Repository) Retrieve(id string) (Session, bool) {
	s, exists := r.snapshot()[id]
	return s, exists
}

func (r *Repository) Update(id string, currDir string) (Session, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	next := cloneSessions(r.snapshot())
	s, exists := next[id]
	if !exists {
		return Session{}, false
	}

	s.CurrDir = currDir
	next[id] = s
	r.sessions.Store(next)
	return s, true
}

func (r *Repository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	next := cloneSessions(r.snapshot())
	if _, exists := next[id]; !exists {
		return false
	}

	delete(next, id)
	r.sessions.Store(next)
	return true
}

func (r *Repository) List() []Session {
	sessions := r.snapshot()
	result := make([]Session, 0, len(sessions))
	for _, s := range sessions {
		result = append(result, s)
	}

	return result
}

func (r *Repository) snapshot() map[string]Session {
	return r.sessions.Load().(map[string]Session)
}

func cloneSessions(src map[string]Session) map[string]Session {
	dst := make(map[string]Session, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
