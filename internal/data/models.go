package data

import (
	"database/sql"
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
	ErrForbiddenAccess = errors.New("user does not have access")
	ErrCacheNotFound   = errors.New("cache not found")
)

type Models struct {
	Banners BannerModel
	Users   UserModel
	Tokens  TokenModel
}

func NewModels(db *sql.DB, rd *redis.Client) Models {
	return Models{
		Banners: BannerModel{DB: db, RD: rd},
		Users:   UserModel{DB: db},
		Tokens:  TokenModel{DB: db},
	}
}
