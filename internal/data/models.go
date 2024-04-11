package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
	ErrForbiddenAccess = errors.New("user does not have access")
)

type Models struct {
	Banners BannerModel
	Users   UserModel
	Tokens  TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Banners: BannerModel{DB: db},
		Users:   UserModel{DB: db},
		Tokens:  TokenModel{DB: db},
	}
}
