package main

import (
    "net/http"

    "github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
    router := httprouter.New()

    router.NotFound = http.HandlerFunc(app.notFoundResponse)
    router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

    router.HandlerFunc(http.MethodGet, "/v1/user_banner", app.showBannerHandler)
    router.HandlerFunc(http.MethodGet, "/v1/banner", app.showAllBannersHandler)
    router.HandlerFunc(http.MethodPost, "/v1/banner", app.createBannerHandler)
    router.HandlerFunc(http.MethodPatch, "/v1/banner/:id", app.updateBannerHandler)
    router.HandlerFunc(http.MethodDelete, "/v1/banner/:id", app.deleteBannerHandler)

    return app.recoverPanic(router)
}
