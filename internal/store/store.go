package store

import (
	"errors"
	"sync"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(id int) string {
	if id == 0 {
		return "0"
	}
	var str []byte
	for id > 0 {
		r := id % 62
		str = append(str, alphabet[r])
		id /= 62
	}
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	return string(str)
}

type User struct {
	ID           int
	Login        string
	PasswordHash string `json:"-"`
}

type Store struct {
	mu         sync.Mutex
	urls       map[string]string
	nextID     int
	users      map[string]User
	nextUserID int
}

func (s *Store) SaveUser(login, passwordHash string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[login]; ok {
		return User{}, errors.New("User already exists")
	}
	u := User{Login: login, PasswordHash: passwordHash}
	id := s.nextUserID
	u.ID = id + 1
	s.users[u.Login] = u
	s.nextID++
	return u, nil
}
func (s *Store) GetUser(login string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[login]
	if !ok {
		return User{}, errors.New("User not found")
	}
	return u, nil
}

func New() *Store {
	return &Store{
		urls:   make(map[string]string),
		nextID: 1,
		users:  make(map[string]User),
	}
}
func (s *Store) Save(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	code := encodeBase62(s.nextID)
	s.urls[code] = url
	s.nextID++
	return code
}

func (s *Store) Get(code string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, ok := s.urls[code]
	if !ok {
		return "", false
	}
	return url, true
}
