package transaction_route_v1

import (
	controllerV1Path "anti-fraud/transaction-service/controllers/v1"

	middlewareHandlerPathV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
)

type TransactionRoutes struct {
	controller        *controllerV1Path.TransactionController
	muxRouter         *mux.Router
	middlewareHandler *middlewareHandlerPathV1.MiddlewareHandler
}

func NewTransactionRoutes(controller *controllerV1Path.TransactionController, router *mux.Router, middlewareHandler *middlewareHandlerPathV1.MiddlewareHandler) *TransactionRoutes {
	return &TransactionRoutes{controller: controller, muxRouter: router, middlewareHandler: middlewareHandler}

}

func (routes *TransactionRoutes) Init() {
	handlerFunc := routes.middlewareHandler.MiddlewareHandlerFunc

	routes.muxRouter.HandleFunc("/transactions", handlerFunc(routes.controller.CreateTransaction)).Methods("POST")
}
