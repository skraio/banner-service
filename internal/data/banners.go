package data

import (
	"context"
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

func (c *Content) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	return json.Unmarshal(src.([]byte), c)
}

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
	v.Check(validator.Matches(banner.Content.URL, validator.UrlRX), "content.url", "must be a valid URL")
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).Scan(&banner.BannerID, &banner.CreatedAt, &banner.UpdatedAt)
}

func (b BannerModel) Get(filters UserFilters, userRole Role) (*Banner, error) {
	if filters.TagID < 1 || filters.FeatureID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT banner_id, content, created_at, updated_at, is_active
        FROM banners
        WHERE $1 = ANY(tag_ids) AND feature_id = $2`

	var banner Banner

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, filters.TagID, filters.FeatureID).Scan(
		&banner.BannerID,
		&banner.Content,
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

	if !banner.IsActive && userRole == RoleUser {
		return nil, ErrForbiddenAccess
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, banner_id).Scan(
		&banner.BannerID,
		&banner.Content,
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

	return &banner, nil
}

func (b BannerModel) GetAll(filters AdminFilters) ([]*Banner, Metadata, error) {
	query := `
        SELECT count(*) OVER(), banner_id, tag_ids, feature_id, content, is_active, created_at, updated_at
        FROM banners
        WHERE (feature_id = $1 OR $1 = 0)
            AND ($2 = ANY(tag_ids) OR $2 = 0)
        LIMIT $3 OFFSET $4`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	args := []any{filters.FeatureID, filters.TagID, filters.Limit, filters.Offset}

	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	banners := []*Banner{}

	for rows.Next() {
		var banner Banner

		err := rows.Scan(
			&totalRecords,
			&banner.BannerID,
			pq.Array(&banner.TagIDs),
			&banner.FeatureID,
			&banner.Content,
			&banner.IsActive,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		banners = append(banners, &banner)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Offset, filters.Limit)

	return banners, metadata, nil
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = b.DB.QueryRowContext(ctx, query, args...).Scan(&banner.UpdatedAt)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, banner_id)
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
