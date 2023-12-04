package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/drills", app.requirePermission("drills:read", app.listDrillsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/drills", app.requirePermission("drills:write", app.createDrillHandler))
	router.HandlerFunc(http.MethodGet, "/v1/drills/:id", app.requirePermission("drills:read", app.showDrillHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/drills/:id", app.requirePermission("drills:write", app.updateDrillHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/drills/:id", app.requirePermission("drills:write", app.deleteDrillHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}