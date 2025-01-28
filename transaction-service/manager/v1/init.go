package transaction_manager_v1

import (
	controllerV1Package "anti-fraud/transaction-service/controllers/v1"
	coreV1Package "anti-fraud/transaction-service/core/v1"
	repoV1Package "anti-fraud/transaction-service/repository/v1"
	routerV1Package "anti-fraud/transaction-service/routes/v1"

	accountClientV1Package "anti-fraud/mediator-service/account-service-client"
	operationClientV1Package "anti-fraud/mediator-service/operation-service-client"
	middlewareHandlerPackageV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TransactionManager wires all components required to run transaction-service.
type TransactionManager struct {
	db              *gorm.DB
	router          *mux.Router
	logger          *logrus.Logger
	operationClient operationClientV1Package.IOperationClient
	accountClient   accountClientV1Package.IAccountClient
}

// NewTransactionManager create and return new instance of TransactionManager.
func NewTransactionManager(db *gorm.DB, router *mux.Router, logger *logrus.Logger, operationClient operationClientV1Package.IOperationClient, accountClient accountClientV1Package.IAccountClient) *TransactionManager {

	return &TransactionManager{db: db, router: router, logger: logger, operationClient: operationClient, accountClient: accountClient}
}

// Init wire all components, register routes for transaction-service.
func (mw *TransactionManager) Init() {

	middlewareHandler := middlewareHandlerPackageV1.NewMiddlewareHandler(mw.logger)
	repoV1 := repoV1Package.NewTransactionRepository(mw.logger)
	coreV1 := coreV1Package.NewTransactionCore(repoV1, mw.logger, mw.operationClient, mw.accountClient)
	controllerV1 := controllerV1Package.NewTransactionController(repoV1, coreV1, mw.db, mw.logger)
	router := routerV1Package.NewTransactionRoutes(controllerV1, mw.router, middlewareHandler)
	router.Init()
}
