package filebased

import "server/storage"

func (fdb *FilebasedDB) CreateUser(username, password string) error {
	fdb.UserMux.Lock()
	defer func() {
		fdb.UserMux.Unlock()
		appendData(fdb)
	}()

	fdb.NextUserID++
	newUser := storage.User{
		ID:           fdb.NextUserID,
		Name:         username,
		Password:     password,
		Transactions: map[int]struct{}{},
	}
	fdb.Users[username] = newUser

	return nil
}

func (fdb *FilebasedDB) UserExists(username string) (bool, error) {
	fdb.UserMux.RLock()
	defer fdb.UserMux.RUnlock()

	_, exists := fdb.Users[username]
	return exists, nil
}

func (fdb *FilebasedDB) AuthenticateUser(username, password string) (bool, error) {
	fdb.UserMux.RLock()
	defer fdb.UserMux.RUnlock()

	user, exists := fdb.Users[username]
	if exists && user.Password == password {
		return true, nil
	}
	return false, nil
}
