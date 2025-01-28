package account_route_v1

import (
	controllerV1Package "anti-fraud/account-service/controllers/v1"

	middlewareHandlerPackageV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
)

type AccountRoutes struct {
	controller        controllerV1Package.IAccountController
	muxRouter         *mux.Router
	middlewareHandler *middlewareHandlerPackageV1.MiddlewareHandler
}

// NewAccountRoutes create and return an instance of AccountRoutes.
func NewAccountRoutes(controller controllerV1Package.IAccountController, router *mux.Router, middlewareHandler *middlewareHandlerPackageV1.MiddlewareHandler) *AccountRoutes {
	return &AccountRoutes{controller: controller, muxRouter: router, middlewareHandler: middlewareHandler}

}

// Init register route for account-service.
func (routes *AccountRoutes) Init() {
	handlerFunc := routes.middlewareHandler.MiddlewareHandlerFunc

	routes.muxRouter.HandleFunc("/accounts/v1", handlerFunc(routes.controller.CreateAccount)).Methods("POST")
	routes.muxRouter.HandleFunc("/accounts/v1/{accountId}", handlerFunc(routes.controller.GetAccountDetails)).Methods("GET")
}
