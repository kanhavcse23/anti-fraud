package account_route_v1

import (
	controllerV1Package "anti-fraud/account-service/controllers/v1"

	middlewareHandlerPackageV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
)

type AccountRoutes struct {
	controller        *controllerV1Package.AccountController
	muxRouter         *mux.Router
	middlewareHandler *middlewareHandlerPackageV1.MiddlewareHandler
}

func NewAccountRoutes(controller *controllerV1Package.AccountController, router *mux.Router, middlewareHandler *middlewareHandlerPackageV1.MiddlewareHandler) *AccountRoutes {
	return &AccountRoutes{controller: controller, muxRouter: router, middlewareHandler: middlewareHandler}

}

func (routes *AccountRoutes) Init() {
	handlerFunc := routes.middlewareHandler.MiddlewareHandlerFunc

	routes.muxRouter.HandleFunc("/accounts", handlerFunc(routes.controller.CreateAccount)).Methods("POST")
	routes.muxRouter.HandleFunc("/accounts/{accountId}", handlerFunc(routes.controller.GetAccountDetails)).Methods("GET")
}
