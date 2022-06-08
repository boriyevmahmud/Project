package storage

import (
	"github.com/mahmud3253/Project/User_Service/storage/postgres"
	"github.com/mahmud3253/Project/User_Service/storage/repo"

	"github.com/jmoiron/sqlx"
)

//IStorage ...
type IStorage interface {
	User() repo.UserStorageI
}

type storagePg struct {
	db       *sqlx.DB
	userRepo repo.UserStorageI
}

//NewStoragePg ...
func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		userRepo: postgres.NewUserRepo(db),
	}
}

func (s storagePg) User() repo.UserStorageI {
	return s.userRepo
}
