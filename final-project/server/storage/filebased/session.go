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

	user := fdb.Users[username]
	if user.SessionToken != "" {
		delete(fdb.Sessions, user.SessionToken)
	}
	user.SessionToken = token

	newSession := storage.Session{
		Token:     token,
		Timestamp: time.Now().Unix(),
		UserID:    user.ID,
	}

	fdb.Sessions[token] = newSession

	return nil
}

func (fdb *FilebasedDB) FindSession(token string) (int, error) {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	session, exists := fdb.Sessions[token]
	if exists && session.Token == token {
		return session.UserID, nil
	}
	return 0, nil
}
