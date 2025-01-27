package account_manager_v1

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

type AccountManager struct {
	db     *gorm.DB
	router *mux.Router
	logger *logrus.Logger
	coreV1 coreV1Package.IAccountCore
}

func NewAccountManager(db *gorm.DB, router *mux.Router, logger *logrus.Logger) *AccountManager {

	return &AccountManager{db: db, router: router, logger: logger}
}
func (mw *AccountManager) Init() {

	managerHandler := middlewareHandlerPackageV1.NewMiddlewareHandler(mw.logger)
	repoV1 := repoV1Package.NewAccountRepository(mw.logger)
	mw.coreV1 = coreV1Package.NewAccountCore(repoV1, mw.logger)
	controllerV1 := controllerV1Package.NewAccountController(repoV1, mw.coreV1, mw.db, mw.logger)
	router := routerV1Package.NewAccountRoutes(controllerV1, mw.router, managerHandler)
	router.Init()
}

func (mw *AccountManager) ConfigureClient(client clientV1Package.IAccountClient) {
	client.SetupCore(mw.coreV1)
}
