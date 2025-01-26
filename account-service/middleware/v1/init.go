package account_middleware_v1

import (
	controllerV1Package "anti-fraud/account-service/controllers/v1"
	coreV1Package "anti-fraud/account-service/core/v1"
	repoV1Package "anti-fraud/account-service/repository/v1"
	routerV1Package "anti-fraud/account-service/routes/v1"

	clientV1Package "anti-fraud/mediator-service/account-service-client"
	middlewareHandlerPackageV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountMiddleware struct {
	db     *gorm.DB
	router *mux.Router
	logger *logrus.Logger
	coreV1 *coreV1Package.AccountCore
}

func NewAccountMiddleware(db *gorm.DB, router *mux.Router, logger *logrus.Logger) *AccountMiddleware {

	return &AccountMiddleware{db: db, router: router, logger: logger}
}
func (mw *AccountMiddleware) Init() {

	middlewareHandler := middlewareHandlerPackageV1.NewMiddlewareHandler(mw.logger)
	repoV1 := repoV1Package.NewAccountRepository(mw.logger)
	mw.coreV1 = coreV1Package.NewAccountCore(repoV1, mw.logger)
	controllerV1 := controllerV1Package.NewAccountController(repoV1, mw.coreV1, mw.db, mw.logger)
	router := routerV1Package.NewAccountRoutes(controllerV1, mw.router, middlewareHandler)
	router.Init()
}

func (mw *AccountMiddleware) ConfigureClient(client *clientV1Package.AccountClient) {
	client.SetupCore(mw.coreV1)
}
