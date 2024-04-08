package data

import (
	"time"

	"github.com/skraio/banner-service/internal/validator"
)

type Banner struct {
	BannerID  int64     `json:"id"`
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

func ValidateBanner(v *validator.Validator, banner *Banner) {
    v.Check(banner.BannerID >= 0, "banner_id", "must be non-negative")

    v.Check(banner.TagIDs != nil, "tag_id", "must be provided")
    v.Check(len(banner.TagIDs) >= 1, "tag_ids", "must contain at least 1 tag")
    v.Check(len(banner.TagIDs) <= 1000, "tag_ids", "must not contain more than 1000 tags")
    v.Check(validator.Unique(banner.TagIDs), "tag_ids", "must not contain duplicate values")
    v.Check(validator.NonNegative(banner.TagIDs), "tag_ids", "must be non-negative")

    v.Check(banner.FeatureID >= 0, "feature_id", "must be non-negative")

	v.Check(banner.Content.Title != "", "content.title", "must be provided")
	v.Check(len(banner.Content.Title) <= 50, "title", "must not be more than 50 bytes long")

	v.Check(banner.Content.Text != "", "content.text", "must be provided")
	v.Check(len(banner.Content.Text) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(banner.Content.URL != "", "content.url", "must be provided")
	// add regex check for url
}
