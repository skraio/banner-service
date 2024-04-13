package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/skraio/banner-service/internal/data"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/user_banner", app.requireRole(app.showBannerHandler, data.RoleUser, data.RoleAdmin))
	router.HandlerFunc(http.MethodGet, "/banner", app.requireRole(app.listFilteredBannersHandler, data.RoleAdmin))
	router.HandlerFunc(http.MethodPost, "/banner", app.requireRole(app.createBannerHandler, data.RoleAdmin))
	router.HandlerFunc(http.MethodPatch, "/banner/:id", app.requireRole(app.updateBannerHandler, data.RoleAdmin))
	router.HandlerFunc(http.MethodDelete, "/banner/:id", app.requireRole(app.deleteBannerHandler, data.RoleAdmin))

	router.HandlerFunc(http.MethodPost, "/user", app.createUserHandler)
	router.HandlerFunc(http.MethodPost, "/token", app.createTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
