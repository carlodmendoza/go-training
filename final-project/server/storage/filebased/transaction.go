package filebased

import (
	"server/storage"
	"sort"
)

func (fdb *FilebasedDB) CreateTransaction(tr storage.Transaction) error {
	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	fdb.NextTransactionID++
	newTransaction := storage.Transaction{
		ID:         fdb.NextTransactionID,
		Amount:     tr.Amount,
		Date:       tr.Date,
		Notes:      tr.Notes,
		Username:   tr.Username,
		CategoryID: tr.CategoryID,
	}
	fdb.Transactions[fdb.NextTransactionID] = newTransaction

	user := fdb.Users[tr.Username]
	user.Transactions[fdb.NextTransactionID] = struct{}{}
	fdb.Users[tr.Username] = user

	return nil
}

func (fdb *FilebasedDB) GetTransactions(username string) ([]storage.Transaction, error) {
	transactions := []storage.Transaction{}

	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	for k := range fdb.Users[username].Transactions {
		transactions = append(transactions, fdb.Transactions[k])
	}
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].ID < transactions[j].ID
	})

	return transactions, nil
}

func (fdb *FilebasedDB) UpdateTransaction(tr storage.Transaction) error {
	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	fdb.Transactions[tr.ID] = tr

	return nil
}

func (fdb *FilebasedDB) DeleteTransactions(username string) (bool, error) {
	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	user := fdb.Users[username]
	if len(user.Transactions) == 0 {
		return false, nil
	}

	for k := range user.Transactions {
		delete(fdb.Transactions, k)
	}
	user.Transactions = map[int]struct{}{}
	fdb.Users[username] = user

	return true, nil
}

func (fdb *FilebasedDB) DeleteTransaction(username string, tid int) error {
	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	delete(fdb.Transactions, tid)
	delete(fdb.Users[username].Transactions, tid)

	return nil
}

func (fdb *FilebasedDB) FindTransaction(username string, tid int) (storage.Transaction, bool, error) {
	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	transaction, exists := fdb.Transactions[tid]
	if exists && transaction.Username == username {
		return transaction, true, nil
	}
	return storage.Transaction{}, false, nil
}
