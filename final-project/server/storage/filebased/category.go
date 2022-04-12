package filebased

import (
	"server/storage"
	"sort"
)

func (fdb *FilebasedDB) GetCategories() ([]storage.Category, error) {
	categories := []storage.Category{}

	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	for _, v := range fdb.Categories {
		categories = append(categories, v)
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].ID < categories[j].ID
	})

	return categories, nil
}

func (fdb *FilebasedDB) FindCategory(cid int) (bool, error) {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	_, exists := fdb.Categories[cid]
	return exists, nil
}
