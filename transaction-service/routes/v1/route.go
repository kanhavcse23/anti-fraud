package transaction_route_v1

import (
	controllerV1Package "anti-fraud/transaction-service/controllers/v1"

	middlewareHandlerPackageV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
)

type TransactionRoutes struct {
	controller        *controllerV1Package.TransactionController
	muxRouter         *mux.Router
	middlewareHandler *middlewareHandlerPackageV1.MiddlewareHandler
}

func NewTransactionRoutes(controller *controllerV1Package.TransactionController, router *mux.Router, middlewareHandler *middlewareHandlerPackageV1.MiddlewareHandler) *TransactionRoutes {
	return &TransactionRoutes{controller: controller, muxRouter: router, middlewareHandler: middlewareHandler}

}

func (routes *TransactionRoutes) Init() {
	handlerFunc := routes.middlewareHandler.MiddlewareHandlerFunc

	routes.muxRouter.HandleFunc("/transactions/v1", handlerFunc(routes.controller.CreateTransaction)).Methods("POST")
}
