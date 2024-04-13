package main

import (
	"errors"
	"net/http"

	"github.com/skraio/banner-service/internal/data"
	"github.com/skraio/banner-service/internal/validator"
)

func (app *application) showBannerHandler(w http.ResponseWriter, r *http.Request) {

	var filters data.UserFilters

	qs := r.URL.Query()
	v := validator.New()

	filters.TagID = app.readInt(qs, "tag_id", 0, data.ReadIntOptions{Required: true, IsID: true}, v)
	filters.FeatureID = app.readInt(qs, "feature_id", 0, data.ReadIntOptions{Required: true, IsID: true}, v)
	filters.UseLastRevision = app.readBool(qs, "use_last_revision", false, v)

	if data.ValidateRequiredFilters(v, filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)

	banner, err := app.models.Banners.Get(filters, user.Role)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		case errors.Is(err, data.ErrForbiddenAccess):
			app.forbiddenAccessResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"content": banner.Content}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listFilteredBannersHandler(w http.ResponseWriter, r *http.Request) {
	var filters data.AdminFilters

	v := validator.New()

	qs := r.URL.Query()

	filters.FeatureID = app.readInt(qs, "feature_id", 0, data.ReadIntOptions{Required: false, IsID: true}, v)
	filters.TagID = app.readInt(qs, "tag_id", 0, data.ReadIntOptions{Required: false, IsID: true}, v)
	filters.Limit = app.readInt(qs, "limit", 10, data.ReadIntOptions{Required: false, IsID: false}, v)
	filters.Offset = app.readInt(qs, "offset", 0, data.ReadIntOptions{Required: false, IsID: false}, v)

	if data.ValidateNonRequiredFilters(v, filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	banners, metadata, err := app.models.Banners.GetAll(filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "banners": banners}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createBannerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TagIDs    []int64      `json:"tag_ids"`
		FeatureID int64        `json:"feature_id"`
		Content   data.Content `json:"content"`
		IsActive  bool         `json:"is_active"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	banner := &data.Banner{
		TagIDs:    input.TagIDs,
		FeatureID: input.FeatureID,
		Content:   input.Content,
		IsActive:  input.IsActive,
	}

	v := validator.New()

	if data.ValidateBanner(v, banner); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Banners.Insert(banner)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"banner_id": banner.BannerID}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
	}

	banner, err := app.models.Banners.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		TagIDs    []int64       `json:"tag_ids"`
		FeatureID *int64        `json:"feature_id"`
		Content   *data.Content `json:"content"`
		IsActive  *bool         `json:"is_active"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.TagIDs != nil {
		banner.TagIDs = input.TagIDs
	}

	if input.FeatureID != nil {
		banner.FeatureID = *input.FeatureID
	}

	if input.Content != nil {
		if input.Content.Title != "" {
			banner.Content.Title = input.Content.Title
		}
		if input.Content.Text != "" {
			banner.Content.Text = input.Content.Text
		}
		if input.Content.URL != "" {
			banner.Content.URL = input.Content.URL
		}
	}

	if input.IsActive != nil {
		banner.IsActive = *input.IsActive
	}

	v := validator.New()

	if data.ValidateBanner(v, banner); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Banners.Update(banner)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"banner": banner}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
	}

	err = app.models.Banners.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
