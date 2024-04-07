package main

import (
    "fmt"
    "net/http"
)


// GET /v1/user_banner
func (app *application) showBannerHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "show banner for user")
}

// GET /v1/banner
func (app *application) showAllBannersHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "show all banners filtered by feature and/or tag")
}

// POST /v1/banner 
func (app *application) createBannerHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "create a new banner")
}

// PATCH /v1/banner/{id} 
func (app *application) updateBannerHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "update the banner %d content\n", id)
}

// DELETE /v1/banner/{id}
func (app *application) deleteBannerHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "delete a banner %d\n", id)
}
