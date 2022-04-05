package filebased

import "server/storage"

func (fdb *FilebasedDB) GetCategories() []storage.Category {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	return fdb.Categories
}

func (fdb *FilebasedDB) FindCategory(cid int) bool {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	for _, category := range fdb.Categories {
		if cid == category.ID {
			return true
		}
	}
	return false
}
