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
}

func NewAccountMiddleware(db *gorm.DB, router *mux.Router) *AccountMiddleware {

	return &AccountMiddleware{db: db, router: router}
}

func (mw *AccountMiddleware) Init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{}) // Use JSON format for logs
	logger.SetLevel(logrus.InfoLevel)

	middlewareHandler := middlewareHandlerPathV1.NewMiddlewareHandler(logger)
	repoV1 := repoV1Path.NewAccountRepository(logger)
	coreV1 := coreV1Path.NewAccountCore(repoV1, logger)
	controllerV1 := controllerV1Path.NewAccountController(repoV1, coreV1, mw.db, logger)
	router := routerV1Path.NewAccountRoutes(controllerV1, mw.router, middlewareHandler)
	router.Init()
}
