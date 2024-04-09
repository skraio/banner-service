package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/skraio/banner-service/internal/data"
	"github.com/skraio/banner-service/internal/validator"
)

// GET /v1/user_banner
// curl -i -H "token: user_token" "localhost:4000/v1/user_banner?tag_id=111&feature_id=777"
func (app *application) showBannerHandler(w http.ResponseWriter, r *http.Request) {
	tagID, err := app.readTagIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	featureID, err := app.readFeatureIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	useLastRevision, err := app.readUseLastRevisionParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// userToken

	banner, err := app.models.Banners.Get(tagID, featureID, useLastRevision)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
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

// GET /v1/banner
func (app *application) showAllBannersHandler(w http.ResponseWriter, r *http.Request) {
}

// POST /v1/banner
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

// PATCH /v1/banner/{id}
func (app *application) updateBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
	}

	fmt.Fprintf(w, "update the banner %d content\n", id)
}

// DELETE /v1/banner/{id}
func (app *application) deleteBannerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
	}

	fmt.Fprintf(w, "delete a banner %d\n", id)
}
