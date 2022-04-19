package filebased

import (
	"server/storage"
	"sort"
)

func (fdb *FilebasedDB) GetCategories() ([]storage.Category, error) {
	categories := []storage.Category{}

	fdb.CategoryMux.RLock()
	defer fdb.CategoryMux.RUnlock()

	for _, v := range fdb.Categories {
		categories = append(categories, v)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].ID < categories[j].ID
	})

	return categories, nil
}

func (fdb *FilebasedDB) CategoryExists(cid int) (bool, error) {
	fdb.CategoryMux.RLock()
	defer fdb.CategoryMux.RUnlock()

	_, exists := fdb.Categories[cid]
	return exists, nil
}
