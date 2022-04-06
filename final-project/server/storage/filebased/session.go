package filebased

import (
	"server/auth"
	"server/storage"
	"time"
)

func (fdb *FilebasedDB) CreateSession(uid int) storage.Session {
	newSession := storage.Session{
		Token:     auth.GenerateSessionToken(),
		Timestamp: time.Now().Unix(),
		UserID:    uid,
	}

	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	sessionExists := false
	for i, session := range fdb.Sessions {
		if uid == session.UserID {
			fdb.Sessions[i] = newSession
			sessionExists = true
			break
		}
	}

	if !sessionExists {
		fdb.Sessions = append(fdb.Sessions, newSession)
	}

	return newSession
}

func (fdb *FilebasedDB) FindSession(token string) int {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	for _, session := range fdb.Sessions {
		if token == session.Token {
			return session.UserID
		}
	}
	return 0
}
