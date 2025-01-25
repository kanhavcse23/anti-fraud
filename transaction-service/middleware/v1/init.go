package transaction_middleware_v1

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionMiddleware struct {
	db *mongo.Database
}

func NewTransactionMiddleware(db *mongo.Database) *TransactionMiddleware {

	return &TransactionMiddleware{db: db}
}

func (mw *TransactionMiddleware) Init() {
	// controllerV1 := &controllerV1Path.AccountController{}
	// repoV1 := &repoV1Path.AccountRepository{}
}
