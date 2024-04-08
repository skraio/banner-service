package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/skraio/banner-service/internal/data"
	"github.com/skraio/banner-service/internal/validator"
)

// GET /v1/user_banner
func (app *application) showBannerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "show banner for user")

	// dummy
	banner := data.Banner{
		BannerID:  123,
		FeatureID: 111,
		TagIDs:    []int64{7, 8, 9},
		Content: data.Content{
			Title: "Cakes",
			Text:  "Homemade cakes for birthdays.",
			URL:   "https://example.com/",
		},
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"banner": banner}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GET /v1/banner
func (app *application) showAllBannersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "show all banners filtered by feature and/or tag")
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

	fmt.Fprintf(w, "%+v\n", input)
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
