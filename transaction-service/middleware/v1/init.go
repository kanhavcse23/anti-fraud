package transaction_middleware_v1

import (
	controllerV1Path "anti-fraud/transaction-service/controllers/v1"
	coreV1Path "anti-fraud/transaction-service/core/v1"
	repoV1Path "anti-fraud/transaction-service/repository/v1"
	routerV1Path "anti-fraud/transaction-service/routes/v1"

	middlewareHandlerPathV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionMiddleware struct {
	db     *gorm.DB
	router *mux.Router
	logger *logrus.Logger
}

func NewTransactionMiddleware(db *gorm.DB, router *mux.Router, logger *logrus.Logger) *TransactionMiddleware {

	return &TransactionMiddleware{db: db, router: router, logger: logger}
}

func (mw *TransactionMiddleware) Init() {

	middlewareHandler := middlewareHandlerPathV1.NewMiddlewareHandler(mw.logger)
	repoV1 := repoV1Path.NewTransactionRepository(mw.logger)
	coreV1 := coreV1Path.NewTransactionCore(repoV1, mw.logger)
	controllerV1 := controllerV1Path.NewTransactionController(repoV1, coreV1, mw.db, mw.logger)
	router := routerV1Path.NewTransactionRoutes(controllerV1, mw.router, middlewareHandler)
	router.Init()
}
