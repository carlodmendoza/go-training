package filebased

import "server/storage"

func (fdb *FilebasedDB) CreateUser(username, password string) {
	fdb.Mu.Lock()

	fdb.NextUserID++
	newUser := storage.User{
		ID:       fdb.NextUserID,
		Name:     username,
		Password: password,
	}
	fdb.Users = append(fdb.Users, newUser)

	fdb.Mu.Unlock()

	updateDatabase(fdb)
}

func (fdb *FilebasedDB) FindUser(username string) bool {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	for _, user := range fdb.Users {
		if username == user.Name {
			return true
		}
	}
	return false
}

func (fdb *FilebasedDB) AuthenticateUser(username, password string) (int, bool) {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	for _, user := range fdb.Users {
		if username == user.Name && password == user.Password {
			return user.ID, true
		}
	}
	return 0, false
}
