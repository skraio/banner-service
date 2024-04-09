package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/skraio/banner-service/internal/validator"
)

type Banner struct {
	BannerID  int64     `json:"banner_id"`
	FeatureID int64     `json:"feature_id"`
	TagIDs    []int64   `json:"tag_ids"`
	Content   Content   `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

// func (c *Content) Scan(src interface{}) error {
//     if src == nil {
//         return nil
//     }
//     return json.Unmarshal(src.([]byte), c)
// }

func ValidateBanner(v *validator.Validator, banner *Banner) {
	v.Check(banner.BannerID >= 0, "banner_id", "must be positive")

	v.Check(banner.TagIDs != nil, "tag_id", "must be provided")
	v.Check(len(banner.TagIDs) >= 1, "tag_ids", "must contain at least 1 tag")
	v.Check(len(banner.TagIDs) <= 1000, "tag_ids", "must not contain more than 1000 tags")
	v.Check(validator.Unique(banner.TagIDs), "tag_ids", "must not contain duplicate values")
	v.Check(validator.Positive(banner.TagIDs), "tag_ids", "must be positive")

	v.Check(banner.FeatureID > 0, "feature_id", "must be positive")

	v.Check(banner.Content.Title != "", "content.title", "must be provided")
	v.Check(len(banner.Content.Title) <= 50, "title", "must not be more than 50 bytes long")

	v.Check(banner.Content.Text != "", "content.text", "must be provided")
	v.Check(len(banner.Content.Text) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(banner.Content.URL != "", "content.url", "must be provided")
	// add regex check for url
}

type BannerModel struct {
	DB *sql.DB
}

func (b BannerModel) Insert(banner *Banner) error {
	contentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO banners (tag_ids, feature_id, content, is_active)
        VALUES ($1, $2, $3, $4)
        RETURNING banner_id, created_at, updated_at`

	args := []interface{}{pq.Array(banner.TagIDs), banner.FeatureID, string(contentJSON), banner.IsActive}

	return b.DB.QueryRow(query, args...).Scan(&banner.BannerID, &banner.CreatedAt, &banner.UpdatedAt)
}

func (b BannerModel) Get(tagID, featureID int64, useLastRevision bool) (*Banner, error) {
	if tagID < 1 || featureID < 1 {
		return nil, ErrRecordNotFound
	}

	// if useLastRevision {
	// }

	query := `
        SELECT banner_id, content, created_at, updated_at, is_active
        FROM banners
        WHERE is_active = true AND $1 = ANY(tag_ids) AND feature_id = $2`

	var banner Banner
	var contentJSON []byte

	err := b.DB.QueryRow(query, tagID, featureID).Scan(
		&banner.BannerID,
		&contentJSON,
		&banner.CreatedAt,
		&banner.UpdatedAt,
		&banner.IsActive,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	err = json.Unmarshal(contentJSON, &banner.Content)
	if err != nil {
		return nil, err
	}

	return &banner, nil
}

func (b BannerModel) GetByID(banner_id int64) (*Banner, error) {
	if banner_id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT banner_id, content, created_at, updated_at, is_active, tag_ids, feature_id
        FROM banners
        WHERE banner_id = $1`

	var banner Banner
	var contentJSON []byte

	err := b.DB.QueryRow(query, banner_id).Scan(
		&banner.BannerID,
		&contentJSON,
		&banner.CreatedAt,
		&banner.UpdatedAt,
		&banner.IsActive,
		pq.Array(&banner.TagIDs),
		&banner.FeatureID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	err = json.Unmarshal(contentJSON, &banner.Content)
	if err != nil {
		return nil, err
	}

	return &banner, nil
}

func (b BannerModel) Gets() ([]*Banner, error) {
	return nil, nil
}

func (b BannerModel) Update(banner *Banner) error {
    query := `
        UPDATE banners
        SET
            tag_ids = $1,
            feature_id = $2,
            content = $3,
            is_active = $4,
            updated_at = NOW()
        WHERE banner_id = $5 AND updated_at = $6
        RETURNING updated_at`

	contentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		return err
	}

	args := []any{
		pq.Array(banner.TagIDs),
		banner.FeatureID,
        contentJSON,
		banner.IsActive,
		banner.BannerID,
        banner.UpdatedAt,
	}

    err = b.DB.QueryRow(query, args...).Scan(&banner.UpdatedAt)
    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return ErrEditConflict
        default:
            return err
        }
    }

    return nil
}

func (b BannerModel) Delete(banner_id int64) error {
	if banner_id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM banners
        WHERE banner_id = $1`

	result, err := b.DB.Exec(query, banner_id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
