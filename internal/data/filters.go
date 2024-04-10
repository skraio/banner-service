package data

import (
	"github.com/skraio/banner-service/internal/validator"
)

type ReadIntOptions struct {
	Required bool
	IsID     bool
}

type UserFilters struct {
	TagID           int
	FeatureID       int
	UseLastRevision bool
}

type AdminFilters struct {
	TagID     int
	FeatureID int
	Limit     int
	Offset    int
}

func ValidateUserFilters(v *validator.Validator, f UserFilters) {
	v.Check(f.TagID > 0, "tag_id", "must be greater than zero")
	v.Check(f.FeatureID > 0, "feature_id", "must be greater than zero")
}

func ValidateAdminFilters(v *validator.Validator, f AdminFilters) {
	if f.TagID != 0 {
		v.Check(f.TagID > 0, "tag_id", "must be greater than zero")
	}
	if f.FeatureID != 0 {
		v.Check(f.FeatureID > 0, "feature_id", "must be greater than zero")
	}
	v.Check(f.Limit > 0, "limit", "must be greater than zero")
	v.Check(f.Limit <= 100, "limit", "must be a maximum of 100")
	v.Check(f.Offset >= 0, "offset", "must be non-negative")
}

type Metadata struct {
	Offset       int `'json:"offset,omitempty"`
	Limit        int `'json:"limit,omitempty"`
	TotalRecords int `'json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, offset, limit int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		Offset:       offset,
		Limit:        limit,
		TotalRecords: totalRecords,
	}
}
