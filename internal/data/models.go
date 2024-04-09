package data

import (
	"database/sql"
	"errors"
)

var (
    ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
    Banners BannerModel
}

func NewModels(db *sql.DB) Models {
    return Models{
        Banners: BannerModel{DB: db},
    }
}
