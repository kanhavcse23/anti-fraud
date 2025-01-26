package account_middleware_v1

import (
	controllerV1Path "anti-fraud/account-service/controllers/v1"
	coreV1Path "anti-fraud/account-service/core/v1"
	repoV1Path "anti-fraud/account-service/repository/v1"
	routerV1Path "anti-fraud/account-service/routes/v1"

	middlewareHandlerPathV1 "anti-fraud/utils-server/middleware/v1"

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

	middlewareHandler := middlewareHandlerPathV1.NewMiddlewareHandler(mw.logger)
	repoV1 := repoV1Path.NewAccountRepository(mw.logger)
	coreV1 := coreV1Path.NewAccountCore(repoV1, mw.logger)
	controllerV1 := controllerV1Path.NewAccountController(repoV1, coreV1, mw.db, mw.logger)
	router := routerV1Path.NewAccountRoutes(controllerV1, mw.router, middlewareHandler)
	router.Init()
}
