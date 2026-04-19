package user

import (
	"bytes"
	"crypto/sha512"
	"hash"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type User struct {
	Login    string
	Password []byte
}

type Repository struct {
	mu    sync.RWMutex
	hash  hash.Hash
	users map[string]User
}

func NewRepository() *Repository {
	return &Repository{
		users: make(map[string]User),
		hash:  sha512.New(),
	}
}

func (r *Repository) Create(user User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.Login] = user

	if err := os.MkdirAll(filepath.Join("users_files", user.Login), 0o755); err != nil {
		log.Printf("erro ao criar diretorio do usuario %q: %v", user.Login, err)
	}
}

func (r *Repository) Retrieve(login string) (User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[login]
	return u, exists
}

func (r *Repository) Delete(login string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.users, login)
}

func (r *Repository) HandleAuth(login string, password []byte) bool {
	user, exists := r.Retrieve(login)
	if !exists {
		return false
	}
	if !bytes.Equal(user.Password, password) {
		return false
	}
	return true
}
