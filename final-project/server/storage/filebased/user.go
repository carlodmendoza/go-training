package filebased

import "server/storage"

func (fdb *FilebasedDB) CreateUser(username, password string) error {
	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	fdb.NextUserID++
	newUser := storage.User{
		ID:           fdb.NextUserID,
		Name:         username,
		Password:     password,
		Transactions: map[int]struct{}{},
	}
	fdb.Users[username] = &newUser

	return nil
}

func (fdb *FilebasedDB) FindUser(username string) (bool, error) {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	_, exists := fdb.Users[username]
	return exists, nil
}

func (fdb *FilebasedDB) AuthenticateUser(username, password string) (bool, error) {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	user, exists := fdb.Users[username]
	if exists && user.Password == password {
		return true, nil
	}
	return false, nil
}
