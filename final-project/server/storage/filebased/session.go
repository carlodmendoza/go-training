package filebased

import (
	"time"

	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

func (fdb *FilebasedDB) CreateSession(username, token string) error {
	fdb.SessionMux.Lock()

	newSession := storage.Session{
		Token:     token,
		Timestamp: time.Now().Unix(),
		Username:  username,
	}
	fdb.Sessions[token] = newSession

	fdb.UserMux.Lock()
	user := fdb.Users[username]
	if user.SessionToken != "" {
		delete(fdb.Sessions, user.SessionToken)
	}
	user.SessionToken = token
	fdb.Users[username] = user
	fdb.UserMux.Unlock()

	fdb.SessionMux.Unlock()

	err := appendData(fdb)
	if err != nil {
		return err
	}

	return nil
}

func (fdb *FilebasedDB) FindSession(token string) (storage.Session, error) {
	fdb.SessionMux.RLock()
	defer fdb.SessionMux.RUnlock()

	return fdb.Sessions[token], nil
}
