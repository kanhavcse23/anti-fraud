package account_route_v1

import (
	controllerV1Path "anti-fraud/account-service/controllers/v1"

	middlewareHandlerPathV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
)

type AccountRoutes struct {
	controller        *controllerV1Path.AccountController
	muxRouter         *mux.Router
	middlewareHandler *middlewareHandlerPathV1.MiddlewareHandler
}

func NewAccountRoutes(controller *controllerV1Path.AccountController, router *mux.Router, middlewareHandler *middlewareHandlerPathV1.MiddlewareHandler) *AccountRoutes {
	return &AccountRoutes{controller: controller, muxRouter: router, middlewareHandler: middlewareHandler}

}

func (routes *AccountRoutes) Init() {

	routes.muxRouter.HandleFunc("/accounts", routes.middlewareHandler.MiddlewareHandlerFunc(routes.controller.CreateAccount)).
		Methods("POST")
	routes.muxRouter.HandleFunc("/accounts/{accountId}", routes.middlewareHandler.MiddlewareHandlerFunc(routes.controller.GetAccountDetails)).
		Methods("GET")
}
