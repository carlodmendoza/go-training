package filebased

import (
	"server/storage"
	"time"
)

func (fdb *FilebasedDB) CreateSession(username, token string) error {
	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	newSession := storage.Session{
		Token:     token,
		Timestamp: time.Now().Unix(),
		Username:  username,
	}
	fdb.Sessions[token] = newSession

	user := fdb.Users[username]
	if user.SessionToken != "" {
		delete(fdb.Sessions, user.SessionToken)
	}
	user.SessionToken = token
	fdb.Users[username] = user

	return nil
}

func (fdb *FilebasedDB) FindSession(token string) (storage.Session, error) {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	return fdb.Sessions[token], nil
}
