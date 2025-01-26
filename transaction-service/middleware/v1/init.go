package transaction_middleware_v1

import (
	controllerV1Package "anti-fraud/transaction-service/controllers/v1"
	coreV1Package "anti-fraud/transaction-service/core/v1"
	repoV1Package "anti-fraud/transaction-service/repository/v1"
	routerV1Package "anti-fraud/transaction-service/routes/v1"

	operationClientV1Package "anti-fraud/mediator-service/operation-service-client"
	middlewareHandlerV1Package "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionMiddleware struct {
	db              *gorm.DB
	router          *mux.Router
	logger          *logrus.Logger
	operationClient *operationClientV1Package.OperationClient
}

func NewTransactionMiddleware(db *gorm.DB, router *mux.Router, logger *logrus.Logger, operationClient *operationClientV1Package.OperationClient) *TransactionMiddleware {

	return &TransactionMiddleware{db: db, router: router, logger: logger, operationClient: operationClient}
}

func (mw *TransactionMiddleware) Init() {

	middlewareHandler := middlewareHandlerV1Package.NewMiddlewareHandler(mw.logger)
	repoV1 := repoV1Package.NewTransactionRepository(mw.logger)
	coreV1 := coreV1Package.NewTransactionCore(repoV1, mw.logger, mw.operationClient)
	controllerV1 := controllerV1Package.NewTransactionController(repoV1, coreV1, mw.db, mw.logger)
	router := routerV1Package.NewTransactionRoutes(controllerV1, mw.router, middlewareHandler)
	router.Init()
}
