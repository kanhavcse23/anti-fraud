package account_middleware_v1

import (
	controllerV1Package "anti-fraud/account-service/controllers/v1"
	coreV1Package "anti-fraud/account-service/core/v1"
	repoV1Package "anti-fraud/account-service/repository/v1"
	routerV1Package "anti-fraud/account-service/routes/v1"

	middlewareHandlerV1Package "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountMiddleware struct {
	db     *gorm.DB
	router *mux.Router
	logger *logrus.Logger
}

func NewAccountMiddleware(db *gorm.DB, router *mux.Router, logger *logrus.Logger) *AccountMiddleware {

	return &AccountMiddleware{db: db, router: router, logger: logger}
}

func (mw *AccountMiddleware) Init() {

	middlewareHandler := middlewareHandlerV1Package.NewMiddlewareHandler(mw.logger)
	repoV1 := repoV1Package.NewAccountRepository(mw.logger)
	coreV1 := coreV1Package.NewAccountCore(repoV1, mw.logger)
	controllerV1 := controllerV1Package.NewAccountController(repoV1, coreV1, mw.db, mw.logger)
	router := routerV1Package.NewAccountRoutes(controllerV1, mw.router, middlewareHandler)
	router.Init()
}
