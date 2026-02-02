package storage

type GofferStorage struct {
	DBStorage DBStorage
}

func NewGofferStorage(dataBaseDsn string) *GofferStorage {
	return &GofferStorage{
		DBStorage: newDatabase(dataBaseDsn),
	}
}