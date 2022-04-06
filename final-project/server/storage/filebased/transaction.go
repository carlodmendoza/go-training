package filebased

import "server/storage"

func (fdb *FilebasedDB) CreateTransaction(tr storage.Transaction) {
	fdb.Mu.Lock()

	fdb.NextTransactionID++
	newTransaction := storage.Transaction{
		ID:         fdb.NextTransactionID,
		Amount:     tr.Amount,
		Date:       tr.Date,
		Notes:      tr.Notes,
		UserID:     tr.UserID,
		CategoryID: tr.CategoryID,
	}
	fdb.Transactions = append(fdb.Transactions, newTransaction)

	fdb.Mu.Unlock()

	updateDatabase(fdb)
}

func (fdb *FilebasedDB) GetTransactions(uid int) []storage.Transaction {
	transactions := []storage.Transaction{}

	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	for _, transaction := range fdb.Transactions {
		if uid == transaction.UserID {
			transactions = append(transactions, transaction)
		}
	}

	return transactions
}

func (fdb *FilebasedDB) UpdateTransaction(tr storage.Transaction) {
	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	for i, transaction := range fdb.Transactions {
		if tr.ID == transaction.ID {
			fdb.Transactions[i] = tr
			break
		}
	}
}

func (fdb *FilebasedDB) DeleteTransactions(uid int) bool {
	if len(fdb.GetTransactions(uid)) == 0 {
		return false
	}

	transactions := []storage.Transaction{}

	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	for _, transaction := range fdb.Transactions {
		if uid != transaction.UserID {
			transactions = append(transactions, transaction)
		}
	}

	fdb.Transactions = transactions
	return true
}

func (fdb *FilebasedDB) DeleteTransaction(uid, tid int) {
	fdb.Mu.Lock()
	defer func() {
		fdb.Mu.Unlock()
		updateDatabase(fdb)
	}()

	for i, transaction := range fdb.Transactions {
		if tid == transaction.ID {
			fdb.Transactions = append(fdb.Transactions[:i], fdb.Transactions[i+1:]...)
			break
		}
	}
}

func (fdb *FilebasedDB) FindTransaction(uid, tid int) (storage.Transaction, bool) {
	transactions := fdb.GetTransactions(uid)

	fdb.Mu.Lock()
	defer fdb.Mu.Unlock()

	for _, transaction := range transactions {
		if tid == transaction.ID {
			return transaction, true
		}
	}
	return storage.Transaction{}, false
}
